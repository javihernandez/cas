/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package verify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/bom"
	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/bom/docker"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/codenotary/cas/pkg/uri"
	immuschema "github.com/codenotary/immudb/pkg/api/schema"
)

var trustLevelMap = map[string]artifact.TrustLevel{
	"trusted":     artifact.Trusted,
	"t":           artifact.Trusted,
	"unknown":     artifact.Unknown,
	"unk":         artifact.Unknown,
	"uk":          artifact.Unknown,
	"unsupported": artifact.Unsupported,
	"uns":         artifact.Unsupported,
	"us":          artifact.Unsupported,
	"untrusted":   artifact.Untrusted,
	"unt":         artifact.Untrusted,
	"ut":          artifact.Untrusted,
}

var bomTrustLevelToMeta = map[artifact.TrustLevel]meta.Status{
	artifact.Untrusted:   meta.StatusUntrusted,
	artifact.Unsupported: meta.StatusUnsupported,
	artifact.Unknown:     meta.StatusUnknown,
	artifact.Trusted:     meta.StatusTrusted,
}

var ErrInsufficientTrustLevel = errors.New("some dependencies have insufficient trust level")

func processBOM(lcUser *api.LcUser, signerID, output, hash, path string) (artifact.Artifact, error) {
	trustLevel, ok := trustLevelMap[viper.GetString("bom-trust-level")]
	if !ok {
		return nil, fmt.Errorf("invalid BOM trust level, supported values are trusted/unknown/unsupported/untrusted")
	}

	outputOpts := artifact.Progress
	if viper.GetBool("silent") || output != "" {
		outputOpts = artifact.Silent
	}

	var bomArtifact artifact.Artifact
	var deps []artifact.Dependency
	var err error
	if hash != "" {
		// hash specified - resolve dependencies from DB
		bomArtifact, err = loadBomFromDb(hash, signerID, lcUser)
		if err != nil {
			return nil, err
		}
		deps = bomArtifact.Dependencies()
	} else {
		// resolve dependencies from the asset
		u, err := uri.Parse(path)
		if err != nil {
			return nil, err
		}
		if _, ok := bom.BomSchemes[u.Scheme]; !ok {
			return nil, fmt.Errorf("unsupported URI %s for --bom option", path)
		}
		if u.Scheme != "" {
			path = strings.TrimPrefix(u.Opaque, "//")
		}
		if u.Scheme == "docker" {
			bomArtifact, err = docker.New(path)
			if err != nil {
				return nil, err
			}
		} else {
			path, err = filepath.Abs(path)
			if err != nil {
				return nil, err
			}
			bomArtifact = bom.New(path)
		}
		if bomArtifact == nil {
			return nil, fmt.Errorf("unsupported artifact format/language")
		}
		if signerID == "" {
			signerID = api.GetSignerIDByApiKey(lcUser.Client.ApiKey)
		}

		if outputOpts != artifact.Silent {
			fmt.Printf("Resolving dependencies...\n")
		}
		deps, err = bomArtifact.ResolveDependencies(outputOpts)
		if err != nil {
			return nil, fmt.Errorf("cannot get dependencies: %w", err)
		}
	}

	if outputOpts != artifact.Silent {
		fmt.Printf("Authenticating dependencies...\n")
	}
	threshold := viper.GetFloat64("bom-max-unsupported")
	unsupportedCount := 0
	failed := false

	var bar *progressbar.ProgressBar
	if len(deps) > 1 && output == "" && outputOpts == artifact.Progress {
		bar = progressbar.Default(int64(len(deps)))
	}

	progressCallback := func(processedDeps []artifact.Dependency) {
		if bar != nil {
			bar.Add(len(processedDeps))
		}
	}

	bomBatchSize := int(viper.GetUint("bom-batch-size"))

	errs, err := artifact.AuthenticateDependencies(lcUser, signerID, deps, bomBatchSize, progressCallback)
	if err != nil {
		return nil, fmt.Errorf("error authenticating dependencies: %w", err)
	}

	lowestLevel := artifact.Trusted
	for i := range deps { // Authenticate mutates the dependency, so use the index
		if errs[i] != nil {
			fmt.Fprintf(os.Stderr, "cannot authenticate %s@%s dependency: %v\n",
				deps[i].Name, deps[i].Version, errs[i])
			continue
		}
		if deps[i].TrustLevel < trustLevel {
			if deps[i].TrustLevel < lowestLevel {
				lowestLevel = deps[i].TrustLevel
			}
			if deps[i].TrustLevel == artifact.Unsupported || deps[i].TrustLevel == artifact.Unknown {
				unsupportedCount++
			} else {
				failed = true // keep going - process all
			}
		}
	}
	if outputOpts != artifact.Silent {
		artifact.Display(bomArtifact, artifact.ColNameVersion|artifact.ColHash|artifact.ColTrustLevel)
	}
	if threshold < 100 && unsupportedCount > int(float64(len(deps))*threshold/100) {
		failed = true // keep going - user still may need output files
	}

	err = bom.Output(bomArtifact)
	if err != nil {
		// show warning, but not error, because authentication finished
		fmt.Fprintln(os.Stderr, err)
	}

	if failed {
		viper.Set("exit-code", strconv.Itoa(bomTrustLevelToMeta[lowestLevel].Int()))
		return bomArtifact, ErrInsufficientTrustLevel
	}

	return bomArtifact, nil
}

func loadBomFromDb(hash string, signerID string, lcUser *api.LcUser) (artifact.Artifact, error) {
	md := metadata.Pairs(meta.CasPluginTypeHeaderName, meta.CasPluginTypeHeaderValue)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	if signerID == "" {
		signerID = api.GetSignerIDByApiKey(lcUser.Client.ApiKey)
	}

	ar, err := artifact.LoadFromDb(hash, signerID, lcUser)
	if err != nil {
		return nil, err
	}

	deps, err := getDeps(hash, signerID, lcUser, ctx, ".0", artifact.DepDirect)
	if err != nil {
		return nil, err
	}
	deps2, err := getDeps(hash, signerID, lcUser, ctx, "", artifact.DepDirect)
	if err != nil {
		return nil, err
	}
	deps = append(deps, deps2...)

	depsTransient, err := getDeps(hash, signerID, lcUser, ctx, ".1", artifact.DepTransient)
	if err != nil {
		return nil, err
	}
	if len(depsTransient) > 0 {
		deps = append(deps, depsTransient...)
	}

	ar.Deps = deps
	return ar, nil
}

func getDeps(hash string, signerID string, lcUser *api.LcUser, ctx context.Context, suffix string, depType artifact.DepType) ([]artifact.Dependency, error) {
	var trustLevelMap = map[meta.Status]artifact.TrustLevel{
		meta.StatusUntrusted:   artifact.Untrusted,
		meta.StatusUnsupported: artifact.Unsupported,
		meta.StatusUnknown:     artifact.Unknown,
		meta.StatusTrusted:     artifact.Trusted,
	}

	zItems, err := lcUser.Client.ZScanExt(ctx, &immuschema.ZScanRequest{
		// "included_by_vcn" is the prefix used by CNIL server
		Set:    []byte("included_by_vcn." + signerID + "." + hash + suffix),
		NoWait: true,
	})
	if err != nil {
		return nil, err
	}

	deps := make([]artifact.Dependency, 0, len(zItems.Items))
	for _, v := range zItems.Items {
		var p pkg
		err := json.Unmarshal(v.Item.Entry.Value, &p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot parse JSON: %v\n", err)
			continue
		}
		level, ok := trustLevelMap[meta.Status(p.Status)]
		if !ok {
			level = artifact.Unknown
		}

		deps = append(deps, artifact.Dependency{
			Name:       p.Name,
			Version:    p.Md.Version,
			SignerID:   signerID,
			HashType:   artifact.HashTypeByName(p.Md.HashType),
			Hash:       p.Hash,
			TrustLevel: level,
			Kind:       p.Kind,
			Timestamp:  time.Unix(v.Timestamp.GetSeconds(), int64(v.Timestamp.GetNanos())).UTC(),
			Type:       depType,
		})
	}

	return deps, nil
}

func DepsToPackageDetails(deps []artifact.Dependency) []api.PackageDetails {
	res := make([]api.PackageDetails, 0, len(deps))
	for _, deps := range deps {
		res = append(res, api.PackageDetails{
			Name:    deps.Name,
			Version: deps.Version,
			Hash:    deps.Hash,
			Status:  bomTrustLevelToMeta[deps.TrustLevel],
			License: deps.License,
		})
	}
	return res
}

type diff struct {
	Counters counters  `json:"counters"`
	Pkg      []pkgDiff `json:"packages"`
}

type pkgDiff struct {
	Name string      `json:"name"`
	Base *pkgDetails `json:"base,omitempty"`
	Diff *pkgDetails `json:"diff,omitempty"`
}

type pkgDetails struct {
	Action    string `json:"action,omitempty"` // only in `diff`, always empty in `base`
	Version   string `json:"version"`
	Status    string `json:"status"`
	Timestamp string `json:"when,omitempty"`
}

type counters struct {
	Unchanged   int `json:"unchanged"`
	Removed     int `json:"removed"`
	Added       int `json:"added"`
	VerChanged  int `json:"version_changed"`
	HashChanged int `json:"hash_changed"`
}

// compare BOMs, output the difference
// TODO currently this function assumes that artifact may have only one instance of the dependency,
// multiple instances with different versions are not supported
func diffBOMs(first, second artifact.Artifact) error {
	firstDeps := first.Dependencies()
	if len(firstDeps) == 0 {
		return fmt.Errorf("artifact %s has no dependencies - nothing to compare", first.Path())
	}
	secondDeps := second.Dependencies()
	if len(secondDeps) == 0 {
		return fmt.Errorf("artifact %s has no dependencies - nothing to compare", second.Path())
	}
	depMap := make(map[string]*artifact.Dependency, len(firstDeps))
	for i := range firstDeps {
		depMap[firstDeps[i].Name] = &firstDeps[i]
	}

	var res diff
	for _, newDep := range secondDeps {
		pkg := pkgDiff{Name: newDep.Name}
		pkg.Diff = &pkgDetails{Version: newDep.Version, Status: artifact.TrustLevelName(newDep.TrustLevel)}
		if newDep.TrustLevel != artifact.Unknown {
			pkg.Diff.Timestamp = newDep.Timestamp.Format(time.RFC3339)
		}
		base, ok := depMap[newDep.Name]
		if !ok {
			pkg.Diff.Action = "added"
			res.Counters.Added++
			res.Pkg = append(res.Pkg, pkg)
		} else {
			if base.Hash != newDep.Hash {
				pkg.Base = &pkgDetails{
					Version:   base.Version,
					Status:    artifact.TrustLevelName(base.TrustLevel),
					Timestamp: base.Timestamp.Format(time.RFC3339),
				}
				if base.Version != newDep.Version {
					res.Counters.VerChanged++
					pkg.Diff.Action = "changed"
				} else {
					res.Counters.HashChanged++
					pkg.Diff.Action = "hash_changed"
				}
				res.Pkg = append(res.Pkg, pkg)
			} else {
				res.Counters.Unchanged++
			}
			delete(depMap, newDep.Name) // to make it has been processed
		}
	}

	// anything left in the map was removed
	for k, v := range depMap {
		res.Pkg = append(res.Pkg, pkgDiff{
			Name: k,
			Base: &pkgDetails{
				Version:   v.Version,
				Status:    artifact.TrustLevelName(v.TrustLevel),
				Timestamp: v.Timestamp.Format(time.RFC3339),
			},
			Diff: &pkgDetails{
				Version: "none",
				Status:  "none",
				Action:  "removed",
			},
		})
		res.Counters.Removed++
	}

	out, _ := json.MarshalIndent(res, "", "  ")
	fmt.Printf("%s\n", out)

	return nil
}

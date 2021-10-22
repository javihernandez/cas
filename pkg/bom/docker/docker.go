/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package docker

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/bom/executor"
)

const AssetType = "docker"

const (
	APK   = "apk"
	DPKG  = "dpkg"
	RPM   = "rpm"
	Image = "image"
)

type DockerArtifact struct {
	artifact.GenericArtifact
	image   string
	ex      executor.Executor
	pkg     pkgManager
	pkgType string
}

// New returns new DockerArtifact object
// unlike other environments, this New() is not called from bom.New()
func New(path string) (*DockerArtifact, error) {
	executor, err := executor.NewDockerExecutor(path)
	if err != nil {
		return nil, err
	}

	pkg, err := probePackageManager(executor)
	if err != nil {
		return nil, fmt.Errorf("error identifying package manager for the image: %w", err)
	}
	if pkg == nil {
		return nil, fmt.Errorf("cannot identify package manager for the container image")
	}

	ret := DockerArtifact{
		image:   path,
		ex:      executor,
		pkg:     pkg,
		pkgType: pkg.Type(),
	}
	return &ret, nil
}

func (p DockerArtifact) Type() string {
	return p.pkgType
}

func (p DockerArtifact) Path() string {
	return p.image
}

func (a *DockerArtifact) ResolveDependencies(output artifact.OutputOptions) ([]artifact.Dependency, error) {
	if a.Deps != nil {
		return a.Deps, nil
	}
	defer a.ex.Close()

	result, err := a.pkg.AllPackages(a.ex, output)
	if err != nil {
		return nil, err
	}

	result = a.filterOutDuplicates(result)

	a.Deps = result
	return result, nil
}

func (a *DockerArtifact) filterOutDuplicates(deps []artifact.Dependency) []artifact.Dependency {
	ret := make([]artifact.Dependency, 0, len(deps))
	hashesFound := make(map[string]struct{}, len(deps))
	for _, dep := range deps {
		if _, exists := hashesFound[dep.Hash]; !exists {
			ret = append(ret, dep)
			hashesFound[dep.Hash] = struct{}{}
		}
	}
	return ret
}

func (a *DockerArtifact) FileHash(name string) (string, error) {
	buf, err := a.ex.ReadFile(name)
	if err != nil {
		return "", nil
	}

	checksum := sha256.Sum256(buf)
	return hex.EncodeToString(checksum[:]), nil
}

func IsDocker(a artifact.Artifact) bool {
	aType := a.Type()
	if aType == APK || aType == RPM || aType == DPKG {
		return true
	}
	return false
}

/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package artifact

import (
	"errors"
	"fmt"
	"time"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/meta"
)

type HashType uint

const (
	HashInvalid HashType = iota
	HashSHA1
	HashSHA224
	HashSHA256
	HashSHA384
	HashSHA512
	HashMD2
	HashMD4
	HashMD5
	HashMD6
	minHash = HashSHA1
	maxHash = HashMD6
)

var hashText = [maxHash + 1]string{"Invalid", "SHA1", "SHA224", "SHA256", "SHA384", "SHA512", "MD2", "MD4", "MD5", "MD6"}

type TrustLevel uint

const (
	Invalid TrustLevel = iota
	Untrusted
	Unsupported
	Unknown
	Trusted
	MinTrustLevel = Untrusted
	MaxTrustLevel = Trusted
)

var levelText = [MaxTrustLevel + 1]string{"", "Untrusted", "Unsupported", "Unknown", "Trusted"}

const MaxGoroutines = 8 // used by other packages that query components from external sources

type DepType bool

const (
	DepDirect    DepType = false
	DepTransient DepType = true
)

// Dependency is a single building block, used for building the Artifact
type Dependency struct {
	Name       string
	Version    string
	Hash       string
	Kind       string
	HashType   HashType
	TrustLevel TrustLevel // set by Notorize/Authenticate
	SignerID   string     // set by Notorize/Authenticate
	License    string
	Timestamp  time.Time
	Type       DepType
}

func HashTypeName(hashType HashType) string {
	if hashType > maxHash {
		return hashText[HashInvalid]
	}
	return hashText[hashType]
}

func HashTypeByName(text string) HashType {
	for i := range hashText {
		if hashText[i] == text {
			return HashType(i)
		}
	}
	return HashInvalid
}

func TrustLevelName(level TrustLevel) string {
	if level > MaxTrustLevel {
		return "" // return empty string to avoid printing SPDX comment for invalid level
	}
	return levelText[level]
}

// AuthenticateDependencies ...
func AuthenticateDependencies(
	lcUser *api.LcUser,
	signerID string,
	deps []Dependency,
	batchSize int,
	progressCallback func([]Dependency),
) ([]error, error) {
	if len(deps) == 0 {
		return nil, nil
	}

	hashes := make([]string, 0, len(deps))
	for _, dep := range deps {
		hashes = append(hashes, dep.Hash)
	}

	if batchSize <= 0 {
		batchSize = len(hashes)
	}
	nbBatches := len(hashes) / int(batchSize)
	if nbBatches*batchSize < len(hashes) {
		nbBatches++
	}

	var artifacts []*api.LcArtifact
	var verified []bool
	var errs []error

	for i := 0; i < nbBatches; i++ {
		startAt := i * batchSize
		endBefore := startAt + batchSize
		if endBefore > len(hashes) {
			endBefore = len(hashes)
		}

		currArtifacts, currVerified, currErrs, err := lcUser.LoadArtifacts(signerID, hashes[startAt:endBefore], nil)
		if progressCallback != nil {
			progressCallback(deps[startAt:endBefore])
		}
		if err != nil {
			return nil, err
		}

		artifacts = append(artifacts, currArtifacts...)
		verified = append(verified, currVerified...)
		errs = append(errs, currErrs...)
	}

	retErrs := make([]error, len(deps))
	for i := 0; i < len(deps); i++ {
		level := Unknown
		err := errs[i]
		if err == nil {
			switch {
			case !verified[i]:
				return nil, errors.New("the ledger is compromised")
			case artifacts[i].Status == meta.StatusUntrusted || (artifacts[i].Revoked != nil && !artifacts[i].Revoked.IsZero()):
				level = Untrusted
			case artifacts[i].Status == meta.StatusUnsupported:
				level = Unsupported
			default:
				level = Trusted
			}
			deps[i].Timestamp = artifacts[i].Timestamp.UTC()
		} else if err != api.ErrNotFound {
			retErrs[i] = err
		}

		deps[i].TrustLevel = level
		if level != Unknown {
			deps[i].SignerID = signerID
		}
	}

	return retErrs, nil
}

// NotarizeDependencies ...
func NotarizeDependencies(
	lcUser *api.LcUser,
	kinds []string,
	deps []*Dependency,
	batchSize int,
	progressCallback func([]*Dependency),
) error {
	if len(deps) == 0 {
		return nil
	}

	if len(kinds) != len(deps) {
		return fmt.Errorf("number of kinds (%d) and dependencies (%d) must match", len(kinds), len(deps))
	}

	if batchSize <= 0 {
		batchSize = len(deps)
	}
	nbBatches := len(deps) / int(batchSize)
	if nbBatches*batchSize < len(deps) {
		nbBatches++
	}

	for i := 0; i < nbBatches; i++ {
		startAt := i * batchSize
		endBefore := startAt + batchSize
		if endBefore > len(deps) {
			endBefore = len(deps)
		}
		currDeps := deps[startAt:endBefore]

		artifacts := make([]*api.Artifact, len(currDeps))
		options := make([][]api.LcSignOption, len(currDeps))
		for i, dep := range currDeps {
			artifacts[i] = ToApiArtifact(kinds[i], dep.Name, dep.Version, dep.Hash, dep.HashType)
			options[i] = []api.LcSignOption{api.LcSignWithStatus(meta.StatusTrusted)}
		}

		_, err := lcUser.SignMulti(artifacts, options)
		if progressCallback != nil {
			progressCallback(currDeps)
		}
		if err != nil {
			return fmt.Errorf("notarization of %d dependencies failed: %w", len(deps), err)
		}
	}

	signerID := api.GetSignerIDByApiKey(lcUser.Client.ApiKey)
	for i := range deps {
		deps[i].TrustLevel = Trusted
		deps[i].SignerID = signerID
	}

	return nil
}

// ToApiArtifact ...
func ToApiArtifact(kind, name, version, hash string, hashType HashType) *api.Artifact {
	return &api.Artifact{
		Kind: kind,
		Name: name,
		Hash: hash,
		Size: 0,
		Metadata: map[string]interface{}{
			"version":  version,
			"hashtype": HashTypeName(hashType)}}
}

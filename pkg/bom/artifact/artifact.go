/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package artifact

import (
	"fmt"
	"sort"

	"github.com/fatih/color"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/meta"
)

type ColumnID uint

const (
	ColNameVersion ColumnID = 1 << iota
	ColHash
	ColTrustLevel
	MaxColumn = iota
)

type OutputOptions uint

const (
	Silent OutputOptions = iota
	Progress
)

// Artifact is a result of build process.
// It is a language- and/or environment-specific interface which finds dependencies
type Artifact interface {
	Path() string
	Type() string
	Dependencies() []Dependency
	ResolveDependencies(output OutputOptions) ([]Dependency, error)
}

type GenericArtifact struct {
	Deps []Dependency
}

type LoadedArtifact struct {
	GenericArtifact
	path string
	kind string
}

func (a GenericArtifact) Dependencies() []Dependency {
	return a.Deps
}

func LoadFromDb(hash string, signerID string, lcUser *api.LcUser) (*LoadedArtifact, error) {
	ar, _, err := lcUser.LoadArtifact(hash, signerID, "", 0, nil)
	if err != nil {
		return nil, err
	}

	return &LoadedArtifact{kind: ar.Kind, path: ar.Name}, nil
}

func Display(a Artifact, columns ColumnID) {
	deps := a.Dependencies()

	maxColWidth := make([]int, MaxColumn)
	sort.SliceStable(deps, (func(q, b int) bool { return deps[q].Name < deps[b].Name }))
	for _, dep := range deps {
		width := len(dep.Name) + len(dep.Version) + 1
		if width > maxColWidth[0] {
			maxColWidth[0] = width
		}
		width = len(dep.Hash)
		if width > maxColWidth[1] {
			maxColWidth[1] = width
		}
	}

	for _, dep := range deps {
		if ColNameVersion&columns != 0 {
			fmt.Printf("%-*s", maxColWidth[0]+1, dep.Name+"@"+dep.Version)
		}
		if ColHash&columns != 0 {
			fmt.Printf("%-*s", maxColWidth[1]+1, dep.Hash)
		}
		if ColTrustLevel&columns != 0 {
			switch dep.TrustLevel {
			case Trusted:
				color.Set(meta.StyleSuccess())
			case Unknown:
				color.Set(meta.StyleWarning())
			case Unsupported, Untrusted:
				color.Set(meta.StyleError())
			}
			fmt.Print(TrustLevelName(dep.TrustLevel))
			color.Unset()
		}
		fmt.Println()
	}
}

func (a LoadedArtifact) Type() string {
	return a.kind
}

func (a LoadedArtifact) Path() string {
	return a.path
}

func (a LoadedArtifact) Dependencies() []Dependency {
	return a.Deps
}

func (a LoadedArtifact) ResolveDependencies(output OutputOptions) ([]Dependency, error) {
	return a.Deps, nil
}

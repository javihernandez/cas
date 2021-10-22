/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package api

import (
	"fmt"

	"github.com/codenotary/cas/pkg/meta"
)

// Artifact represents the set of all relevant information gathered from a digital asset.
type Artifact struct {
	Kind        string
	Name        string
	Hash        string
	Size        uint64
	ContentType string
	IncludedIn  []PackageDetails
	Deps        []PackageDetails
	Metadata
}

type PackageDetails struct {
	Name    string      `json:"name" yaml:"name" cas:"name"`
	Version string      `json:"version,omitempty" yaml:"version,omitempty" cas:"version"`
	Hash    string      `json:"hash" yaml:"hash" cas:"hash"`
	Status  meta.Status `json:"status" yaml:"status" cas:"status"`
	License string      `json:"license,omitempty" yaml:"license"`
}

// Copy returns a deep copy of the artifact.
func (a Artifact) Copy() Artifact {
	c := a
	if a.Metadata != nil {
		c.Metadata = nil
		c.Metadata.SetValues(a.Metadata)
	}
	return c
}

// ArtifactResponse holds artifact values returned by the platform.
type ArtifactResponse struct {
	// root fields
	Kind        string `json:"kind" yaml:"kind" cas:"Kind"`
	Name        string `json:"name" yaml:"name" cas:"Name"`
	Hash        string `json:"hash" yaml:"hash" cas:"Hash"`
	Size        uint64 `json:"size" yaml:"size" cas:"Size"`
	ContentType string `json:"contentType" yaml:"contentType" cas:"ContentType"`
	URL         string `json:"url" yaml:"url" cas:"URL"`

	// custom metadata
	Metadata Metadata `json:"metadata" yaml:"metadata" cas:"Metadata"`

	// ArtifactResponse specific
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
}

func (a ArtifactResponse) String() string {
	return fmt.Sprintf("Name:\t%s\nHash:\t%s\nStatus:\t%s\n\n",
		a.Name, a.Hash, a.Status)
}

// Artifact returns an new *Artifact from a
func (a ArtifactResponse) Artifact() *Artifact {
	return &Artifact{
		// root fields
		Kind:        a.Kind,
		Name:        a.Name,
		Hash:        a.Hash,
		Size:        a.Size,
		ContentType: a.ContentType,

		// custom metadata
		Metadata: a.Metadata,
	}
}

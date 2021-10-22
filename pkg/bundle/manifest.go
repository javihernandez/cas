/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package bundle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	// See https://github.com/opencontainers/go-digest#usage
	_ "crypto/sha256"
	_ "crypto/sha512"

	digest "github.com/opencontainers/go-digest"
)

const (
	// ManifestSchemaVersion is the current manifest schema version.
	ManifestSchemaVersion = 1

	// ManifestFilename is the default filename for manifest when stored.
	ManifestFilename = ".cas.manifest.json"

	// ManifestDigestAlgo is the only supported digest's algorithm by current manifest schema version.
	ManifestDigestAlgo = digest.SHA256
)

// Manifest provides bundle structure when marshalled to JSON.
//
// Specifications (version 1):
//  - `schemaVersion` is the version number of the current specification (MUST be always 1 in this case)
//  - fields order is defined as per code structs definitions, orderning MUST NOT be changed
//  - `items` MUST be sorted by its digest's value (lexically byte-wise)
//  - multiple `items` MUST NOT have the same digest value
//  - `items.paths` MUST be sorted by value (lexically byte-wise)
//  - across the same manifest multiple `items.paths`'s elements MUST NOT have the value
//  - json representation of the manifest MUST NOT be indented
//  - sha256 is the only digest's algorithm that MUST be used
//
// The Normalize() method provides sorting funcionality and specification enforcement. It's implictly called
// when the manifest is marshalled.
type Manifest struct {
	// SchemaVersion is the manifest schema that this bundle follows
	SchemaVersion uint `json:"schemaVersion"`

	// Items is an ordered list of items referenced by the manifest.
	Items []Descriptor `json:"items"`
}

// MarshalJSON implements the json.Marshaler interface.
func (m *Manifest) MarshalJSON() ([]byte, error) {
	if err := m.Normalize(); err != nil {
		return nil, err
	}

	type alias Manifest
	mm := alias(*m)
	return json.Marshal(mm)
}

// Normalize deduplicates and sorts items and items's paths in accordance with manifest's schema specs.
// An error is returned when duplicate paths across different items are found, or if digest's algo
// does not match sha256.
func (m *Manifest) Normalize() error {
	if m == nil {
		return fmt.Errorf("nil manifest")
	}

	if m.SchemaVersion != ManifestSchemaVersion {
		return fmt.Errorf("unsupported bundle.Manifest schema version: %d", m.SchemaVersion)
	}

	// make unique index
	idx := make(map[string]Descriptor, len(m.Items))
	for _, d := range m.Items {
		k := d.Digest.String()
		if dd, ok := idx[k]; ok {
			if d.Size != dd.Size {
				return fmt.Errorf(
					"distinct sizes found for same digest (%s): %d, %d",
					d.Digest.String(),
					d.Size,
					dd.Size,
				)
			}
			dd.Paths = append(dd.Paths, d.Paths...)
			idx[k] = dd
		} else {
			idx[k] = d
		}
	}

	// recreate unique digest list and sort paths
	m.Items = make([]Descriptor, len(idx))
	paths := make(map[string]bool)
	i := 0
	for _, d := range idx {
		d.sortUnique()
		m.Items[i] = d
		i++

		// specs enforcement:
		// - the only allowed digest's algo is SHA256
		// - within the same manifest multiple paths elements with same value are NOT allowed
		if algo := d.Digest.Algorithm(); algo != ManifestDigestAlgo {
			return fmt.Errorf("unsupported digest algorithm: %s", string(algo))
		}
		for _, p := range d.Paths {
			if paths[p] {
				return fmt.Errorf("duplicate path in manifest: %s", p)
			}
			paths[p] = true
		}
	}

	// finally, sort items by digest
	sort.SliceStable(m.Items, func(k, j int) bool {
		return m.Items[k].Digest.String() < m.Items[j].Digest.String()
	})
	return nil
}

// Digest digests the JSON encoded m and returns a digest.Digest.
func (m *Manifest) Digest() (digest.Digest, error) {
	b, err := json.Marshal(m) // sorting is implicitly called by Marshal
	if err != nil {
		return "", err
	}

	return digest.SHA256.FromBytes(b), nil
}

// NewManifest returns a new empty Manifest.
func NewManifest(items ...Descriptor) *Manifest {
	if items == nil {
		items = make([]Descriptor, 0)
	}
	return &Manifest{
		SchemaVersion: ManifestSchemaVersion,
		Items:         items,
	}
}

// WriteManifest writes manifest's data to a file named by filename.
func WriteManifest(manifest Manifest, filename string) error {
	data, err := json.Marshal(&manifest)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// ReadManifest reads the file named by filename and returns the decoded manifest.
func ReadManifest(filename string) (*Manifest, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	d := digest.SHA256.FromBytes(data)

	m := Manifest{}
	json.Unmarshal(data, &m)

	dd, err := m.Digest()
	if err != nil {
		return nil, err
	}
	if dd != d {
		return nil, fmt.Errorf("manifest integrity check failed")
	}

	return &m, nil
}

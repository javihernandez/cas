/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package extractor

import (
	"fmt"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/uri"
)

var extractors = map[string]Extractor{}

// Extractor extract an api.Artifact referenced by the given uri.URI.
type Extractor func(*uri.URI, ...Option) ([]*api.Artifact, error)

// Register the Extractor e for the given scheme
func Register(scheme string, e Extractor) {
	extractors[scheme] = e
}

// Schemes returns the list of registered schemes.
func Schemes() []string {
	schemes := make([]string, len(extractors))
	i := 0
	for scheme := range extractors {
		schemes[i] = scheme
		i++
	}
	return schemes
}

// Extract returns an []*api.Artifact for the given rawURIs.
func Extract(rawURIs []string, options ...Option) ([]*api.Artifact, error) {
	artifacts := make([]*api.Artifact, 0)
	for _, ru := range rawURIs {
		u, err := uri.Parse(ru)
		if err != nil {
			return nil, err
		}
		if e, ok := extractors[u.Scheme]; ok {
			ars, err := e(u, options...)
			if err != nil {
				return nil, err
			}
			artifacts = append(artifacts, ars...)
		} else {
			return nil, fmt.Errorf("%s scheme not yet supported", u.Scheme)
		}
	}
	return artifacts, nil
}

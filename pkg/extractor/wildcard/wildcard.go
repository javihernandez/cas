/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package wildcard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codenotary/cas/pkg/extractor/file"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/extractor"
	"github.com/codenotary/cas/pkg/uri"
)

// Scheme for dir
const Scheme = "wildcard"

// ManifestKey is the metadata's key for storing the manifest
const ManifestKey = "manifest"

// PathKey is the metadata's key for the directory path
const PathKey = "path"

type opts struct {
	initIgnoreFile    bool
	skipIgnoreFileErr bool
}

// Artifact returns a file *api.Artifact from a given u
func Artifact(u *uri.URI, options ...extractor.Option) ([]*api.Artifact, error) {

	if u.Scheme != "" && u.Scheme != Scheme {
		return nil, nil
	}

	opts := &opts{}
	if err := extractor.Options(options).Apply(opts); err != nil {
		return nil, err
	}

	path := strings.TrimPrefix(u.Opaque, "//")
	wildcard := filepath.Base(path)
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// provided path is a file
	if fileInfo, err := os.Stat(p); err == nil {
		if fileInfo.IsDir() {
			return nil, fmt.Errorf("folder notarization is not allowed")
		}
		u, err := uri.Parse("file://" + p)
		if err != nil {
			return nil, err
		}
		return file.Artifact(u)
	}

	root := filepath.Dir(p)

	// build a list of all files matching the wildcard provided. Method is based on filepath.Glob
	var filePaths []string

	i, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	err = buildFilePaths(wildcard, &filePaths)(root, i, nil)
	if err != nil {
		return nil, err
	}

	if len(filePaths) == 0 {
		return nil, errors.New("no matching files found")
	}

	arst := []*api.Artifact{}
	// convert files path list to artifacts
	for _, fp := range filePaths {
		u, err := uri.Parse("file://" + fp)
		if err != nil {
			return nil, err
		}
		ars, err := file.Artifact(u)
		if err != nil {
			return nil, err
		}
		arst = append(arst, ars...)
	}

	return arst, nil
}

func buildFilePaths(wildcard string, filePaths *[]string) func(ele string, info os.FileInfo, err error) error {
	return func(ele string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fpd, err := filepath.Glob(filepath.Join(ele, wildcard))
			if err != nil {
				return err
			}
			if len(fpd) > 0 {
				for _, fp := range fpd {
					info, err = os.Stat(fp)
					if err != nil {
						return err
					}
					if !info.IsDir() {
						*filePaths = append(*filePaths, fp)
					}
				}
			}
		}
		return nil
	}
}

/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package store

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"

	"github.com/codenotary/cas/pkg/bundle"
)

func ManifestFilepath(kind string, target string) (string, error) {
	target, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, defaultManifestsDir)
	if err := ensureDir(path); err != nil {
		return "", err
	}
	id := sha256.Sum256([]byte(target))

	return filepath.Join(path, fmt.Sprintf("%s_%x.json", kind, id)), nil
}

func SaveManifest(kind string, target string, manifest bundle.Manifest) error {
	path, err := ManifestFilepath(kind, target)
	if err != nil {
		return err
	}
	return bundle.WriteManifest(manifest, path)
}

func ReadManifest(kind string, target string) (*bundle.Manifest, error) {
	path, err := ManifestFilepath(kind, target)
	if err != nil {
		return nil, err
	}
	return bundle.ReadManifest(path)
}

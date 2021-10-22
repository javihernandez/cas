/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package file

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/extractor"
	"github.com/codenotary/cas/pkg/uri"
)

// Scheme for file
const Scheme = "file"

// Artifact returns a file *api.Artifact from a given u
func Artifact(u *uri.URI, options ...extractor.Option) ([]*api.Artifact, error) {

	if u.Scheme != Scheme {
		return nil, nil
	}

	path := strings.TrimPrefix(u.Opaque, "//")

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Metadata container
	m := api.Metadata{}

	// Hash
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	checksum := h.Sum(nil)

	// Name and Size
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// ContentType
	ct, err := contentType(f)
	if err != nil {
		return nil, err
	}

	// Infer version from filename
	if version := inferVer(stat.Name()); version != "" {
		m["version"] = version
	}

	// Sniff executable info, if any
	if ok, data, _ := xInfo(f, &ct); ok {
		m.SetValues(data)
	}

	return []*api.Artifact{{
		Kind:        Scheme,
		Name:        stat.Name(),
		Hash:        hex.EncodeToString(checksum),
		Size:        uint64(stat.Size()),
		ContentType: ct,
		Metadata:    m,
	}}, nil
}

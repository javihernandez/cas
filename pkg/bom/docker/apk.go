/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package docker

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/bom/executor"
)

// apk implements packageManager interface
type apk struct {
	cache []artifact.Dependency
	index map[string]*artifact.Dependency
}

const (
	checksumTag = 'C'
	packageTag  = 'P'
	versionTag  = 'V'
	licenseTag  = 'L'
)

func (pkg apk) Type() string {
	return APK
}

// AllPackages finds all installed packages
func (pkg *apk) AllPackages(e executor.Executor, output artifact.OutputOptions) ([]artifact.Dependency, error) {
	if pkg.cache == nil {
		err := pkg.buildCache(e)
		if err != nil {
			return nil, err
		}
	}

	return pkg.cache, nil
}

// build package cache and index - see https://wiki.alpinelinux.org/wiki/Apk_spec
func (pkg *apk) buildCache(e executor.Executor) error {
	buf, err := e.ReadFile("/lib/apk/db/installed")
	if err != nil {
		return fmt.Errorf("error reading file from container: %w", err)
	}

	pkg.cache = make([]artifact.Dependency, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(buf))
	curPkg := artifact.Dependency{Type: artifact.DepDirect}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// end of package section
			pkg.cache = append(pkg.cache, curPkg)
			curPkg = artifact.Dependency{}
			continue
		}
		tag := line[0]
		switch tag {
		case packageTag:
			curPkg.Name = line[2:]
		case versionTag:
			curPkg.Version = line[2:]
		case checksumTag:
			hash, err := base64.StdEncoding.DecodeString(line[4:]) // first two bytes contain Q1 to indicate SHA1
			if err != nil {
				return fmt.Errorf("malformed package checksum")
			}
			curPkg.Hash = hex.EncodeToString(hash)
			curPkg.HashType = artifact.HashSHA1
		case licenseTag:
			curPkg.License = line[2:]
		}
	}
	// last line in 'installed' file is empty, therefore last created curPkg just discarded

	pkg.index = make(map[string]*artifact.Dependency, len(pkg.cache))
	for i := range pkg.cache {
		pkg.index[pkg.cache[i].Name] = &pkg.cache[i]
	}

	return nil
}

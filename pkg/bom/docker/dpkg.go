/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/bom/executor"
)

// dpkg implements packageManager interface
type dpkg struct {
	cache []artifact.Dependency
	index map[string]*artifact.Dependency
}

var (
	licensePattern           = regexp.MustCompile(`^License: (\S*)`)
	commonLicensePathPattern = regexp.MustCompile(`/usr/share/common-licenses/([0-9A-Za-z_.\-]+)`)
)

func (pkg dpkg) Type() string {
	return DPKG
}

// AllPackages finds all installed packages
func (pkg *dpkg) AllPackages(e executor.Executor, output artifact.OutputOptions) ([]artifact.Dependency, error) {
	if pkg.cache == nil {
		err := pkg.buildCache(e)
		if err != nil {
			return nil, err
		}
	}

	return pkg.cache, nil
}

func (pkg *dpkg) buildCache(e executor.Executor) error {
	buf, err := e.ReadFile("/var/lib/dpkg/status")
	if err != nil {
		buf, err = e.ReadFile("/var/lib/dpkg/status.d")
		if err != nil {
			return fmt.Errorf("error reading file from container: %w", err)
		}
	}

	pkg.cache = make([]artifact.Dependency, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(buf))
	curPkg := artifact.Dependency{Type: artifact.DepDirect}
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.SplitN(line, ": ", 2)
		if len(fields) != 2 {
			continue
		}
		switch fields[0] {
		case "Package":
			if curPkg.Name != "" {
				pkg.cache = append(pkg.cache, curPkg)
			}
			curPkg = artifact.Dependency{Name: fields[1]}
		case "Version":
			curPkg.Version = fields[1]
		}
	}
	if curPkg.Name != "" {
		pkg.cache = append(pkg.cache, curPkg)
	}

	pkg.index = make(map[string]*artifact.Dependency, len(pkg.cache))
	for i := range pkg.cache {
		pkg.index[pkg.cache[i].Name] = &pkg.cache[i]
	}

	// calculate hashes for all packages
	hashReader, err := e.ReadDir("/var/lib/dpkg/info")
	if err != nil {
		return err
	}
	defer hashReader.Close()
	tr := tar.NewReader(hashReader)

	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if !strings.HasSuffix(hdr.Name, ".md5sums") {
			continue
		}
		// file name has a form of '/var/lib/dpkg/info/<pkg>[:arch].md5sums'
		fields := strings.Split(strings.TrimSuffix(filepath.Base(hdr.Name), ".md5sums"), ":")
		p, ok := pkg.index[fields[0]]
		if !ok {
			continue // unknown package - maybe md5sums file is a leftover
		}
		h := sha256.New()
		if _, err := io.Copy(h, tr); err != nil {
			return err
		}
		p.Hash = hex.EncodeToString(h.Sum(nil))
		p.HashType = artifact.HashSHA256
	}

	// collect license info
	licReader, err := e.ReadDir("/usr/share/doc")
	if err != nil {
		return fmt.Errorf("error reading file from container: %w", err)
	}
	defer licReader.Close()

	tr = tar.NewReader(licReader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading file from container: %w", err)
		}
		fields := strings.Split(hdr.Name, "/")
		if len(fields) != 3 { // expect doc/<package_name>/copyright
			continue
		}
		if fields[2] != "copyright" {
			continue
		}
		pkg, ok := pkg.index[fields[1]]
		if !ok {
			continue
		}
		pkg.License = findLicense(tr)
	}
	return nil
}

// see https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/#license-syntax for details
func findLicense(reader io.Reader) string {
	scanner := bufio.NewScanner(reader)

	license := ""
	for scanner.Scan() {
		line := scanner.Text()
		match := licensePattern.FindStringSubmatch(line)
		if len(match) > 0 {
			license = match[1]
			break
		}
		match = commonLicensePathPattern.FindStringSubmatch(line)
		if len(match) > 0 {
			license = match[1]
			break
		}
	}

	return license
}

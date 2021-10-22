/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package bom

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/codenotary/cas/pkg/bom/artifact"
)

const (
	optional = iota
	mandatory
)

type headerLine struct {
	tag      string
	presense int
	fn       func(artifact.Artifact) (string, error)
}

var headerContent = []headerLine{
	{"SPDXVersion", mandatory, func(artifact.Artifact) (string, error) {
		return "SPDX-2.2", nil
	}},
	{"DataLicense", mandatory, func(artifact.Artifact) (string, error) {
		return "CC0-1.0", nil
	}},
	{"SPDXID", mandatory, func(artifact.Artifact) (string, error) {
		return "SPDXRef-DOCUMENT", nil
	}},
	{"DocumentName", mandatory, func(p artifact.Artifact) (string, error) {
		path, err := filepath.Abs(p.Path())
		if err != nil {
			return "", err
		}
		return filepath.Base(path), nil
	}},
	{"DocumentNamespace", mandatory, func(p artifact.Artifact) (string, error) {
		path, err := filepath.Abs(p.Path())
		if err != nil {
			return "", err
		}
		return "http://spdx.org/spdxdocs/" + filepath.Base(path) + "-" + uuid.NewString(), nil
	}},
	{"Creator", mandatory, func(artifact.Artifact) (string, error) {
		return "Tool: Codenotary cas", nil
	}},
	{"Created", mandatory, func(artifact.Artifact) (string, error) {
		return time.Now().UTC().Format(time.RFC3339), nil
	}},
}

type componentLine struct {
	tag      string
	presense int
	fn       func(artifact.Artifact, artifact.Dependency, int) (string, error)
}

var componentContent = []componentLine{
	{"PackageName", mandatory, func(a artifact.Artifact, d artifact.Dependency, seq int) (string, error) {
		return d.Name, nil
	}},
	{"SPDXID", mandatory, func(a artifact.Artifact, d artifact.Dependency, seq int) (string, error) {
		return "SPDXRef-Package-" + strconv.Itoa(seq), nil
	}},
	{"PackageVersion", optional, func(a artifact.Artifact, d artifact.Dependency, seq int) (string, error) {
		return d.Version, nil
	}},
	{"PackageDownloadLocation", mandatory, noAssertion},
	// FilesAnalysed is optional, but by default it is true, which requires presence of many other fields
	{"FilesAnalyzed", mandatory, func(a artifact.Artifact, d artifact.Dependency, seq int) (string, error) {
		return "false", nil
	}},
	{"PackageChecksum", mandatory, func(a artifact.Artifact, d artifact.Dependency, seq int) (string, error) {
		return artifact.HashTypeName(d.HashType) + ": " + d.Hash, nil
	}},
	{"PackageSourceInfo", optional, func(a artifact.Artifact, d artifact.Dependency, seq int) (string, error) {
		return "<text>" + Purl(a, d) + "</text>", nil
	}},
	{"PackageLicenseConcluded", mandatory, func(a artifact.Artifact, d artifact.Dependency, seq int) (string, error) {
		if d.License != "" {
			return d.License, nil
		}
		return "NOASSERTION", nil
	}},
	{"PackageLicenseDeclared", mandatory, noAssertion},
	{"PackageCopyrightText", mandatory, noAssertion},
	{"PackageComment", optional, func(a artifact.Artifact, d artifact.Dependency, seq int) (string, error) {
		text := "<text>"
		l := artifact.TrustLevelName(d.TrustLevel)
		if l != "" {
			text += l + ", "
		}
		text += DepLinkType(a, d) + ", " + DepType(d) + "</text>"
		return text, nil
	}},
}

func noAssertion(a artifact.Artifact, d artifact.Dependency, seq int) (string, error) {
	return "NOASSERTION", nil
}

// Output info about package and its components in SPDX text (tag:value) format, according to
// SPDX spec 2.2: https://spdx.dev/wp-content/uploads/sites/41/2020/08/SPDX-specification-2-2.pdf
func OutputSpdxText(a artifact.Artifact, filename string) error {
	deps := a.Dependencies()
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	// SPDX header
	for _, line := range headerContent {
		value, err := line.fn(a)
		if err != nil {
			if line.presense == mandatory {
				return fmt.Errorf("cannot get value for tag %s: %w", line.tag, err)
			}
			continue // optional tag - ignore error
		}
		if value == "" {
			if line.presense == mandatory {
				return fmt.Errorf("no value found for mandatory header tag %s", line.tag)
			}
			continue // optional
		}
		if _, err = fmt.Fprintf(f, "%s: %s\n", line.tag, value); err != nil {
			return err
		}
	}

	if _, err = fmt.Fprintf(f, "\n##### Software components\n\n"); err != nil {
		return err
	}

	for i, dep := range deps {
		for _, line := range componentContent {
			value, err := line.fn(a, dep, i+1)
			if err != nil {
				if line.presense == mandatory {
					return fmt.Errorf("cannot get value for tag %s for component %s: %w", line.tag, dep.Name, err)
				}
				continue // optional tag - ignore error
			}
			if value == "" {
				if line.presense == mandatory {
					return fmt.Errorf("no value found for mandatory component tag %s for component %s", line.tag, dep.Name)
				}
				continue // optional
			}
			if _, err = fmt.Fprintf(f, "%s: %s\n", line.tag, value); err != nil {
				return err
			}
		}
		if _, err = fmt.Fprintln(f); err != nil {
			return err
		}
	}

	return nil
}

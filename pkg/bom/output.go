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

	cyclonedx "github.com/CycloneDX/cyclonedx-go"
	"github.com/spf13/viper"

	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/bom/docker"
)

const (
	dynamicLinkage = "Dynamic"
	staticLinkage  = "Static"
	directDep      = "Direct"
	transientDep   = "Transient"
)

func Output(a artifact.Artifact) error {
	filename := viper.GetString("bom-spdx")
	if filename != "" {
		err := OutputSpdxText(a, filename)
		if err != nil {
			return fmt.Errorf("cannot output SPDX: %w", err)
		}
	}

	filename = viper.GetString("bom-cdx-json")
	if filename != "" {
		err := OutputCycloneDX(a, filename, cyclonedx.BOMFileFormatJSON)
		if err != nil {
			return fmt.Errorf("cannot output  CycloneDX JSON: %w", err)
		}
	}

	filename = viper.GetString("bom-cdx-xml")
	if filename != "" {
		err := OutputCycloneDX(a, filename, cyclonedx.BOMFileFormatXML)
		if err != nil {
			return fmt.Errorf("cannot output  CycloneDX XML: %w", err)
		}
	}

	return nil
}

// DepLinkType returns the link type
// static linkage only in following cases:
// - artifact is a container
// - artifact is a Go binary and dependency is Go
func DepLinkType(a artifact.Artifact, d artifact.Dependency) string {

	if docker.IsDocker(a) {
		return staticLinkage
	}

	return dynamicLinkage
}

func DepType(d artifact.Dependency) string {
	if d.Type == artifact.DepDirect {
		return directDep
	}
	return transientDep
}

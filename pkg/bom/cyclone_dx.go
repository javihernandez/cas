/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package bom

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"

	cdx "github.com/CycloneDX/cyclonedx-go"

	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/meta"
)

var hashNames = map[artifact.HashType]cdx.HashAlgorithm{
	artifact.HashMD5:    cdx.HashAlgoMD5,
	artifact.HashSHA1:   cdx.HashAlgoSHA1,
	artifact.HashSHA256: cdx.HashAlgoSHA256,
	artifact.HashSHA384: cdx.HashAlgoSHA384,
	artifact.HashSHA512: cdx.HashAlgoSHA512,
}

func OutputCycloneDX(a artifact.Artifact, filename string, format cdx.BOMFileFormat) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	encoder := cdx.NewBOMEncoder(buf, format)
	encoder.SetPretty(true)

	bom := convertToCyclone(a)

	err = encoder.Encode(bom)
	if err != nil {
		return err
	}

	c := buf.Bytes()

	_, err = f.Write(c)
	return err
}

func convertToCyclone(a artifact.Artifact) *cdx.BOM {
	bom := cdx.NewBOM()

	bom.Metadata = &cdx.Metadata{
		Tools: &[]cdx.Tool{
			{
				Vendor:  "Codenotary",
				Name:    "cas",
				Version: meta.Version(),
			},
		},
	}

	name := filepath.Base(a.Path())
	bom.Metadata.Component = &cdx.Component{
		Name: name,
		Type: cdx.ComponentTypeApplication,
	}

	deps := a.Dependencies()
	comps := make([]cdx.Component, len(deps))
	for i, dep := range deps {
		hashName, ok := hashNames[dep.HashType]
		if !ok {
			hashName = ""
		}
		pkgUrl := Purl(a, dep)
		comps[i] = cdx.Component{
			BOMRef:     name + "-" + strconv.Itoa(i+1),
			Type:       cdx.ComponentTypeLibrary,
			Name:       dep.Name,
			Version:    dep.Version,
			PackageURL: pkgUrl,
			Hashes: &[]cdx.Hash{
				{
					Algorithm: hashName,
					Value:     dep.Hash,
				},
			},
		}
		if dep.License != "" {
			comps[i].Licenses = &cdx.Licenses{cdx.LicenseChoice{Expression: dep.License}}
		}
		props := make([]cdx.Property, 1, 2)
		props[0] = cdx.Property{Name: "LinkType", Value: DepLinkType(a, dep)}
		trustLevel := artifact.TrustLevelName(dep.TrustLevel)
		if trustLevel != "" {
			props = append(props, cdx.Property{Name: "TrustLevel", Value: trustLevel})
		}
		comps[i].Properties = &props
	}
	bom.Components = &comps

	return bom
}

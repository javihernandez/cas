/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package bom

import (
	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/bom/docker"
	purl "github.com/package-url/packageurl-go"
)

var typeMap = map[string]string{
	docker.DPKG:  purl.TypeDebian,
	docker.RPM:   purl.TypeRPM,
	docker.Image: purl.TypeDocker,
}

func Purl(a artifact.Artifact, d artifact.Dependency) string {
	assetType := d.Kind
	if assetType == "" {
		assetType = a.Type()
	}
	assetType, ok := typeMap[assetType]
	if !ok {
		assetType = purl.TypeGeneric
	}
	return purl.NewPackageURL(assetType, "", d.Name, d.Version, nil, "").ToString()
}

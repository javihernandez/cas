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
)

// extractor schemes that can be used to point to BOM source
var BomSchemes = map[string]struct{}{"dir": {}, "git": {}, "docker": {}, "": {}}

// New returns Artifact implementation of type, matching the artifact language/environment
func New(filename string) artifact.Artifact {
	return nil
}

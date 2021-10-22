/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferVer(t *testing.T) {
	testCases := map[string]string{
		// Supported
		"cas-v0.4.0-darwin-10.6-amd64":     "0.4.0",
		"cas-v0.4.0-linux-amd64":           "0.4.0",
		"cas-v0.4.0-windows-4.0-amd64.exe": "0.4.0",

		// Unsupported
		"codenotary_cas_0.4.0_setup.exe": "",
	}

	for filename, ver := range testCases {
		assert.Equal(t, ver, inferVer(filename), "wrong version for %s", filename)
	}

}

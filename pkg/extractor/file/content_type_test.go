/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package file

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContenType(t *testing.T) {
	emptyFile, err := ioutil.TempFile("", "TestContenType")
	if err != nil {
		t.Fatal(err)
	}
	txtFile, err := ioutil.TempFile("", "TestContenType")
	if err != nil {
		t.Fatal(err)
	}
	txtFile.Write([]byte{99, 105, 97, 111})

	ct, err := contentType(emptyFile)
	assert.NoError(t, err)
	assert.Empty(t, ct)

	ct, err = contentType(txtFile)
	assert.NoError(t, err)
	assert.Equal(t, "application/octet-stream", ct)
}

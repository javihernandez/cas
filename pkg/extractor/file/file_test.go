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
	"log"
	"path/filepath"

	"github.com/stretchr/testify/assert"

	"os"
	"testing"

	"github.com/codenotary/cas/pkg/uri"
)

func TestFile(t *testing.T) {
	file, err := ioutil.TempFile("", "cas-test-scheme-file")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	err = ioutil.WriteFile(file.Name(), []byte("123\n"), 0644)
	if err != nil {
		log.Fatal(err)
	}
	u, _ := uri.Parse("file://" + file.Name())

	artifacts, err := Artifact(u)
	assert.NoError(t, err)
	assert.NotNil(t, artifacts)
	assert.Equal(t, Scheme, artifacts[0].Kind)
	assert.Equal(t, filepath.Base(file.Name()), artifacts[0].Name)
	assert.Equal(t, "181210f8f9c779c26da1d9b2075bde0127302ee0e3fca38c9a83f5b1dd8e5d3b", artifacts[0].Hash)

}

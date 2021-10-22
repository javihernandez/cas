/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package store

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mkTmpForConfig(t *testing.T) string {
	tdir, err := ioutil.TempDir("", "cas-test-store-config")
	if err != nil {
		t.Fatal(err)
	}
	return tdir
}

func TestLoadConfig(t *testing.T) {
	tdir := mkTmpForConfig(t)
	SetDir(tdir + "/" + DefaultDirName)

	assert.Nil(t, Config())

	err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, Config())
	assert.Equal(t, filepath.Join(tdir, DefaultDirName, configFilename), ConfigFile())
	assert.FileExists(t, ConfigFile())
	assert.NotNil(t, Config())
}

/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package store

import (
	"os"
	"path/filepath"
)

var dir = DefaultDirName
var configFilepath string

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, DirPerm); err != nil {
			return err
		}
	}
	return nil
}

func defaultConfigFilepath() string {
	return filepath.Join(dir, configFilename)
}

// SetDefaultDir sets the default store working directory (eg. /tmp/.cas)
func SetDefaultDir() error {
	// Find home directory
	tmpDir := os.TempDir()
	cas := DefaultDirName

	SetDir(filepath.Join(tmpDir, cas))
	return nil
}

// SetDir sets the store working directory (eg. /tmp/.cas)
func SetDir(p string) {
	dir = p
}

// ConfigFile returns the config file path
func ConfigFile() string {
	if configFilepath == "" {
		return defaultConfigFilepath()
	}
	return configFilepath
}

// SetConfigFile sets the config file path (e.g. /tmp/.cas/config.json)
func SetConfigFile(filepath string) {
	configFilepath = filepath
}

// CurrentConfigFilePath returns the current config file path (e.g. /tmp/.cas/config.json)
func CurrentConfigFilePath() string {
	return dir
}

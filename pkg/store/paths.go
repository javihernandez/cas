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

	"github.com/mitchellh/go-homedir"
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

// SetDefaultDir sets the default store working directory
func SetDefaultDir() error {
	// Find home directory
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	SetDir(filepath.Join(home, DefaultDirName))
	return nil
}

// SetDir sets the store working directory
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

// SetConfigFile sets the config file path
func SetConfigFile(filepath string) {
	configFilepath = filepath
}

// CurrentConfigFilePath returns the current config file path
func CurrentConfigFilePath() string {
	return dir
}

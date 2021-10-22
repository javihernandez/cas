/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configSchemaVer uint = 3
)

// User holds user's configuration.
type User struct {
	Token    string `json:"token,omitempty"`
	KeyStore string `json:"keystore,omitempty"`
	LcCert   string `json:"lcCert,omitempty"`
}

// ConfigRoot holds root fields of the configuration file.
type ConfigRoot struct {
	SchemaVersion  uint           `json:"schemaVersion"`
	Users          []*User        `json:"users"`
	CurrentContext CurrentContext `json:"currentContext"`
}

type CurrentContext struct {
	LcHost          string `json:"LcHost,omitempty"`
	LcPort          string `json:"LcPort,omitempty"`
	LcCert          string `json:"LcCert,omitempty"`
	LcSkipTlsVerify bool   `json:"LcSkipTlsVerify,omitempty"`
	LcNoTls         bool   `json:"LcNoTls,omitempty"`
}

func (cc *CurrentContext) Clear() {
	cc.LcHost = ""
	cc.LcPort = ""
	cc.LcCert = ""
	cc.LcSkipTlsVerify = false
	cc.LcNoTls = false
}

var cfg *ConfigRoot
var v = viper.New()

// Config returns the global config instance
func Config() *ConfigRoot {
	return cfg
}

func setupConfigFile() string {
	cfgFile := ConfigFile()
	v.SetConfigFile(cfgFile)
	v.SetConfigPermissions(FilePerm)
	return cfgFile
}

// LoadConfig loads the global configuration from file
func LoadConfig() error {

	// Make default config
	c := ConfigRoot{
		SchemaVersion: configSchemaVer,
	}
	cfg = &c

	// Setup config file
	cfgFile := setupConfigFile()

	SetDir(filepath.Dir(cfgFile))

	// Ensure working dir
	if err := ensureDir(dir); err != nil {
		return err
	}

	// Create default file if it does not exist yet
	if ConfigFile() == defaultConfigFilepath() {
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			sErr := SaveConfig()
			return sErr
		}
	}

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&c); err != nil {
		oldFormat := ConfigRootV2{
			SchemaVersion: 2,
		}
		if err := v.Unmarshal(&oldFormat); err != nil {
			return errors.New("unable to parse config file")
		}
		fmt.Println("Upgrading config to new format. Old sessions will expire")
		c.Users = []*User{}
		c.SchemaVersion = 3
	}

	return SaveConfig()
}

// SaveConfig stores the current configuration to file
func SaveConfig() error {
	// Setup config file
	setupConfigFile()

	// Ensure working dir
	if err := ensureDir(dir); err != nil {
		return err
	}

	cfg.SchemaVersion = configSchemaVer
	v.Set("users", cfg.Users)
	v.Set("currentContext", cfg.CurrentContext)
	v.Set("schemaVersion", cfg.SchemaVersion)
	return v.WriteConfig()
}

// User returns an User from the global config matching the given email.
// User returns nil when an empty email is given or c is nil.
func (c *ConfigRoot) NewLcUser(host, port, lcCert string, lcSkipTlsVerify, lcNoTls bool) (u *CurrentContext) {
	defer func() {
		cfg.CurrentContext.Clear()
		cfg.CurrentContext.LcHost = host
		cfg.CurrentContext.LcPort = port
		cfg.CurrentContext.LcCert = lcCert
		cfg.CurrentContext.LcSkipTlsVerify = lcSkipTlsVerify
		cfg.CurrentContext.LcNoTls = lcNoTls
	}()

	return u
}

// ClearContext clean up all auth token for all users and set an empty context.
func (c *ConfigRoot) ClearContext() {
	if c == nil {
		return
	}
	for _, u := range c.Users {
		u.Token = ""
	}
	c.CurrentContext = CurrentContext{}
}

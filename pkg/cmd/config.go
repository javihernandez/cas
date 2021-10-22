/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/codenotary/cas/pkg/extractor"
	"github.com/codenotary/cas/pkg/extractor/docker"
	"github.com/codenotary/cas/pkg/extractor/file"
	"github.com/codenotary/cas/pkg/extractor/git"
	"github.com/codenotary/cas/pkg/extractor/wildcard"

	"github.com/codenotary/cas/pkg/store"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Register metadata extractors
	extractor.Register("", wildcard.Artifact)
	extractor.Register(file.Scheme, file.Artifact)
	extractor.Register(docker.Scheme, docker.Artifact)
	extractor.Register(docker.SchemePodman, docker.Artifact)
	extractor.Register(git.Scheme, git.Artifact)
	extractor.Register(wildcard.Scheme, wildcard.Artifact)

	// Load config
	if cfgFile != "" {
		store.SetConfigFile(cfgFile)
		if output, _ := rootCmd.PersistentFlags().GetString("output"); output == "" {
			fmt.Println("Using config file: ", store.ConfigFile())
		}
	}
	if err := store.LoadConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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
	"strings"

	"github.com/codenotary/cas/pkg/cmd/bom"
	"github.com/codenotary/cas/pkg/cmd/inspect"
	"github.com/codenotary/cas/pkg/cmd/internal/cli"
	"github.com/codenotary/cas/pkg/cmd/internal/types"
	"github.com/codenotary/cas/pkg/cmd/login"
	"github.com/codenotary/cas/pkg/cmd/logout"
	"github.com/codenotary/cas/pkg/cmd/sign"
	"github.com/codenotary/cas/pkg/cmd/verify"
	"github.com/codenotary/cas/pkg/cmd/list"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/codenotary/cas/pkg/store"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     meta.CasPrefix,
	Version: meta.Version(),
	Long:    ``,
}

// Root returns the root &cobra.Command
func Root() *cobra.Command {
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err error
	var cmd *cobra.Command
	var output string
	if cmd, err = rootCmd.ExecuteC(); err != nil {
		if !viper.IsSet("exit-code") {
			viper.Set("exit-code", 1)
		}
		output, _ = rootCmd.PersistentFlags().GetString("output")
		if output != "" && !cmd.SilenceErrors {
			cli.PrintError(output, types.NewError(err))
		}
	}

	exitCode := meta.CasDefaultExitCode
	if viper.IsSet("exit-code") {
		exitCode = viper.GetInt("exit-code")
	}
	os.Exit(exitCode)
}

func init() {

	// Read in environment variables that match
	viper.SetEnvPrefix(strings.ToUpper(meta.CasPrefix))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Set config files directory based on os.TempDir method ( Linux: /temp/.cas, Windows: c:\temp, c:\windows\temp )
	if err := store.SetDefaultDir(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Disable default behavior when started through explorer.exe
	cobra.MousetrapHelpText = ""

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "caspath", "", "config files (default is /tmp/.cas/config.json on linux, c:\\temp\\config.json or c:\\windows\\temp\\config.json on Windows)")
	rootCmd.PersistentFlags().StringP("output", "o", "", "output format, one of: --output=json|--output=''")
	rootCmd.PersistentFlags().BoolP("silent", "S", false, "silent mode, don't show progress spinner, but it will still output the result")
	rootCmd.PersistentFlags().BoolP("quit", "q", true, "if false, ask for confirmation before quitting")
	rootCmd.PersistentFlags().Bool("verbose", false, "if true, print additional information")

	rootCmd.PersistentFlags().MarkHidden("quit")

	// Root command flags
	rootCmd.Flags().BoolP("version", "v", false, "version for cas") // needed for -v shorthand

	// Verification group
	rootCmd.AddCommand(verify.NewCommand())
	rootCmd.AddCommand(inspect.NewCommand())

	// Signing group
	rootCmd.AddCommand(sign.NewCommand())
	rootCmd.AddCommand(sign.NewUntrustCommand())
	rootCmd.AddCommand(sign.NewUnsupportCommand())

	// User group
	rootCmd.AddCommand(login.NewCommand())
	rootCmd.AddCommand(logout.NewCommand())

	// BoM
	rootCmd.AddCommand(bom.NewCommand())

	// List
	rootCmd.AddCommand(list.NewCommand())
}

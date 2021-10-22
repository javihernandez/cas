/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package logout

import (
	"fmt"

	"github.com/codenotary/cas/pkg/meta"
	"github.com/fatih/color"

	"github.com/codenotary/cas/pkg/store"

	"github.com/spf13/cobra"
)

// NewCommand returns the cobra command for `cas logout`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout the current user",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			if store.Config() == nil || store.Config().CurrentContext.LcHost == "" {
				color.Set(meta.StyleWarning())
				fmt.Println("No logged-in user.")
				color.Unset()
				return nil
			}
			if err := Execute(); err != nil {
				return err
			}
			if output == "" {
				color.Set(meta.StyleSuccess())
				fmt.Println("Logout successful.")
				color.Unset()
			}
			return nil
		},
		Args: cobra.NoArgs,
	}

	return cmd
}

// Execute logout action for Immutable Ledger
func Execute() error {
	store.Config().ClearContext()
	if err := store.SaveConfig(); err != nil {
		return err
	}
	return nil
}

/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package sign

import (
	"github.com/codenotary/cas/pkg/meta"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewUnsupportCommand returns the cobra command for `cas unsupport`
func NewUnsupportCommand() *cobra.Command {
	cmd := makeCommand()
	cmd.Use = "unsupport"
	cmd.Aliases = []string{"us"}
	cmd.Short = "Unsupport an asset"
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runSignWithState(cmd, args, meta.StatusUnsupported)
	}
	cmd.Long = `
Change an asset's status so it is equal to UNSUPPORTED.

Unsupport command calculates the SHA-256 hash of a digital asset
(file, directory, container's image).
The hash (not the asset) and the desired status of UNSUPPORTED are then
cryptographically signed by the signer's secret (private key).
Next, these signed objects are sent to the CAS where the signer’s
trust level and a timestamp are added.
When complete, a new CAS entry is created that binds the asset’s
signed hash, signed status, level, and timestamp together.

Note that your assets will not be uploaded. They will be processed locally.

Assets are referenced by passed ARG(s) with unsupport command only accepting
1 ARG at a time.

` + helpMsgFooter

	return cmd
}

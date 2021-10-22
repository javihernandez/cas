/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package sign

import (
	"fmt"

	"github.com/spf13/cobra"
)

func noArgsWhenHashOrPipe(cmd *cobra.Command, args []string) error {
	if hash, _ := cmd.Flags().GetString("hash"); hash != "" {
		if len(args) > 0 {
			return fmt.Errorf("cannot use ARG(s) with --hash")
		}
		return nil
	}
	if pipeMode() {
		return nil
	}
	return cobra.ExactArgs(1)(cmd, args)
}

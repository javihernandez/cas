/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package inspect

import (
	"fmt"
	"strings"

	"github.com/codenotary/cas/pkg/cmd/internal/cli"
	"github.com/codenotary/cas/pkg/cmd/internal/types"

	"github.com/codenotary/cas/pkg/meta"
	"github.com/spf13/viper"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/extractor"
	"github.com/codenotary/cas/pkg/store"
	"github.com/spf13/cobra"
)

// NewCommand returns the cobra command for `cas inspect`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect",
		Aliases: []string{"i"},
		Short:   "Returns the asset history with low-level information",
		Long: `
Returns the asset history with low-level information

Environment variables:
CAS_HOST=
CAS_PORT=
CAS_CERT=
CAS_SKIP_TLS_VERIFY=false
CAS_NO_TLS=false
CAS_API_KEY=
CAS_LEDGER=
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: runInspect,
		Args: func(cmd *cobra.Command, args []string) error {
			if hash, _ := cmd.Flags().GetString("hash"); hash != "" {
				if len(args) > 0 {
					return fmt.Errorf("cannot use ARG(s) with --hash")
				}
				return nil
			}

			first, _ := cmd.Flags().GetUint64("first")
			last, _ := cmd.Flags().GetUint64("last")
			start, _ := cmd.Flags().GetString("start")
			end, _ := cmd.Flags().GetString("end")

			if (first > 0 || last > 0 || start != "" || end != "") &&
				store.Config().CurrentContext.LcHost == "" {
				return fmt.Errorf("time range filter is available only in Ledger Compliance environment")
			}

			if first > 0 && last > 0 {
				return fmt.Errorf("--first and --last are mutual exclusive")
			}
			return cobra.MinimumNArgs(1)(cmd, args)
		},
		Example: `
cas inspect document.pdf --last 1
cas inspect document.pdf --first 1
cas inspect document.pdf --start 2020/10/28-08:00:00 --end 2020/10/28-17:00:00 --first 10
cas inspect document.pdf --signerID CygBE_zb8XnprkkO6ncIrbbwYoUq5T1zfyEF6DhqcAI= --start 2020/10/28-16:00:00 --end 2020/10/28-17:10:00 --last 3
`,
	}

	cmd.SetUsageTemplate(
		strings.Replace(cmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}} ARG", 1),
	)

	cmd.Flags().String("hash", "", "specify a hash to inspect, if set no ARG can be used")
	cmd.Flags().Bool("extract-only", false, "if set, print only locally extracted info")
	// ledger compliance flags
	cmd.Flags().String("host", "", meta.CasHostFlagDesc)
	cmd.Flags().String("port", "", meta.CasPortFlagDesc) // set to default port in GetOrCreateLcUser(), if not available from context
	cmd.Flags().String("cert", "", meta.CasCertPathDesc)
	cmd.Flags().Bool("skip-tls-verify", false, meta.CasSkipTlsVerifyDesc)
	cmd.Flags().Bool("no-tls", false, meta.CasNoTlsDesc)
	cmd.Flags().String("api-key", "", meta.CasApiKeyDesc)
	cmd.Flags().String("ledger", "", meta.CasLedgerDesc)

	cmd.Flags().String("signerID", "", "specify a signerID to refine inspection result on ledger compliance")

	cmd.Flags().Uint64("first", 0, "set the limit for the first elements filter. MAX 10")
	cmd.Flags().Uint64("last", 0, "set the limit for the last elements filter. MAX 10")

	cmd.Flags().String("start", "", "set the start of date and time range filter. Example 2020/10/28-16:00:00")
	cmd.Flags().String("end", "", "set the end of date and time range filter. Example 2020/10/28-16:00:00")

	return cmd
}

func runInspect(cmd *cobra.Command, args []string) error {

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	hash, err := cmd.Flags().GetString("hash")
	if err != nil {
		return err
	}
	hash = strings.ToLower(hash)

	extractOnly, err := cmd.Flags().GetBool("extract-only")
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true

	if hash == "" {
		if len(args) < 1 {
			return fmt.Errorf("no argument")
		}
		if hash, err = extractInfo(args[0], output); err != nil {
			return err
		}
		if output == "" {
			fmt.Print("\n\n")
		}
	}

	if extractOnly {
		return nil
	}

	signerID, err := cmd.Flags().GetString("signerID")
	if err != nil {
		return err
	}

	lcHost := viper.GetString("host")
	lcPort := viper.GetString("port")
	lcCert := viper.GetString("cert")
	lcApiKey := viper.GetString("api-key")
	lcLedger := viper.GetString("ledger")

	lcUser, err := api.GetOrCreateLcUser(lcApiKey, lcLedger, lcHost, lcPort, lcCert, viper.IsSet("skip-tls-verify"), viper.GetBool("skip-tls-verify"), viper.IsSet("no-tls"), viper.GetBool("no-tls"), nil, false)
	if err != nil {
		return err
	}

	first, err := cmd.Flags().GetUint64("first")
	if err != nil {
		return err
	}
	if first > 10 {
		return fmt.Errorf("only first 10 items are allowed when using --first flag")
	}
	last, err := cmd.Flags().GetUint64("last")
	if err != nil {
		return err
	}
	if last > 10 {
		return fmt.Errorf("only last 10 items are allowed when using --last flag")
	}
	start, err := cmd.Flags().GetString("start")
	if err != nil {
		return err
	}
	end, err := cmd.Flags().GetString("end")
	if err != nil {
		return err
	}

	if first == 0 && last == 0 {
		last = 10
		fmt.Printf("no filter is specified. At maximum last 10 items will be returned\n")
	}
	return lcInspect(hash, signerID, lcUser, first, last, start, end, output)
}

func extractInfo(arg string, output string) (hash string, err error) {
	a, err := extractor.Extract([]string{arg})
	if err != nil {
		return "", err
	}
	if len(a) == 0 {
		return "", fmt.Errorf("unable to process the input asset provided: %s", arg)
	}
	if len(a) == 1 {
		hash = a[0].Hash
	}
	if len(a) > 1 {
		return "", fmt.Errorf("info extraction on multiple items is not yet supported")
	}
	if output == "" {
		fmt.Printf("Extracted info from: %s\n\n", arg)
	}
	cli.Print(output, types.NewResult(a[0], nil))
	return
}

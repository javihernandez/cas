/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package login

import (
	"fmt"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/signature"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/codenotary/cas/pkg/meta"
)

// NewCommand returns the cobra command for `cas login`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			return nil
		},
		Use:   "login",
		Short: "Log in to Community Attestation Service",
		Long: `Log in to Community Attestation Service.

Environment variables:
CAS_HOST=
CAS_PORT=
CAS_CERT=
CAS_SKIP_TLS_VERIFY=false
CAS_NO_TLS=false
CAS_API_KEY=
CAS_LEDGER=
`,
		Example: `  # Codenotary.io login:
  ./cas login
  # On-premise service login:
  ./cas login --port 33443 --host lc.vchain.us --cert mycert.pem
  ./cas login --port 3324 --host 127.0.0.1 --no-tls
  ./cas login --port 443 --host lc.vchain.us --skip-tls-verify`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			lcHost := api.DefaultHost
			if len(args) > 0 {
				lcHost = args[0]
			} else if viper.IsSet("host") {
				lcHost = viper.GetString("host")
			}

			lcPort := viper.GetString("port")
			lcCert := viper.GetString("cert")
			lcApiKey := viper.GetString("api-key")
			lcLedger := viper.GetString("ledger")

			signingPubKey, skipLocalPubKeyComp, err := signature.PrepareSignatureParams(
				viper.GetString("signing-pub-key"),
				viper.GetString("signing-pub-key-file"))
			if err != nil {
				return err
			}
			enforceSignatureVerify := viper.GetBool("enforce-signature-verify")

			err = ExecuteLC(lcHost, lcPort, lcCert, lcApiKey, lcLedger, viper.IsSet("skip-tls-verify"), viper.GetBool("skip-tls-verify"), viper.IsSet("no-tls"), viper.GetBool("no-tls"), signingPubKey, skipLocalPubKeyComp, enforceSignatureVerify)
			if err != nil {
				return err
			}
			if output == "" {
				color.Set(meta.StyleSuccess())
				fmt.Println("Login successful.")
				color.Unset()
			}
			return nil
		},
		Args: cobra.MaximumNArgs(2),
	}
	cmd.Flags().String("host", "", meta.CasHostFlagDesc)
	cmd.Flags().String("port", "", meta.CasPortFlagDesc)
	cmd.Flags().String("cert", "", meta.CasCertPathDesc)
	cmd.Flags().Bool("skip-tls-verify", false, meta.CasSkipTlsVerifyDesc)
	cmd.Flags().Bool("no-tls", false, meta.CasNoTlsDesc)
	cmd.Flags().String("api-key", "", meta.CasApiKeyDesc)
	cmd.Flags().String("ledger", "", meta.CasLedgerDesc)
	cmd.Flags().String("signing-pub-key-file", "", meta.CasSigningPubKeyFileNameDesc)
	cmd.Flags().String("signing-pub-key", "", meta.CasSigningPubKeyDesc)
	cmd.Flags().Bool("enforce-signature-verify", false, meta.CasEnforceSignatureVerifyDesc)

	return cmd
}

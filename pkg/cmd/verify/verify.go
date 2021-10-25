/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package verify

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	caserr "github.com/codenotary/cas/internal/errors"
	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/extractor"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/codenotary/cas/pkg/signature"
	"github.com/vchain-us/ledger-compliance-go/schema"
)

type pkg struct {
	Name   string
	Hash   string
	Kind   string
	Md     md `json:"metadata"`
	Status int
}

type md struct {
	Version  string `json:"version,omitempty"`
	HashType string `json:"hashType"`
}

var (
	keyRegExp = regexp.MustCompile("0x[0-9a-z]{40}")
)

func getSignerIDs() []string {
	ids := viper.GetStringSlice("signerID")
	if len(ids) > 0 {
		for i := range ids {
			if strings.Contains(ids[i], "@") {
				// signer is an e-mail - encode it
				ids[i] = base64.StdEncoding.EncodeToString([]byte(ids[i]))
			}
		}
		return ids
	}
	return viper.GetStringSlice("key")
}

// NewCommand returns the cobra command for `cas verify`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "authenticate",
		Example: "  cas authenticate /bin/cas",
		Aliases: []string{"a", "verify", "v"},
		Short:   "Authenticate assets against CAS",
		Long: `
Authenticate assets against the CAS.

Authentication is the process of matching the hash of a local asset to
a hash on CAS.
If matched, the returned result (the authentication) is the CAS
metadata that’s bound to the matching hash.
Otherwise, the returned result status equals UNKNOWN.

Note that your assets will not be uploaded but processed locally.

The exit code will be 0 only if all assets' statuses are equal to TRUSTED.
Otherwise,
	Status Untrusted:     1
	Status Unknown:       2
	Status Unsupported:   3
	Status ApikeyRevoked: 4

Assets are referenced by the passed ARG(s), with authentication accepting
1 or more ARG(s) at a time. Multiple assets can be authenticated at the
same time while passing them within ARG(s).

ARG must be one of:
  <file>
  file://<file>
  git://<repository>
  docker://<image>
  podman://<image>
Environment variables:
CAS_HOST=
CAS_PORT=
CAS_CERT=
CAS_SKIP_TLS_VERIFY=false
CAS_NO_TLS=false
CAS_API_KEY=
CAS_LEDGER=
CAS_SIGNING_PUB_KEY_FILE=
CAS_SIGNING_PUB_KEY=
CAS_ENFORCE_SIGNATURE_VERIFY=
`,
		RunE: runVerify,
		PreRun: func(cmd *cobra.Command, args []string) {
			// Bind to all flags to env vars (after flags were parsed),
			// but only ones retrivied by using viper will be used.
			viper.BindPFlags(cmd.Flags())
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if org := viper.GetString("org"); org != "" {
				if keys := getSignerIDs(); len(keys) > 0 {
					return fmt.Errorf("cannot use both --org and SignerID(s)")
				}
			}

			if hash, _ := cmd.Flags().GetString("hash"); hash != "" {
				if len(args) > 0 {
					return fmt.Errorf("cannot use ARG(s) with --hash")
				}
				return nil
			}

			return cobra.MinimumNArgs(1)(cmd, args)
		},
	}

	cmd.SetUsageTemplate(
		strings.Replace(cmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}} ARG(s)", 1),
	)

	cmd.Flags().StringSliceP("signerID", "s", nil, "accept only authentications matching the passed SignerID(s)")
	cmd.Flags().StringSliceP("key", "k", nil, "")
	cmd.Flags().MarkDeprecated("key", "please use --signerID instead")
	cmd.Flags().String("hash", "", "specify a hash to authenticate, if set no ARG(s) can be used")
	cmd.Flags().Int("exit-code", meta.CasDefaultExitCode, meta.CasExitCode)
	cmd.Flags().String("host", "", meta.CasHostFlagDesc)
	cmd.Flags().String("port", "", meta.CasPortFlagDesc) // set to default port in GetOrCreateLcUser(), if not available from context
	cmd.Flags().String("cert", "", meta.CasCertPathDesc)
	cmd.Flags().Bool("skip-tls-verify", false, meta.CasSkipTlsVerifyDesc)
	cmd.Flags().Bool("no-tls", false, meta.CasNoTlsDesc)
	cmd.Flags().String("api-key", "", meta.CasApiKeyDesc)
	cmd.Flags().String("ledger", "", meta.CasLedgerDesc)
	cmd.Flags().String("uid", "", meta.CasUidDesc)
	cmd.Flags().Bool("bom", false, "link asset to its dependencies from BoM")
	cmd.Flags().String("bom-trust-level", "trusted", "min trust level: untrusted (unt) / unsupported (uns) / unknown (unk) / trusted (t)")
	cmd.Flags().Float64("bom-max-unsupported", 0, "max number (in %) of unsupported dependencies")
	cmd.Flags().Uint("bom-batch-size", 10, "By default BoM dependencies are authenticated/notarized in batches of up to 10 dependencies each. Use this flag to set a different batch size. A value of 0 will disable batching (all dependencies will be authenticated/notarized at once).")
	// BoM output options
	cmd.Flags().String("bom-spdx", "", "name of the file to output BoM in SPDX format")
	cmd.Flags().String("bom-cyclonedx-json", "", "name of the file to output BoM in CycloneDX JSON format")
	cmd.Flags().String("bom-cyclonedx-xml", "", "name of the file to output BoM in CycloneDX XML format")

	cmd.Flags().String("signing-pub-key-file", "", meta.CasSigningPubKeyFileNameDesc)
	cmd.Flags().String("signing-pub-key", "", meta.CasSigningPubKeyDesc)
	cmd.Flags().Bool("enforce-signature-verify", false, meta.CasEnforceSignatureVerifyDesc)
	cmd.Flags().MarkHidden("raw-diff")

	return cmd
}

func runVerify(cmd *cobra.Command, args []string) error {
	hashes := make([]string, 0)
	hash, err := cmd.Flags().GetString("hash")
	if err != nil {
		return err
	}
	if hash != "" {
		hashes = append(hashes, strings.ToLower(hash))
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	cmd.SilenceUsage = true

	lcHost := viper.GetString("host")
	lcPort := viper.GetString("port")
	lcCert := viper.GetString("cert")
	lcApiKey := viper.GetString("api-key")
	lcLedger := viper.GetString("ledger")
	lcUid := viper.GetString("uid")
	lcVerbose := viper.GetBool("verbose")

	signingPubKey, skipLocalPubKeyComp, err := signature.PrepareSignatureParams(
		viper.GetString("signing-pub-key"),
		viper.GetString("signing-pub-key-file"))
	if err != nil {
		return err
	}
	enforceSignatureVerify := viper.GetBool("enforce-signature-verify")

	publicAuth := false
	if lcApiKey == "" {
		publicAuth = true
	}

	var signerID string
	signerIDs := getSignerIDs()
	if len(signerIDs) > 0 {
		signerID = signerIDs[0]
	}
	if lcApiKey == "" && signerID == "" {
		return caserr.ErrPubAuthNoSignerID
	}
	lcUser, err := api.GetOrCreateLcUser(lcApiKey, lcLedger, lcHost, lcPort, lcCert, viper.IsSet("skip-tls-verify"), viper.GetBool("skip-tls-verify"), viper.IsSet("no-tls"), viper.GetBool("no-tls"), signingPubKey, publicAuth)
	if err != nil {
		return err
	}

	if !skipLocalPubKeyComp {
		err = lcUser.CheckConnectionPublicKey(enforceSignatureVerify)
		if err != nil {
			return err
		}
	}

	// any set 'bom-xxx' option, except 'bom-what-includes', implies BoM
	bomFlag := viper.GetBool("bom") ||
		viper.IsSet("bom-trust-level") ||
		viper.IsSet("bom-max-unsupported") ||
		viper.IsSet("bom-spdx") ||
		viper.IsSet("bom-cyclonedx-json") ||
		viper.IsSet("bom-cyclonedx-xml") ||
		viper.IsSet("bom-batch-size")

	if bomFlag {
		err := lcUser.RequireFeatOrErr(schema.FeatBoM)
		if err != nil {
			return err
		}
	}

	var bomArtifact artifact.Artifact
	if bomFlag {
		if len(hashes)+len(args) > 1 {
			return fmt.Errorf("asset selection criteria match several assets - BoM can be processed only for single asset")
		}
		if len(hashes)+len(args) < 1 {
			return fmt.Errorf("asset selection criteria don't match any assets - BoM cannot be processed")
		}

		if len(hashes) > 0 {
			bomArtifact, err = processBOM(lcUser, signerID, output, hashes[0], "")
		} else {
			bomArtifact, err = processBOM(lcUser, signerID, output, "", args[0])
		}
		// in case of diff don't stop if some dependencies have insufficient trust level
		if err != nil && err != ErrInsufficientTrustLevel {
			return err
		}
	}

	if len(hashes) > 0 {
		for _, hash := range hashes {
			err = lcVerify(cmd, &api.Artifact{Hash: hash}, lcUser, signerID, lcUid, lcVerbose, output)
			if err != nil {
				return err
			}
		}
	} else {
		artifacts, err := extractor.Extract([]string{args[0]})
		if err != nil {
			return err
		}
		for _, a := range artifacts {
			if bomArtifact != nil {
				a.Deps = DepsToPackageDetails(bomArtifact.Dependencies())
			}
			err := lcVerify(cmd, a, lcUser, signerID, lcUid, lcVerbose, output)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

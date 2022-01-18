/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package sign

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codenotary/cas/pkg/signature"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/bom"
	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/bom/docker"
	"github.com/codenotary/cas/pkg/cicontext"
	"github.com/codenotary/cas/pkg/cmd/verify"
	"github.com/codenotary/cas/pkg/extractor"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/codenotary/cas/pkg/uri"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vchain-us/ledger-compliance-go/schema"
)

const helpMsgFooter = `
ARG must be one of:
  wildcard
  file
  directory
  file://<file>
  git://<repository>
  docker://<image>
  podman://<image>
  wildcard://"*"
`

// NewCommand returns the cobra command for `cas sign`
func NewCommand() *cobra.Command {
	cmd := makeCommand()
	cmd.Flags().Bool("bom", false, "auto-notarize asset dependencies and link dependencies to the asset")
	cmd.Flags().String("bom-signerID", "", "signerID to use for authenticating dependencies")
	cmd.Flags().Uint("bom-batch-size", 10, "By default BOM dependencies are authenticated/notarized in batches of up to 10 dependencies each. Use this flag to set a different batch size. A value of 0 will disable batching (all dependencies will be authenticated/notarized at once).")
	// BOM output options
	cmd.Flags().String("bom-spdx", "", "name of the file to output BOM in SPDX format")
	cmd.Flags().String("bom-cdx-json", "", "name of the file to output BOM in CycloneDX JSON format")
	cmd.Flags().String("bom-cdx-xml", "", "name of the file to output BOM in CycloneDX XML format")
	return cmd
}

func makeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "notarize",
		Aliases: []string{"n", "sign", "s"},
		Short:   "Notarize an asset onto Community Attestation Service",
		Long: `
Notarize an asset onto the CAS.

Notarization calculates the SHA-256 hash of a digital asset
(file, directory, container's image).
The hash (not the asset) and the desired status of TRUSTED are then
cryptographically signed by the signer's secret (private key).
Next, these signed objects are sent to the CAS where the signer’s
trust level and a timestamp are added.
When complete, a new entry is created that binds the asset’s
signed hash, signed status, level, and timestamp together.

Note that your assets will not be uploaded. They will be processed locally.

Assets are referenced by passed ARG with notarization only accepting
1 ARG at a time.

Pipe mode:
If '-' is provided (echo my-file | cas n -) stdin is read and parsed. Only pipe ARGs are processed.

Environment variables:
CAS_HOST=
CAS_PORT=
CAS_CERT=
CAS_SKIP_TLS_VERIFY=false
CAS_NO_TLS=false
CAS_API_KEY=
CAS_SIGNING_PUB_KEY_FILE=
CAS_SIGNING_PUB_KEY=
CAS_ENFORCE_SIGNATURE_VERIFY=
` + helpMsgFooter,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if pipeMode() && len(args) > 0 && args[0] == "-" {
				args = make([]string, 0)
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Split(bufio.ScanWords)
				for scanner.Scan() {
					token := scanner.Bytes()
					args = append(args, string(token))
				}
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("error parsing stdin input: %s", err)
				}
			}
			return runSignWithState(cmd, args, meta.StatusTrusted)
		},
		Args: noArgsWhenHashOrPipe,
		Example: `cas notarize my-file
echo my-file | cas n -`,
	}

	cmd.Flags().VarP(make(mapOpts), "attr", "a", "add user defined attributes (repeat --attr for multiple entries)")
	cmd.Flags().Bool("ci-attr", false, meta.CasCIAttribDesc)
	cmd.Flags().StringP("name", "n", "", "set the asset name")
	cmd.Flags().String("hash", "", "specify the hash instead of using an asset, if set no ARG(s) can be used")
	cmd.Flags().String("host", "", meta.CasHostFlagDesc)
	cmd.Flags().String("port", "", meta.CasPortFlagDesc) // set to default port in GetOrCreateLcUser(), if not available from context
	cmd.Flags().String("cert", "", meta.CasCertPathDesc)
	cmd.Flags().Bool("skip-tls-verify", false, meta.CasSkipTlsVerifyDesc)
	cmd.Flags().Bool("no-tls", false, meta.CasNoTlsDesc)
	cmd.Flags().String("api-key", "", meta.CasApiKeyDesc)

	cmd.SetUsageTemplate(
		strings.Replace(cmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}} ARG", 1),
	)
	cmd.Flags().String("signing-pub-key-file", "", meta.CasSigningPubKeyFileNameDesc)
	cmd.Flags().String("signing-pub-key", "", meta.CasSigningPubKeyDesc)
	cmd.Flags().Bool("enforce-signature-verify", false, meta.CasEnforceSignatureVerifyDesc)

	return cmd
}

func runSignWithState(cmd *cobra.Command, args []string, state meta.Status) error {
	// default extractors options
	extractorOptions := []extractor.Option{}

	var hash string
	if hashFlag := cmd.Flags().Lookup("hash"); hashFlag != nil {
		var err error
		hash, err = cmd.Flags().GetString("hash")
		if err != nil {
			return err
		}
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	metadata := cmd.Flags().Lookup("attr").Value.(mapOpts).StringToInterface()

	// @todo use dependency injection
	cs := cicontext.NewContextSaver()

	if viper.GetBool("ci-attr") {
		cicontext.ExtendMetadata(metadata, cs.GetCIContextMetadata())
	}

	cmd.SilenceUsage = true

	lcHost := viper.GetString("host")
	lcPort := viper.GetString("port")
	lcCert := viper.GetString("cert")
	lcApiKey := viper.GetString("api-key")

	lcVerbose := viper.GetBool("verbose")

	signingPubKey, skipLocalPubKeyComp, err := signature.PrepareSignatureParams(
		viper.GetString("signing-pub-key"),
		viper.GetString("signing-pub-key-file"))
	if err != nil {
		return err
	}
	enforceSignatureVerify := viper.GetBool("enforce-signature-verify")

	lcUser, err := api.GetOrCreateLcUser(lcApiKey, "", lcHost, lcPort, lcCert, viper.IsSet("skip-tls-verify"), viper.GetBool("skip-tls-verify"), viper.IsSet("no-tls"), viper.GetBool("no-tls"), signingPubKey, false)
	if err != nil {
		return err
	}

	// any set `--bom-xxx` option implies bom mode
	bomFlag := viper.GetBool("bom") ||
		viper.IsSet("bom-signerID") ||
		viper.IsSet("bom-spdx") ||
		viper.IsSet("bom-cdx-json") ||
		viper.IsSet("bom-cdx-xml") ||
		viper.IsSet("bom-batch-size")

	artifacts := make([]*api.Artifact, 0, 1)

	for attr := range metadata {
		if attr == "allowdownload" {
			err := lcUser.RequireFeatOrErr(schema.FeatAllowDownload)
			if err != nil {
				return err
			}
			break
		}
	}

	if !skipLocalPubKeyComp {
		err = lcUser.CheckConnectionPublicKey(enforceSignatureVerify)
		if err != nil {
			return err
		}
	}

	var bomLinks []*schema.VCNDependency

	if bomFlag {
		err := lcUser.RequireFeatOrErr(schema.FeatBoM)
		if err != nil {
			return err
		}
	}
	outputOpts := artifact.Progress
	if viper.GetBool("silent") || output != "" {
		outputOpts = artifact.Silent
	}

	var bomArtifact artifact.Artifact
	if bomFlag {
		// if bom-file specified, use BOM data from file, otherwise resolve dependencies
		if len(args) != 1 {
			return fmt.Errorf("--bom option can be used only with single asset")
		}
		path := args[0]
		u, err := uri.Parse(path)
		if err != nil {
			return err
		}
		if _, ok := bom.BomSchemes[u.Scheme]; !ok {
			return fmt.Errorf("unsupported URI %s for --bom option", path)
		}
		if u.Scheme != "" {
			path = strings.TrimPrefix(u.Opaque, "//")
		}
		if u.Scheme == "docker" {
			bomArtifact, err = docker.New(path)
			if err != nil {
				return err
			}
		} else {
			path, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			bomArtifact = bom.New(path)
		}
		if bomArtifact == nil {
			return fmt.Errorf("unsupported asset format/language")
		}

		if outputOpts != artifact.Silent {
			fmt.Printf("Resolving dependencies...\n")
		}
		deps, err := bomArtifact.ResolveDependencies(outputOpts)
		if err != nil {
			return fmt.Errorf("cannot get dependencies: %w", err)
		}

		bomBatchSize := int(viper.GetUint("bom-batch-size"))

		bomLinks, err = notarizeDeps(lcUser, deps, outputOpts, bomArtifact.Type(), bomBatchSize)
		if err != nil {
			return err
		}

		err = bom.Output(bomArtifact) // process all possible BOM output options
		if err != nil {
			// show warning, but not error, because authentication finished
			fmt.Println(err)
		}
		if outputOpts != artifact.Silent {
			artifact.Display(bomArtifact, artifact.ColNameVersion|artifact.ColHash|artifact.ColTrustLevel)
		}
	}

	// notarize the asset if not instructed otherwise
	if hash != "" {
		hash = strings.ToLower(hash)
		// Load existing artifact, if any, otherwise use an empty artifact
		if ar, _, err := lcUser.LoadArtifact(hash, "", "", 0, nil); err == nil && ar != nil {
			artifacts = append(artifacts, &api.Artifact{
				Kind:        ar.Kind,
				Name:        ar.Name,
				Hash:        ar.Hash,
				Size:        ar.Size,
				ContentType: ar.ContentType,
				Metadata:    ar.Metadata,
			})
		} else {
			if name == "" {
				return fmt.Errorf("please set an asset name, by using --name")
			}
			artifacts = append(artifacts, &api.Artifact{Hash: hash})
		}
	} else {
		artifacts, err = extractor.Extract(args, extractorOptions...)
		if err != nil {
			return err
		}
	}
	if bomArtifact != nil {
		artifacts[0].Deps = verify.DepsToPackageDetails(bomArtifact.Dependencies())
	}
	err = LcSign(lcUser, artifacts, state, output, name, metadata, lcVerbose, bomLinks)
	if err != nil {
		return err
	}

	return nil
}

func pipeMode() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func notarizeDeps(lcUser *api.LcUser, deps []artifact.Dependency, outputOpts artifact.OutputOptions, artType string, batchSize int) ([]*schema.VCNDependency, error) {
	if outputOpts != artifact.Silent {
		fmt.Printf("Authenticating dependencies...\n")
	}

	signerID := viper.GetString("bom-signerID")
	if signerID == "" {
		signerID = api.GetSignerIDByApiKey(lcUser.Client.ApiKey)
	}

	var bar *progressbar.ProgressBar
	if len(deps) > 1 && outputOpts == artifact.Progress {
		bar = progressbar.Default(int64(len(deps)))
	}

	progressCallback := func(processedDeps []artifact.Dependency) {
		if bar != nil {
			bar.Add(len(processedDeps))
		}
	}

	errs, err := artifact.AuthenticateDependencies(lcUser, signerID, deps, batchSize, progressCallback)
	if err != nil {
		return nil, fmt.Errorf("error authenticating dependencies: %w", err)
	}

	var msgs []string
	var depsToNotarize []*artifact.Dependency
	var kinds []string

	for i := range deps { // Authenticate mutates the dependency, so use the index
		if errs[i] != nil {
			return nil, fmt.Errorf("cannot authenticate %s@%s dependency: %w",
				deps[i].Name, deps[i].Version, errs[i])
		}
		if deps[i].TrustLevel < artifact.Unknown {
			msgs = append(msgs, fmt.Sprintf("Dependency %s@%s trust level is %s",
				deps[i].Name, deps[i].Version, artifact.TrustLevelName(deps[i].TrustLevel)))
		}
		if deps[i].TrustLevel < artifact.Trusted {
			depsToNotarize = append(depsToNotarize, &deps[i])
			kinds = append(kinds, artType)
		}
	}

	if len(msgs) > 0 {
		for _, msg := range msgs {
			fmt.Println(msg)
		}
		return nil, fmt.Errorf("some dependencies have insufficient trust level and cannot be automatically notarized")
	}

	// notarize only the dependencies first to make sure all needed keys are present in DB before
	// adding key references to the index
	if len(depsToNotarize) > 0 {
		var bar *progressbar.ProgressBar
		if outputOpts != artifact.Silent {
			ds := "dependencies"
			if len(depsToNotarize) == 1 {
				ds = "dependency"
			}
			fmt.Printf("Notarizing %d %s ...\n", len(depsToNotarize), ds)
			if outputOpts == artifact.Progress {
				bar = progressbar.Default(int64(len(depsToNotarize)))
			}
		}

		progressCallbackN := func(processedDeps []*artifact.Dependency) {
			if bar != nil {
				bar.Add(len(processedDeps))
			}
		}

		err = artifact.NotarizeDependencies(lcUser, kinds, depsToNotarize, batchSize, progressCallbackN)
		if err != nil {
			return nil, fmt.Errorf("error notarizing dependencies: %w", err)
		}
	} else {
		fmt.Printf("No dependencies require notarization\n")
	}

	bom := make([]*schema.VCNDependency, 0, len(deps))
	for i := range deps {
		// add dep key to BOM list for attaching
		depType := schema.VCNDependency_Direct
		if deps[i].Type == artifact.DepTransient {
			depType = schema.VCNDependency_Indirect
		}
		bom = append(bom, &schema.VCNDependency{
			Hash: deps[i].Hash,
			Type: depType,
		})
	}

	return bom, nil
}

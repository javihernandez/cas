package sign

import (
	"fmt"

	"github.com/caarlos0/spin"
	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/cmd/internal/cli"
	"github.com/codenotary/cas/pkg/cmd/internal/types"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/vchain-us/ledger-compliance-go/schema"
)

// LcSign ...
func LcSign(u *api.LcUser, artifacts []*api.Artifact, state meta.Status, output string, name string, metadata map[string]interface{}, verbose bool, bom []*schema.VCNDependency) error {
	if output == "" {
		color.Set(meta.StyleAffordance())
		fmt.Print("Your assets will not be uploaded. They will be processed locally.")
		color.Unset()
		fmt.Println()
		fmt.Println()
	}

	s := spin.New("%s Notarization in progress...")
	s.Set(spin.Spin1)

	var bar *progressbar.ProgressBar
	lenArtifacts := len(artifacts)
	if lenArtifacts > 1 && output == "" {
		bar = progressbar.Default(int64(lenArtifacts))
	}

	// Override the asset's name, if provided by --name
	if len(artifacts) == 1 && name != "" {
		artifacts[0].Name = name
	}

	for _, a := range artifacts {
		// Copy user provided custom attributes
		a.Metadata.SetValues(metadata)

		// @todo mmeloni use verified sign
		tx, err := u.Sign(
			*a,
			api.LcSignWithStatus(state),
			api.LcSignWithBom(bom),
		)
		if err != nil {
			if err == api.ErrNotVerified {
				color.Set(meta.StyleError())
				fmt.Println("the ledger is compromised. Please contact the Community Attestation Service administrators")
				color.Unset()
				fmt.Println()
				return nil
			}
			return err
		}

		if err != nil {
			return cli.PrintWarning(output, err.Error())
		}
		if output == "" && lenArtifacts == 0 {
			fmt.Println()
		}

		artifact, verified, err := u.LoadArtifact(a.Hash, "", "", tx, nil)
		if err != nil {
			if err == api.ErrNotVerified {
				color.Set(meta.StyleError())
				fmt.Println("the ledger is compromised. Please contact the Community Attestation Service administrators")
				color.Unset()
				fmt.Println()
				return nil
			}
			return cli.PrintWarning(output, err.Error())
		}
		artifact.Deps = a.Deps

		if bar != nil {
			if err := bar.Add(1); err != nil {
				return err
			}
		} else {
			var verbInfos *types.LcVerboseInfo
			if verbose {
				verbInfos = &types.LcVerboseInfo{
					LedgerName: artifact.Ledger,
					LocalSID:   api.GetSignerIDByApiKey(u.Client.ApiKey),
					ApiKey:     u.Client.ApiKey,
				}
			}
			cli.PrintLc(output, types.NewLcResult(artifact, verified, verbInfos))
		}
	}
	if lenArtifacts > 1 && output == "" {
		color.Set(meta.StyleSuccess())
		fmt.Printf("notarized %d items", lenArtifacts)
		color.Unset()
		fmt.Println()
	}
	return nil
}

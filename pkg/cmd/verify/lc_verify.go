package verify

import (
	"fmt"
	"strconv"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/cmd/internal/cli"
	"github.com/codenotary/cas/pkg/cmd/internal/types"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func lcVerify(cmd *cobra.Command, a *api.Artifact, user *api.LcUser, signerID string, uid string, verbose bool, output string) (err error) {
	ar, verified, err := user.LoadArtifact(
		a.Hash,
		signerID,
		uid,
		0,
		map[string][]string{meta.CasCmdHeaderName: {meta.CasVerifyCmdHeaderValue}})
	if err != nil {
		if err == api.ErrNotFound {
			err = fmt.Errorf("%s was not notarized", a.Hash)
			viper.Set("exit-code", strconv.Itoa(meta.StatusUnknown.Int()))
		}
		if err == api.ErrNotVerified {
			color.Set(meta.StyleError())
			fmt.Println("the ledger is compromised. Please contact the Community Attestation Service administrators")
			color.Unset()
			fmt.Println()
			viper.Set("exit-code", strconv.Itoa(meta.StatusUnknown.Int()))
		}
		return cli.PrintWarning(output, err.Error())
	}
	if ar.Revoked != nil && !ar.Revoked.IsZero() {
		viper.Set("exit-code", strconv.Itoa(meta.StatusApikeyRevoked.Int()))
		ar.Status = meta.StatusApikeyRevoked
	}

	ar.IncludedIn = a.IncludedIn
	ar.Deps = a.Deps

	if !verified {
		color.Set(meta.StyleError())
		fmt.Println("the ledger is compromised. Please contact the Community Attestation Service administrators")
		color.Unset()
		fmt.Println()
		viper.Set("exit-code", strconv.Itoa(meta.StatusUnknown.Int()))
		ar.Status = meta.StatusUnknown
	}

	exitCode, err := cmd.Flags().GetInt("exit-code")
	if err != nil {
		return err
	}
	// if exitCode == CasDefaultExitCode user didn't specify to use a custom exit code in case of success.
	// In that case we return the ar.Status as exit code.
	// User defined exit code is returned only if the viper exit-code status is == 0 (status trusted)
	if exitCode == meta.CasDefaultExitCode && viper.GetInt("exit-code") == 0 {
		viper.Set("exit-code", strconv.Itoa(ar.Status.Int()))
	}
	var verbInfos *types.LcVerboseInfo
	if verbose {
		verbInfos = &types.LcVerboseInfo{
			LedgerName: ar.Ledger,
			LocalSID:   api.GetSignerIDByApiKey(user.Client.ApiKey),
			ApiKey:     user.Client.ApiKey,
		}
	}
	cli.PrintLc(output, types.NewLcResult(ar, verified, verbInfos))

	return
}

/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package list

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/cmd/internal/cli"
	"github.com/codenotary/cas/pkg/cmd/internal/types"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/codenotary/cas/pkg/store"
	immuschema "github.com/codenotary/immudb/pkg/api/schema"
)

const keyPrefix = "_ITEM.API-KEY-FULL."
const defaultLimit = 10

// NewCommand returns the cobra command for `vcn list`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Returns the history of operations made with API key",
		Long: `
Returns the history of operations made with API key

Environment variables:
CAS_HOST=
CAS_PORT=
CAS_CERT=
CAS_SKIP_TLS_VERIFY=false
CAS_NO_TLS=false
CAS_API_KEY=
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: runList,
		Args: func(cmd *cobra.Command, args []string) error {
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
			return cobra.OnlyValidArgs(cmd, args)
		},
		Example: `
cas list
cas list --api-key <APIkey>
cas list --start 2020/10/28-08:00:00 --end 2020/10/28-17:00:00 --first 10
cas list --start 2020/10/28-16:00:00 --end 2020/10/28-17:10:00 --last 3
`,
	}

	// ledger compliance flags
	cmd.Flags().String("host", "", meta.CasHostFlagDesc)
	cmd.Flags().String("port", "", meta.CasPortFlagDesc)
	cmd.Flags().String("cert", "", meta.CasCertPathDesc)
	cmd.Flags().Bool("skip-tls-verify", false, meta.CasSkipTlsVerifyDesc)
	cmd.Flags().Bool("no-tls", false, meta.CasNoTlsDesc)
	cmd.Flags().String("api-key", "", meta.CasApiKeyDesc)
	cmd.Flags().String("ledger", "", meta.CasLedgerDesc)

	cmd.Flags().Uint64("first", 0, "set the limit for the first elements filter")
	cmd.Flags().Uint64("last", 0, "set the limit for the last elements filter")

	cmd.Flags().String("start", "", "set the start of date and time range filter. Example 2020/10/28-16:00:00")
	cmd.Flags().String("end", "", "set the end of date and time range filter. Example 2020/10/28-16:00:00")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {

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

	lcUser, err := api.GetOrCreateLcUser(lcApiKey, lcLedger, lcHost, lcPort, lcCert, viper.IsSet("skip-tls-verify"), viper.GetBool("skip-tls-verify"), viper.IsSet("no-tls"), viper.GetBool("no-tls"), nil, false)
	if err != nil {
		return err
	}

	fields := strings.Split(lcUser.Client.ApiKey, ".")
	if len(fields) < 2 {
		return fmt.Errorf("malformed API key '%s'", lcUser.Client.ApiKey)
	}

	hashed := sha256.Sum256([]byte(fields[len(fields)-1]))
	ApiKey := keyPrefix + fields[0] + "." + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hashed[:])

	first, err := cmd.Flags().GetUint64("first")
	if err != nil {
		return err
	}
	last, err := cmd.Flags().GetUint64("last")
	if err != nil {
		return err
	}

	var timeStart, timeEnd time.Time
	start, err := cmd.Flags().GetString("start")
	if err != nil {
		return err
	}
	if start != "" {
		timeStart, err = time.Parse(meta.DateShortForm, start)
		if err != nil {
			return fmt.Errorf("invalid start timestamp format: %w", err)
		}
	}
	end, err := cmd.Flags().GetString("end")
	if err != nil {
		return err
	}
	if end != "" {
		timeEnd, err = time.Parse(meta.DateShortForm, end)
		if err != nil {
			return fmt.Errorf("invalid end timestamp format: %w", err)
		}
	}

	if first == 0 && last == 0 {
		last = defaultLimit
		fmt.Printf("no filter is specified. At maximum last %d items will be returned\n\n", defaultLimit)
	}

	return list(lcUser, ApiKey, timeStart, timeEnd, first, last, output)
}

func list(u *api.LcUser, key string, start, end time.Time, first, last uint64, output string) error {

	var startScore *immuschema.Score = nil
	var endScore *immuschema.Score = nil

	if !start.IsZero() {
		startScore = &immuschema.Score{
			Score: float64(start.UnixNano()), // there is no precision loss. 52 bit are enough to represent seconds.
		}
	}
	if !end.IsZero() {
		endScore = &immuschema.Score{
			Score: float64(end.UnixNano()), // there is no precision loss. 52 bit are enough to represent seconds.
		}
	}

	desc := false
	var limit uint64 = 0

	if first > 0 {
		limit = first
	}
	if last > 0 {
		limit = last
		desc = true
		// it is important to set MaxScore for descending query
		if endScore == nil {
			endScore = &immuschema.Score{
				Score: math.MaxFloat64,
			}
		}
	}

	md := metadata.Pairs(meta.CasPluginTypeHeaderName, meta.CasPluginTypeHeaderValue)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	zItems, err := u.Client.ZScanExt(ctx, &immuschema.ZScanRequest{
		Set:      []byte(key),
		Limit:    limit,
		Desc:     desc,
		MinScore: startScore,
		MaxScore: endScore,
	})
	if err != nil {
		return err
	}

	results := make([]*types.LcResult, len(zItems.Items))
	for i, item := range zItems.Items {
		a, err := api.ZItemToLcArtifact(item)
		if err != nil {
			return err
		}
		results[i] = types.NewLcResult(a, true, nil)
	}

	cli.PrintLcSlice(output, results)

	if output == "" {
		fmt.Printf("Found %d operation(s)\n", len(zItems.Items))
	}

	return nil
}

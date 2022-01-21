/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package inspect

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/cmd/internal/cli"
	"github.com/codenotary/cas/pkg/cmd/internal/types"
	"github.com/codenotary/cas/pkg/meta"
	immuschema "github.com/codenotary/immudb/pkg/api/schema"
	"github.com/fatih/color"
	"github.com/vchain-us/ledger-compliance-go/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func lcInspect(hash string, signerID string, u *api.LcUser, first, last uint64, start, end string, output string) (err error) {
	var contextSignerID string

	if signerID == "" {
		if output == "" {
			fmt.Println("no signer ID provided. Full history of the item is returned")
		}
	} else {
		contextSignerID = signerID
	}

	results, err := GetLcResults(hash, signerID, u, first, last, start, end)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.ResourceExhausted {
				return fmt.Errorf("too many notarizations are returned. Try to use --first or --last filter or datetime range filter")
			}
		}
		return err
	}
	l := len(results)
	if output == "" {
		fmt.Printf(`current signerID `)
		color.Set(meta.StyleAffordance())
		fmt.Printf("%s\n", contextSignerID)
		color.Unset()
		fmt.Printf(`%d notarizations found for "%s"`, l, hash)
		fmt.Println()
		fmt.Println()
	}

	return cli.PrintLcSlice(output, results)
}

func GetLcResults(hash, signerID string, u *api.LcUser, first, last uint64, start, end string) (results []*types.LcResult, err error) {
	md := metadata.Pairs(meta.CasPluginTypeHeaderName, meta.CasPluginTypeHeaderValue)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var key []byte
	if signerID == "" {
		key = []byte(hash)
	} else {
		key = api.AppendPrefix(meta.CasPluginTypeHeaderValue, []byte(signerID))
		key = api.AppendSignerId(hash, key)
	}

	if start != "" || end != "" {
		if signerID == "" {
			key = append([]byte(meta.IndexDateRangePrefix), key...)
		}
		results, err = getTimeRangedResults(ctx, u, key, first, last, start, end)
		if err != nil {
			return nil, err
		}
	} else {
		if signerID == "" {
			results, err = getSignerResults(ctx, key, u, first, last)
			if err != nil {
				return nil, err
			}
		} else {
			results, err = getHistoryResults(ctx, key, u, first, last)
			if err != nil {
				return nil, err
			}
		}
	}

	return results, nil
}

func getSignerResults(ctx context.Context, key []byte, u *api.LcUser, first, last uint64) ([]*types.LcResult, error) {
	var err error
	var zitems *schema.ZItemExtList

	desc := false
	var limit uint64 = 0
	var endScore *immuschema.Score = nil

	if first > 0 {
		limit = first
	}
	if last > 0 {
		limit = last
		desc = true
		// it is important to set MaxScore for descending query
		endScore = &immuschema.Score{
			Score: math.MaxFloat64,
		}
	}

	zitems, err = u.Client.ZScanExt(ctx, &immuschema.ZScanRequest{
		Set:      key,
		Limit:    limit,
		Desc:     desc,
		MaxScore: endScore,
	})
	if err != nil {
		return nil, err
	}

	results := make([]*types.LcResult, len(zitems.Items))
	var i = 0
	for _, v := range zitems.Items {
		lca, err := api.ZItemToLcArtifact(v)
		if err != nil {
			results[i].AddError(err)
		}
		results[i] = types.NewLcResult(lca, true, nil)

		i++
	}
	return results, nil
}

func getHistoryResults(ctx context.Context, key []byte, u *api.LcUser, first, last uint64) ([]*types.LcResult, error) {
	var err error
	var items *schema.ItemExtList

	desc := false
	var limit uint64 = 0

	if first > 0 {
		limit = first
	}
	if last > 0 {
		limit = last
		desc = true
	}

	items, err = u.Client.HistoryExt(ctx, &immuschema.HistoryRequest{
		Key:   key,
		Limit: int32(limit),
		Desc:  desc,
	})
	if err != nil {
		return nil, err
	}

	results := make([]*types.LcResult, len(items.Items))
	var i = 0
	for _, v := range items.Items {
		lca, err := api.ItemToLcArtifact(v)
		if err != nil {
			return nil, err
		}
		results[i] = types.NewLcResult(lca, true, nil)
		if err != nil {
			results[i].AddError(err)
		}
		i++
	}
	return results, nil
}

func getTimeRangedResults(ctx context.Context, u *api.LcUser, set []byte, first, last uint64, start, end string) ([]*types.LcResult, error) {
	var err error
	var zitems *schema.ZItemExtList

	var startScore *immuschema.Score = nil
	var endScore *immuschema.Score = nil

	if start != "" {
		timeStart, err := time.Parse(meta.DateShortForm, start)
		if err != nil {
			return nil, err
		}
		startScore = &immuschema.Score{
			Score: float64(timeStart.UnixNano()), // there is no precision loss. 52 bit are enough to represent seconds.
		}
	}

	if end != "" {
		timeEnd, err := time.Parse(meta.DateShortForm, end)
		if err != nil {
			return nil, err
		}
		endScore = &immuschema.Score{
			Score: float64(timeEnd.UnixNano()), // there is no precision loss. 52 bit are enough to represent seconds.
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

	zitems, err = u.Client.ZScanExt(ctx, &immuschema.ZScanRequest{
		Set:      set,
		Limit:    limit,
		Desc:     desc,
		MinScore: startScore,
		MaxScore: endScore,
		NoWait:   true,
	})
	if err != nil {
		return nil, err
	}

	results := make([]*types.LcResult, len(zitems.Items))
	var i = 0
	for _, v := range zitems.Items {
		lca, err := api.ZItemToLcArtifact(v)
		if err != nil {
			results[i].AddError(err)
		}
		results[i] = types.NewLcResult(lca, true, nil)

		i++
	}
	return results, nil
}

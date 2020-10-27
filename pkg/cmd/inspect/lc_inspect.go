/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package inspect

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/vchain-us/vcn/pkg/meta"

	"fmt"
	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"
	"google.golang.org/grpc/metadata"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
)

func lcInspect(hash string, signerID string, u *api.LcUser, output string) (err error) {
	md := metadata.Pairs(meta.VcnLCPluginTypeHeaderName, meta.VcnLCPluginTypeHeaderValue)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var items *schema.StructuredItemList

	hasher := sha256.New()
	hasher.Write([]byte(u.LcApiKey()))
	contextSignerID := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if signerID == "" {
		if output == "" {
			fmt.Println("no signer ID provided. Full history of the item is returned")
		}
		items, err = u.Client.ZScan(ctx, []byte(hash))
		if err != nil {
			return err
		}
	} else {
		contextSignerID = signerID
		key := api.AppendPrefix(meta.VcnLCPrefix, []byte(signerID))
		key = api.AppendSignerId(hash, key)
		items, err = u.Client.History(ctx, key)
		if err != nil {
			return err
		}
	}

	l := len(items.Items)
	if output == "" {
		fmt.Printf(
			`current signer ID "%s"
%d notarizations found for "%s"

`,
			contextSignerID, l, hash)
	}

	results := make([]types.LcResult, l)
	var i = 0
	for _, v := range items.Items {
		var lca api.LcArtifact
		err = json.Unmarshal(v.Value.Payload, &lca)
		if err != nil {
			return err
		}
		results[i] = *types.NewLcResult(&lca, true)
		if err != nil {
			results[i].AddError(err)
		}
		i++
	}

	return cli.PrintLcSlice(output, results)
}

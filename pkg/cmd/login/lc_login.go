/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package login

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	caserr "github.com/codenotary/cas/internal/errors"
	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/fatih/color"
	"google.golang.org/grpc/metadata"
)

// Execute the login action
func ExecuteLC(host, port, lcCert, lcApiKey, lcLedger string, skipTlsVerifySet, skipTlsVerify, noTlsSet, noTls bool, signingPubKey *ecdsa.PublicKey, skipSigVerify bool, enforceSignatureVerify bool) error {
	if lcApiKey == "" {
		return caserr.ErrNoLcApiKeyEnv
	}

	color.Set(meta.StyleAffordance())
	fmt.Println("Logging into Community Attestation Service.")
	color.Unset()

	u, err := api.GetOrCreateLcUser(lcApiKey, lcLedger, host, port, lcCert, skipTlsVerifySet, skipTlsVerify, noTlsSet, noTls, signingPubKey, false)
	if err != nil {
		return err
	}

	md := metadata.Pairs(meta.CasPluginTypeHeaderName, meta.CasPluginTypeHeaderValue)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err = u.Client.Health(ctx)
	if err != nil {
		return err
	}

	if !skipSigVerify {
		err = u.CheckConnectionPublicKey(enforceSignatureVerify)
		if err != nil {
			return err
		}
	}

	// shouldn't happen
	return nil
}

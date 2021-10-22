/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package api

import (
	"github.com/codenotary/cas/internal/errors"
	"github.com/codenotary/cas/internal/logs"

	"os"

	"github.com/codenotary/cas/pkg/meta"
	"github.com/sirupsen/logrus"
)

// PublicCNLCVerify allow connection and verification on CNLC ledger with a single call using environment variables.
// LcLedger parameter is used when a cross-ledger key is provided in order to specify the ledger on which future operations will be directed. Empty string is accepted.
// signerID parameter is used to filter result on a specific signer ID. If empty value is provided is used the current logged signerID value.
func LcVerifyEnv(hash, lcLedger, signerID string) (a *LcArtifact, err error) {
	lcHost := os.Getenv(meta.CasHost)
	lcPort := os.Getenv(meta.CasPort)
	lcCert := os.Getenv(meta.CasCert)
	lcSkipTlsVerify := os.Getenv(meta.CasSkipTlsVerify)
	lcNoTls := os.Getenv(meta.CasNoTls)
	return PublicCNLCVerify(hash, lcLedger, signerID, lcHost, lcPort, lcCert, lcSkipTlsVerify == "true", lcNoTls == "true")
}

// PublicCNLCVerify allow connection and verification on CNLC ledger with a single call.
// LcLedger parameter is used when a cross-ledger key is provided in order to specify the ledger on which future operations will be directed. Empty string is accepted
// signerID parameter is used to filter result on a specific signer ID. If empty value is provided is used the current logged signerID value.
func PublicCNLCVerify(hash, lcLedger, signerID, lcHost, lcPort, lcCert string, lcSkipTlsVerify, lcNoTls bool) (a *LcArtifact, err error) {
	logger().WithFields(logrus.Fields{
		"hash": hash,
	}).Trace("LcVerify")

	apiKey := os.Getenv(meta.CasApiKey)
	if apiKey == "" {
		logs.LOG.Trace("Lc api key provided (environment)")
		return nil, errors.ErrNoLcApiKeyEnv
	}

	client, err := NewLcClient(apiKey, lcLedger, lcHost, lcPort, lcCert, lcSkipTlsVerify, lcNoTls, nil)
	if err != nil {
		return nil, err
	}

	lcUser := &LcUser{Client: client}

	err = lcUser.Client.Connect()
	if err != nil {
		return nil, err
	}

	if hash != "" {
		a, _, err = lcUser.LoadArtifact(
			hash,
			signerID,
			"",
			0,
			map[string][]string{meta.CasCmdHeaderName: {meta.CasVerifyCmdHeaderValue}})
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

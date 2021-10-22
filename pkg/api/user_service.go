/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package api

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"

	caserr "github.com/codenotary/cas/internal/errors"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/codenotary/cas/pkg/store"
	sdk "github.com/vchain-us/ledger-compliance-go/grpcclient"
)

// User represent a CodeNotary platform user.
type LcUser struct {
	Client     *sdk.LcClient
	PrivateKey *ed25519.PrivateKey
}

const (
	HttpPort    = "80"
	HttpsPort   = "443"
	DefaultHost = "cas.codenotary.com"
)

// GetOrCreateLcUser returns a new User instance configured with provided parameters or an error.
// Before creating a new user it looks for a context one
// LcLedger parameter is used when a cross-ledger key is provided in order to specify the ledger on which future operations will be directed. Empty string is accepted
func GetOrCreateLcUser(lcApiKey, lcLedger, host, port, lcCert string, skipTlsVerifySet, skipTlsVerify, noTlsSet, noTls bool, signingPubKey *ecdsa.PublicKey, publicAuth bool) (*LcUser, error) {
	if lcApiKey == "" && !publicAuth {
		return nil, caserr.ErrNoLcApiKeyEnv
	}
	context := store.Config().CurrentContext
	if port == "" {
		port = context.LcPort
		if port == "" {
			// set port for set up a connection to a Community Attestation Service (default 443). If --no-tls is provided default port will be 80
			if noTlsSet && noTls {
				port = HttpPort
			} else {
				port = HttpsPort
			}
		}
	}
	if host == "" {
		if context.LcHost == "" {
			if publicAuth {
				host = DefaultHost
			} else {
				return nil, caserr.ErrLoginRequired
			}
		} else {
			host = context.LcHost
		}
	}
	if lcLedger == "" {
		lcLedger = context.LcCert // it seems context.LcCert is in fact stores ledger. Confirm it and rename the field
	}
	if !noTlsSet {
		noTls = context.LcNoTls
	}
	if !skipTlsVerifySet {
		skipTlsVerify = context.LcSkipTlsVerify
	}
	//check if an lcUser is present inside the context
	client, err := NewLcClient(lcApiKey, lcLedger, host, port, lcCert, skipTlsVerify, noTls, signingPubKey)
	if err != nil {
		return nil, err
	}

	if err := client.Connect(); err != nil {
		return nil, err
	}
	// if client connects parameters are stored on cas config file
	store.Config().NewLcUser(host, port, lcCert, skipTlsVerify, noTls)
	if err := store.SaveConfig(); err != nil {
		return nil, err
	}
	return &LcUser{
		Client: client,
	}, nil
}

// NewLcUserVolatile returns a new User instance without a backing cfg file.
func NewLcUserVolatile(lcApiKey, lcLedger string, host string, port string) *LcUser {
	p, _ := strconv.Atoi(port)
	return &LcUser{
		Client: sdk.NewLcClient(
			sdk.ApiKey(lcApiKey),
			sdk.MetadataPairs([]string{
				meta.CasLedgerHeaderName, lcLedger,
				meta.CasVersionHeaderName, meta.Version(),
			}),
			sdk.Host(host),
			sdk.Port(p),
			sdk.Dir(store.CurrentConfigFilePath())),
	}
}

// Config returns the User configuration object (see store.User), if any.
// It returns nil if the User is not properly initialized.
func (u User) User() *store.User {
	if u.cfg != nil {
		return u.cfg
	}
	return nil
}

func GetSignerIDByApiKey(lcApiKey string) string {
	ris := strings.Split(lcApiKey, ".")
	// new apikey format {friendlySignerID}.{secret}
	if len(ris) > 1 {
		return strings.Join(ris[:len(ris)-1], ".")
	}
	// old apikey format {secret}
	hash := sha256.Sum256([]byte(lcApiKey))
	return base64.URLEncoding.EncodeToString(hash[:])
}

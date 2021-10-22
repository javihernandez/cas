/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package meta

import (
	"fmt"
	"log"
	"runtime"

	"github.com/fatih/color"
)

// Level is the type for all possible signature levels
type Level int64

// Status is the type for all possible asset statuses
type Status int64

// Allowed Level values
const (
	LevelDisabled         Level = -1
	LevelUnknown          Level = 0
	LevelEmailVerified    Level = 1
	LevelSocialVerified   Level = 2
	LevelIDVerified       Level = 3
	LevelLocationVerified Level = 4
	LevelCNLC             Level = 98
	LevelVchain           Level = 99
)

// Allowed Status values
const (
	StatusTrusted       Status = 0
	StatusUntrusted     Status = 1
	StatusUnknown       Status = 2
	StatusUnsupported   Status = 3
	StatusApikeyRevoked Status = 4
)

// cas environment variable names
const (
	CasApiKey        string = "CAS_API_KEY"
	CasHost          string = "CAS_HOST"
	CasPort          string = "CAS_PORT"
	CasCert          string = "CAS_CERT"
	CasNoTls         string = "CAS_NO_TLS"
	CasSkipTlsVerify string = "CAS_SKIP_TLS_VERIFY"
)

const CasExitCode string = "override default exit codes in case of success"

const CasPrefix string = "cas"

// Community Attestation Service
const CasPluginTypeHeaderName string = "plugin-type"
const CasLedgerHeaderName string = "ledger"
const CasVersionHeaderName string = "version"
const CasPluginTypeHeaderValue string = "vcn"
const CasCmdHeaderName = "cas-command"
const CasNotarizeCmdHeaderValue = "notarize"
const CasVerifyCmdHeaderValue = "verify"

const CasHostFlagDesc string = "if set with host, action will be route to a Community Attestation Service"
const CasPortFlagDesc string = "set port for set up a connection to a Community Attestation Service (default 443). If --no-tls is provided default port will be 80"
const CasCertPathDesc string = "local or absolute path to a certificate file needed to set up tls connection to a Community Attestation Service"
const CasSkipTlsVerifyDesc string = "disables tls certificate verification when connecting to a Community Attestation Service"
const CasNoTlsDesc string = "allow insecure connections when connecting to a Community Attestation Service"
const CasApiKeyDesc string = "Community Attestation Service api key"
const CasLedgerDesc string = "Community Attestation Service ledger. Required when a multi-ledger API key is used."
const CasCIAttribDesc string = "detect CI environment variables context if presents and inject "
const CasUidDesc string = "authenticate on a specific artifact uid"
const CasSigningPubKeyFileNameDesc string = "specify a public key file path to verify signature in messages when connected to a Community Attestation Service. If no public key file is specified but server is signig messages is possible an interactive confirmation of the fingerprint. When confirmed the public key is stored in ~/.cas-trusted-signing-pub-key file."
const CasSigningPubKeyDesc string = "specify a public key to verify signature in messages when connected to a Community Attestation Service. It's required a valid ECDSA key content without header and footer. Ex: --signing-pub-key=\"MFkwE...y5i4w==\""
const CasEnforceSignatureVerifyDesc string = "if this flag is provided cas will disable signature auto trusting when connecting to a new Community Attestation Service"

const BomEntryKeyName string = "BOM"

// UserAgent returns the CAS User-Agent string
func UserAgent() string {
	// Syntax reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent#Syntax
	return fmt.Sprintf("%s/%s (%s; %s)", CasPluginTypeHeaderValue, Version(), runtime.GOOS, runtime.GOARCH)
}

// String returns the name of the given level as string.
func (l Level) String() string {
	switch l {
	case LevelDisabled:
		return "DISABLED"
	case LevelUnknown:
		return "0 - UNKNOWN"
	case LevelEmailVerified:
		return "1 - EMAIL_VERIFIED"
	case LevelSocialVerified:
		return "2 - SOCIAL_VERIFIED"
	case LevelIDVerified:
		return "3 - ID_VERIFIED"
	case LevelLocationVerified:
		return "4 - LOCATION_VERIFIED"
	case LevelVchain:
		return "99 - VCHAIN"
	default:
		log.Fatal("unsupported level: ", int64(l))
		return ""
	}
}

// String returns the name of the given status as string
func (s Status) String() string {
	switch s {
	case StatusTrusted:
		return "TRUSTED"
	case StatusUntrusted:
		return "UNTRUSTED"
	case StatusUnknown:
		return "UNKNOWN"
	case StatusUnsupported:
		return "UNSUPPORTED"
	case StatusApikeyRevoked:
		return "REVOKED"
	default:
		log.Fatal("unsupported status: ", int64(s))
		return ""
	}
}
func (s Status) Int() int {
	return int(s)
}

// StatusNameStyled returns the colorized name of the given status as string
func StatusNameStyled(status Status) string {
	c, s, b := StatusColor(status)
	return color.New(c, s, b).Sprintf(status.String())
}

const DateShortForm = "2006/1/2-15:04:05"
const IndexDateRangePrefix = "_INDEX.ITEM.INSERTION-DATE."

const CasDefaultExitCode = 0

const CasSigningPubKeyFileName = ".cas-trusted-signing-pub-key"

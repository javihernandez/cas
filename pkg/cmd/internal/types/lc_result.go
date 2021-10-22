/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package types

import (
	"github.com/codenotary/cas/pkg/api"
)

type LcResult struct {
	api.LcArtifact `yaml:",inline"`
	Verified       bool           `json:"verified" yaml:"verified" cas:"Verified"`
	Verbose        *LcVerboseInfo `yaml:"verbose,omitempty" cas:"Verbose"`
	Errors         []error        `json:"error,omitempty" yaml:"error,omitempty"`
}

type LcVerboseInfo struct {
	LedgerName string `json:"ledgerName" yaml:"ledgerName" cas:"LedgerName"`
	LocalSID   string `json:"localSID" yaml:"localSID" cas:"LocalSID"`
	ApiKey     string `json:"apiKey" yaml:"apiKey" cas:"ApiKey"`
}

func (r *LcResult) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

func NewLcResult(lca *api.LcArtifact, verified bool, verbose *LcVerboseInfo) *LcResult {

	var r LcResult

	switch true {
	case lca != nil:
		r = LcResult{LcArtifact: *lca, Verified: verified, Verbose: verbose}
	default:
		r = LcResult{}
	}

	return &r
}

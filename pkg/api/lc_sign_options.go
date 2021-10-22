/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package api

import (
	"github.com/codenotary/cas/pkg/meta"
	"github.com/vchain-us/ledger-compliance-go/schema"
)

// SignOption is a functional option for signing operations
type LcSignOption func(*lcSignOpts) error

type lcSignOpts struct {
	status meta.Status
	bom    []*schema.VCNDependency
}

func makeLcSignOpts(opts ...LcSignOption) (o *lcSignOpts, err error) {
	o = &lcSignOpts{
		status: meta.StatusTrusted,
	}

	for _, option := range opts {
		if option == nil {
			continue
		}
		if err := option(o); err != nil {
			return nil, err
		}
	}

	return
}

// SignWithStatus returns the functional option for the given status.
func LcSignWithStatus(status meta.Status) LcSignOption {
	return func(o *lcSignOpts) error {
		o.status = status
		return nil
	}
}

func LcSignWithBom(bom []*schema.VCNDependency) LcSignOption {
	return func(o *lcSignOpts) error {
		o.bom = bom
		return nil
	}
}

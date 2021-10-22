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
)

// SignOption is a functional option for signing operations
type SignOption func(*signOpts) error

type signOpts struct {
	status     meta.Status
}

// SignWithStatus returns the functional option for the given status.
func SignWithStatus(status meta.Status) SignOption {
	return func(o *signOpts) error {
		o.status = status
		return nil
	}
}

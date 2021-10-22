/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package api

import (
	"github.com/codenotary/cas/pkg/store"
)

// User represent a CodeNotary platform user.
type User struct {
	cfg *store.User
}

// ClearAuth deletes the stored authentication token.
func (u *User) ClearAuth() {
	if u != nil && u.cfg != nil {
		u.cfg.Token = ""
	}
}

// Config returns the User configuration object (see store.User), if any.
// It returns nil if the User is not properly initialized.
func (u User) Config() *store.User {
	if u.cfg != nil {
		return u.cfg
	}
	return nil
}

// UserByCfg configures current user with a custom values
func (u *User) UserByCfg(cfg *store.User) {
	u.cfg = cfg
}

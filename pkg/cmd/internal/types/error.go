/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package types

type Error struct {
	Error string `json:"error"`
}

func (e *Error) String() string {
	if e != nil {
		return e.Error
	}
	return ""
}

func NewError(err error) *Error {
	return &Error{
		Error: err.Error(),
	}
}

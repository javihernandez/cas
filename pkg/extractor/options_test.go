/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package extractor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyOpts struct {
	flag bool
}

func withTestOption() Option {
	return func(o interface{}) error {
		if oo, ok := o.(*dummyOpts); ok {
			oo.flag = true
		}
		return nil
	}
}

func withTestOptionWithError() Option {
	return func(o interface{}) error {
		return errors.New("some error")
	}
}

func TestApply(t *testing.T) {

	opts := &dummyOpts{}

	err := Options([]Option{
		withTestOption(),
	}).Apply(opts)
	assert.NoError(t, err)
	assert.True(t, opts.flag)

	err = Options([]Option{
		withTestOptionWithError(),
	}).Apply(opts)
	assert.Error(t, err)
}

/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package sign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapOpts(t *testing.T) {
	m := mapOpts{}

	err := m.Set("key=value")
	assert.NoError(t, err)

	assert.Equal(t, mapOpts{"key": "value"}, m)
	assert.Equal(t, `{"key":"value"}`, m.String())
	assert.Equal(t, map[string]interface{}{"key": "value"}, m.StringToInterface())
}

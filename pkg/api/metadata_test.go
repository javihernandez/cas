/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	m := Metadata{}

	// Set/Get
	m.Set("key", "value")
	assert.Equal(t, "value", m.Get("key", nil))
	assert.Equal(t, "default", m.Get("nonExistingKey", "default"))

	// Multiple values
	m.SetValues(map[string]interface{}{"key": "newValue", "a": "one", "b": 2})
	assert.Equal(t, "newValue", m.Get("key", nil))
	assert.Equal(t, "one", m.Get("a", nil))
	assert.Equal(t, 2, m.Get("b", nil))

	// Swipe
	v := m.swipeString("key")
	assert.Equal(t, "newValue", v)
	assert.Nil(t, m.Get("key", nil))
}

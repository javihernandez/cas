/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package bundle

import (
	"strings"
	"testing"

	// See https://github.com/opencontainers/go-digest#usage
	_ "crypto/sha256"
	_ "crypto/sha512"

	"github.com/stretchr/testify/assert"
)

func TestDescriptor(t *testing.T) {
	d, err := NewDescriptor("path", strings.NewReader("qwertyuiopasdfghjklzxcvbnm"))
	assert.NoError(t, err)
	assert.Equal(t, "sha256:8e5eb603482f00768b60cb17f947e273d6aa7c82ffaf8e589a06f6e841c3cef8", d.Digest.String())
	assert.Equal(t, uint64(26), d.Size)
	assert.Equal(t, []string{"path"}, d.Paths)
}

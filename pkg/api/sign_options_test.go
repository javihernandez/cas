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

	"github.com/codenotary/cas/pkg/meta"

	"github.com/stretchr/testify/assert"
)

func TestSignWithStatus(t *testing.T) {
	o := &signOpts{}
	SignWithStatus(meta.StatusUnsupported)(o)

	assert.Equal(t, meta.StatusUnsupported, o.status)
}

/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package file

import (
	"os"
	"strings"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/extractor/file/internal/sniff"
)

func xInfo(file *os.File, contentType *string) (bool, api.Metadata, error) {
	if strings.HasPrefix(*contentType, "application/") {
		d, err := sniff.File(file)
		if err != nil {
			return false, nil, err
		}
		*contentType = d.ContentType()
		return true, api.Metadata{
			"architecture": strings.ToLower(d.Arch),
			"platform":     d.Platform,
			"file":         d,
		}, nil
	}
	return false, nil, nil
}

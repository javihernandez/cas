/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package file

import (
	"io"
	"net/http"
	"os"

	"github.com/h2non/filetype"
)

func contentType(file *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	file.Seek(0, 0)
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		if n == 0 && err == io.EOF {
			// empty file, no content type
			return "", nil
		}
		return "", err
	}

	kind, err := filetype.Match(buf)
	if err == nil && kind != filetype.Unknown {
		return kind.MIME.Value, nil
	}

	// As fallback, use the net/http package's handy DectectContentType function.
	// Always returns a valid content-type by returning "application/octet-stream"
	// if no others seemed to match.
	contentType := http.DetectContentType(buf)

	return contentType, nil
}

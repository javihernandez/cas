/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package executor

import (
	"io"
)

type Executor interface {
	Exec(cmd []string) (stdout, stderr []byte, exitCode int, err error)
	ReadFile(path string) ([]byte, error)
	ReadDir(path string) (io.ReadCloser, error)
	Close() error
}

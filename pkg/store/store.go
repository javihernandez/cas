/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package store

import (
	"os"
)

// FilePerm holds permission bits that are used for all files that store creates.
const FilePerm os.FileMode = 0600

// DirPerm holds permission bits that are used for all directories that store creates.
const DirPerm os.FileMode = 0700

// DefaultDirName is the name of the store working directory.
const DefaultDirName = ".cas"

const configFilename = "config.json"

const defaultManifestsDir = "manifests"

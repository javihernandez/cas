/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package meta

var version = ""

var static = ""

var gitCommit = ""
var gitBranch = ""

// Version returns the current CodeNotary cas version string
func Version() string {
	return version
}

// StaticBuild returns when the current cas executable has been statically linked against libraries
func StaticBuild() bool {
	return static == "static"
}

// GitRevision returns the current CodeNotary cas git revision string
func GitRevision() string {
	rev := gitCommit
	if gitBranch != "" {
		rev += " (" + gitBranch + ")"
	}
	return rev
}

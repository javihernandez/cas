/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package file

import (
	"regexp"

	"github.com/blang/semver"
)

var verRe = regexp.MustCompile(`v[0-9]+\.[0-9]\.+[0-9]+`)

func inferVer(filename string) string {
	match := verRe.FindStringSubmatch(filename)
	if len(match) > 0 {
		v, err := semver.ParseTolerant(match[0])
		if err == nil {
			return v.String()
		}
	}
	return ""
}

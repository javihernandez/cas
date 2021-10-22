/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package meta

import (
	"github.com/fatih/color"
)

// StatusColor returns color.Attribute(s) for the given status
func StatusColor(status Status) (color.Attribute, color.Attribute, color.Attribute) {
	switch status {
	case StatusTrusted:
		return StyleSuccess()
	case StatusUnknown:
		return StyleWarning()
	case StatusApikeyRevoked:
		return StyleWarning()
	default:
		return StyleError()
	}
}

// StyleAffordance returns color.Attribute(s) for affordance messages
func StyleAffordance() (color.Attribute, color.Attribute, color.Attribute) {
	return color.FgGreen, color.Bold, color.BgBlack
}

// StyleError returns color.Attribute(s) for error messages
func StyleError() (color.Attribute, color.Attribute, color.Attribute) {
	return color.FgRed, color.Bold, color.BgBlack
}

// StyleWarning returns color.Attribute(s) for warning messages
func StyleWarning() (color.Attribute, color.Attribute, color.Attribute) {
	return color.FgYellow, color.Bold, color.BgBlack
}

// StyleSuccess returns color.Attribute(s) for success messages
func StyleSuccess() (color.Attribute, color.Attribute, color.Attribute) {
	return color.FgGreen, color.Bold, color.BgBlack
}

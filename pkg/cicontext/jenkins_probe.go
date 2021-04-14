/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package cicontext

import (
	"os"
)

type jenkisnProbe struct {
	name string
}

func (p *jenkisnProbe) Detect() bool {
	_, ok := os.LookupEnv(CI_JENKINS_KEY)
	return ok
}

func (p *jenkisnProbe) GetName() string {
	return p.name
}

func NewJenkinsProbe() *jenkisnProbe {
	return &jenkisnProbe{
		name: CI_JENKINS_DESC,
	}
}

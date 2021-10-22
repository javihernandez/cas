/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
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

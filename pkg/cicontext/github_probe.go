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
	"strings"
)

type githubProbe struct {
	name string
}

func (p *githubProbe) Detect() bool {
	for _, v := range os.Environ() {
		kv := strings.SplitN(v, "=", 2)
		if strings.Contains(kv[0], CI_GITHUB_PREFIX) {
			return true
		}
	}
	return false
}

func (p *githubProbe) GetName() string {
	return p.name
}

func NewGithubProbe() *githubProbe {
	return &githubProbe{
		name: CI_GITHUB_DESC,
	}
}

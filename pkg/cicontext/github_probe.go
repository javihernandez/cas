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

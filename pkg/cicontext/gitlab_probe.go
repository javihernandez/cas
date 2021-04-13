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

type gitlabProbe struct {
	name string
}

func (p *gitlabProbe) Detect() bool {
	for _, v := range os.Environ() {
		kv := strings.SplitN(v, "=", 2)
		if strings.Contains(kv[0], CI_GITLAB_PREFIX) {
			return true
		}
	}
	return false
}

func (p *gitlabProbe) GetName() string {
	return p.name
}

func NewGitlabProbe() *gitlabProbe {
	return &gitlabProbe{
		name: CI_GITLAB_DESC,
	}
}

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

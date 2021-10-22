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

type contextSaver struct {
	probes []Probe
}

func NewContextSaver() *contextSaver {
	return &contextSaver{
		probes: []Probe{NewGithubProbe(), NewGitlabProbe(), NewJenkinsProbe()},
	}
}

// GetCIContextMetadata returns the CI context metadata
func (cs *contextSaver) GetCIContextMetadata() map[string]interface{} {
	r := map[string]interface{}{}
	for _, k := range CIEnvWhiteList {
		if val, exist := os.LookupEnv(k); exist {
			r[k] = val
		}
	}
	for _, probe := range cs.probes {
		if probe.Detect() {
			r[CI_TYPE_KEY_NAME] = probe.GetName()
			break
		}
	}
	return r
}

// ExtendMetadata extends parent metadata with new data
func ExtendMetadata(parent map[string]interface{}, data map[string]interface{}) map[string]interface{} {
	for k, v := range data {
		parent[k] = v
	}
	return parent
}

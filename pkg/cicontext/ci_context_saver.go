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

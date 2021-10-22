/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package docker

import (
	"os/exec"
	"testing"

	"github.com/codenotary/cas/pkg/uri"
	"github.com/stretchr/testify/assert"
)

func TestDocker(t *testing.T) {
	_, err := exec.Command("docker", "pull", "hello-world").Output()
	if err != nil {
		t.Skip("docker not available")
	}

	u, _ := uri.Parse("docker://hello-world")
	artifacts, err := Artifact(u)
	assert.NoError(t, err)
	assert.NotNil(t, artifacts)
	assert.Equal(t, "docker://hello-world:latest", artifacts[0].Name)
	assert.Regexp(t, "[0-9a-f]{64}", artifacts[0].Hash)
	assert.NotZero(t, artifacts[0].Size)
}

func TestInferVer(t *testing.T) {
	testCases := map[string]string{
		"golang:1.12-stretch": "1.12-stretch",
		"golang:latest":       "",
	}

	for tag, ver := range testCases {
		i := image{
			RepoTags: []string{tag},
		}
		assert.Equal(
			t,
			ver,
			i.inferVer(),
			"wrong version for %s", tag,
		)
	}
}

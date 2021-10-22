/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package docker

import (
	"errors"

	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/bom/executor"
)

type pkgManager interface {
	AllPackages(e executor.Executor, output artifact.OutputOptions) ([]artifact.Dependency, error)
	Type() string
}

var ErrNotFound = errors.New("not found")

func probePackageManager(e executor.Executor) (pkgManager, error) {
	_, _, exitCode, err := e.Exec([]string{"apk", "--version"})
	if err != nil {
		return nil, err
	}
	if exitCode == 0 {
		return &apk{}, nil
	}

	_, _, exitCode, err = e.Exec([]string{"dpkg", "--version"})
	if err != nil {
		return nil, err
	}
	if exitCode == 0 {
		return &dpkg{}, nil
	}

	_, _, exitCode, err = e.Exec([]string{"rpm", "--version"})
	if err != nil {
		return nil, err
	}
	if exitCode == 0 {
		return &rpm{}, nil
	}

	return nil, nil // cannot identify
}

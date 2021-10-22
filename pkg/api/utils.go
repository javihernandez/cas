/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package api

import (
	"fmt"

	"github.com/codenotary/cas/internal/logs"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/dghubble/sling"
	"github.com/sirupsen/logrus"
)

func logger() *logrus.Logger {
	return logs.LOG
}

func makeError(msg string, fields logrus.Fields) error {
	err := fmt.Errorf(msg)
	logger().WithFields(fields).Error(err)
	return err
}

func makeFatal(msg string, fields logrus.Fields) error {
	err := fmt.Errorf(msg)
	logger().WithFields(fields).Fatal(err)
	return err
}

func contains(xs []string, x string) bool {
	for _, a := range xs {
		if a == x {
			return true
		}
	}
	return false
}

func newSling(token string) (s *sling.Sling) {
	s = sling.New()
	s.Add("User-Agent", meta.UserAgent())
	if token != "" {
		s = s.Add("Authorization", "Bearer "+token)
	}
	return s
}

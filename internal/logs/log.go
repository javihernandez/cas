/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

// LOG is the global instance of logrus.Logger
var LOG = logrus.New()

func init() {
	LOG.ExitFunc = func(code int) {} // prevent from exiting immediately on fatal
	ll := os.Getenv("LOG_LEVEL")
	switch ll {
	case "TRACE":
		LOG.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		LOG.SetLevel(logrus.DebugLevel)
	case "INFO":
		LOG.SetLevel(logrus.InfoLevel)
	case "WARN":
		LOG.SetLevel(logrus.WarnLevel)
	case "ERROR":
		LOG.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		LOG.SetLevel(logrus.FatalLevel)
	case "PANIC":
		LOG.SetLevel(logrus.PanicLevel)
	default:
		LOG.SetLevel(logrus.FatalLevel)
	}
}

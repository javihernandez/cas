/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package errors

import "errors"

var ErrNoLcApiKeyEnv = errors.New(`no API key configured. Please set the environment variable CAS_API_KEY=<API-KEY> or use --api-key flag on each request before running any commands`)
var ErrPubAuthNoSignerID = errors.New(`when no api key is provided signerID is mandatory`)
var ErrNoHost = errors.New(`a host name is required in order to login`)
var ErrLoginRequired = errors.New(`you need to be logged in`)

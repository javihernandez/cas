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

// GetCIContextMetadata returns the CI context metadata
func GetCIContextMetadata() map[string]interface{} {
	r := map[string]interface{}{}
	for _, v := range os.Environ() {
		kv := strings.SplitN(v, "=", 2)
		if _, ok := CIEnvWhiteList[kv[0]]; ok {
			if len(kv) == 1 {
				r[kv[0]] = ""
			} else {
				r[kv[0]] = kv[1]
			}
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

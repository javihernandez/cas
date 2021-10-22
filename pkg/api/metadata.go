/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package api

// Metadata holds custom artifact attributes
type Metadata map[string]interface{}

func (m *Metadata) init() {
	if (*m) == nil {
		(*m) = Metadata{}
	}
}

// SetValues sets given values into this Metadata instance
func (m *Metadata) SetValues(values map[string]interface{}) {
	m.init()
	for k, v := range values {
		(*m)[k] = v
	}
}

// Set sets the value for given key
func (m *Metadata) Set(key string, value interface{}) {
	m.init()
	(*m)[key] = value
}

// Get returns the value for the given key, if any, otherwise returns defaultValue
func (m Metadata) Get(key string, defaultValue interface{}) interface{} {
	if m == nil {
		return defaultValue
	}
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

func (m Metadata) swipeString(key string) string {
	v := m.Get(key, "")
	delete(m, key)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

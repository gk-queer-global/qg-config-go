/*
Copyright 2022 The Queer Global Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package configuration

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const configTTLField = "ConfigTTL"

// ParseConfigTTL decodes a ConfigTTL value from JSON init file content.
//
// Accepted forms:
//   - Go duration strings (for example "10s", "5m") via time.ParseDuration.
//   - JSON numbers interpreted as nanoseconds, matching encoding/json behavior
//     for time.Duration (for example 10000000000 is ten seconds).
//
// Zero duration is treated as unset and returns nil. Negative values are rejected.
func ParseConfigTTL(raw json.RawMessage) (*time.Duration, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		return parseConfigTTLString(asString)
	}

	var asNumber json.Number
	if err := json.Unmarshal(raw, &asNumber); err == nil {
		return parseConfigTTLNumber(asNumber)
	}

	return nil, fmt.Errorf(
		"%s: invalid value %s", configTTLField, string(raw))
}

func parseConfigTTLString(value string) (*time.Duration, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("%s: value cannot be empty", configTTLField)
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: invalid duration %q: %w", configTTLField, value, err)
	}
	return finalizeConfigTTL(value, d)
}

func parseConfigTTLNumber(value json.Number) (*time.Duration, error) {
	if strings.ContainsAny(value.String(), ".eE") {
		f, err := value.Float64()
		if err != nil {
			return nil, fmt.Errorf(
				"%s: invalid numeric value %q: %w", configTTLField, value.String(), err)
		}
		if f != float64(int64(f)) {
			return nil, fmt.Errorf(
				"%s: numeric value must be a whole number of nanoseconds, got %q",
				configTTLField, value.String())
		}
	}

	ns, err := value.Int64()
	if err != nil {
		return nil, fmt.Errorf(
			"%s: numeric value overflows int64 nanoseconds %q: %w",
			configTTLField, value.String(), err)
	}

	return finalizeConfigTTL(value.String(), time.Duration(ns))
}

func finalizeConfigTTL(display string, d time.Duration) (*time.Duration, error) {
	if d < 0 {
		return nil, fmt.Errorf(
			"%s: duration must not be negative, got %q", configTTLField, display)
	}
	if d == 0 {
		return nil, nil
	}
	return &d, nil
}

type initParamsJSON struct {
	ConfigType  string                  `json:"ConfigType"`
	Environment *string                 `json:"Environment,omitempty"`
	ConfigTTL   json.RawMessage         `json:"ConfigTTL,omitempty"`
	Aliases     *map[string]string      `json:"Aliases,omitempty"`
	InitValues  *map[string]interface{} `json:"InitValues,omitempty"`
}

func decodeInitParamsJSON(data []byte) (*InitParams, error) {
	var raw initParamsJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	ttl, err := ParseConfigTTL(raw.ConfigTTL)
	if err != nil {
		return nil, err
	}

	return &InitParams{
		ConfigType:  raw.ConfigType,
		Environment: raw.Environment,
		ConfigTTL:   ttl,
		Aliases:     raw.Aliases,
		InitValues:  raw.InitValues,
	}, nil
}

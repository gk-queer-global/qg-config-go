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
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config is a local instance of some config values
type Config map[string]interface{}

// Init contains the information we need to initialzie a configuration fetcher
// for a given config type
type InitParams struct {
	//The type of configuation source we'll be reading from
	ConfigType string `json:"ConfigType"`

	//Environment value appended as a prefix
	Environment *string `json:"Environment,omitempty"`

	//The time a configuration is allowed to exist in memory before
	//being reread from source default is no max / never refresh
	ConfigTTL *time.Duration `json:"ConfigTTL,omitempty"`

	//An map of the aliases for the variables we'll be  reading
	//from our configuration store, where the keys are the aliases
	//and the values are variable names / paths
	Aliases *map[string]string `json:"Aliases,omitempty"`

	//InitValues is an object containing additional params for initialization
	InitValues *map[string]interface{} `json:"InitValues,omitempty"`
}

// ConfigFetcher is an interface for config-fetching / reading types
type ConfigFetcher interface {
	FetchConfig() (*Config, error)
}

func GetConfig(initFile *string) (*Config, error) {

	initPath := "init.json"
	if initFile != nil {
		initPath = *initFile
	}
	init, err := NewInitParams(&initPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	configFetcher, err := NewConfigFetcher(init)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	config, err := configFetcher.FetchConfig()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return config, nil
}

func NewConfigFetcher(params *InitParams) (ConfigFetcher, error) {

	cfgType := strings.ToLower(params.ConfigType)

	switch cfgType {

	case "json":
		{
			jsonFetcher, err := NewLocalJSONFetcher(params)
			if err != nil {
				return nil, fmt.Errorf(
					"could not load local JSON config fetcher: %w", err)
			}
			return jsonFetcher, err
		}

	case "envvar":
		{
			envFetcher, err := NewEnvVarFetcher(params)
			if err != nil {
				return nil, fmt.Errorf(
					"could not load environment variable config fetcher: %w", err)
			}
			return envFetcher, err
		}

	default:
		return nil, fmt.Errorf("ConfigType %s not found", cfgType)
	}
}

func NewInitParams(initPath *string) (*InitParams, error) {

	init := InitParams{}

	if initPath != nil {

		initFile, err := os.Open(*initPath)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to open config file %s: %w", *initPath, err)
		}
		defer initFile.Close()

		initBytes, err := io.ReadAll(initFile)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to read config file %s: %w", *initPath, err)
		}

		decoded, err := decodeInitParamsJSON(initBytes)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to parse config file %s: %w", *initPath, err)
		}
		init = *decoded

	}
	return &init, nil
}

func (cfg *Config) GetMap() *map[string]interface{} {
	configMap := (map[string]interface{})(*cfg)
	return &configMap
}

// Get returns the interface{} value stored at a given key
func (cfg *Config) Get(key string) (interface{}, error) {

	value, ok := ((*cfg)[key])
	if !ok {
		return nil, nil
	}
	return value, nil
}

func (cfg *Config) GetString(key string) (*string, error) {

	value, err := cfg.Get(key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	response, ok := value.(string)
	if !ok {
		response = fmt.Sprintf("%v", value)
	}
	return &response, nil
}

func (cfg *Config) GetInt(key string) (*int, error) {

	value, err := cfg.Get(key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	response, ok := value.(int)
	if !ok {
		intstr := fmt.Sprintf("%v", value)
		resp, err := strconv.Atoi(intstr)
		if err != nil {
			return nil, err
		}
		response = resp
	}
	return &response, nil
}

func (cfg *Config) GetInt32(key string) (*int32, error) {

	value, err := cfg.Get(key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	response, ok := value.(int32)
	if !ok {
		intstr := fmt.Sprintf("%v", value)
		intval, err := strconv.Atoi(intstr)
		if err != nil {
			return nil, err
		}
		response = int32(intval)
	}
	return &response, nil
}

func (cfg *Config) GetInt64(key string) (*int64, error) {

	value, err := cfg.Get(key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	response, ok := value.(int64)
	if !ok {
		intstr := fmt.Sprintf("%v", value)
		n, err := strconv.ParseInt(intstr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"value %v could not be converted to int64: %w", value, err)
		}
		response = n
	}
	return &response, nil
}

func (cfg *Config) GetFloat64(key string) (*float64, error) {

	var response float64

	value, err := cfg.Get(key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	switch resp := value.(type) {

	case float64:
		return &resp, nil

	case float32:
		response = float64(resp)
		return &response, nil

	case int32:
		response = float64(resp)
		return &response, nil

	case int64:
		response = float64(resp)
		return &response, nil

	default:
		return nil, errors.New("tried to convert float64 from incompatible type")
	}
}

func (cfg *Config) GetFloat32(key string) (*float32, error) {

	var response float32

	value, err := cfg.Get(key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	switch resp := value.(type) {

	case float64:
		response = float32(resp)
		return &response, nil

	case float32:
		return &resp, nil

	case int32:
		response = float32(resp)
		return &response, nil

	case int64:
		response = float32(resp)
		return &response, nil

	default:
		return nil, errors.New("tried to convert float64 from incompatible type")
	}
}

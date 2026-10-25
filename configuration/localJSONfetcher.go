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
	"errors"
	"fmt"
	"io"
	"os"
)

type LocalJSONFetcher struct {
	path string
	init *InitParams
}

func NewLocalJSONFetcher(init *InitParams) (*LocalJSONFetcher, error) {

	fetcher := LocalJSONFetcher{}
	fetcher.init = init

	initVals := init.InitValues

	if initVals != nil {

		path, ok := (*initVals)["Path"]
		if !ok {
			path, ok = (*initVals)["path"]
			if !ok {
				return nil, errors.New("\"Path\" value must be provided in InitValues for this type")
			}
		}
		pathstr, ok := path.(string)
		if !ok {
			return nil, errors.New("unable to read string from InitValues.Path")
		}
		fetcher.path = pathstr
	} else {
		return nil, errors.New(
			"a valid InitValues object must be provided for this type")
	}
	return &fetcher, nil

}

func (fetcher *LocalJSONFetcher) FetchConfig() (*Config, error) {

	configFile, err := os.Open(fetcher.path)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to open config file %s", fetcher.path)
	}
	defer configFile.Close()

	cfgBytes, err := io.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read config file %s: %w", fetcher.path, err)
	}

	var cfgValues Config
	if err := json.Unmarshal(cfgBytes, &cfgValues); err != nil {
		return nil, fmt.Errorf(
			"unable to parse config file %s: %w", fetcher.path, err)
	}

	if fetcher.init.Aliases != nil {
		for alias, variable := range *fetcher.init.Aliases {
			cfgValues[alias] = cfgValues[variable]
		}
	}

	return &cfgValues, nil

}

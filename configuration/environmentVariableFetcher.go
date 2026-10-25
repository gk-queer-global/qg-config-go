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
	"fmt"
	"os"
)

type EnvVarFetcher struct {
	init *InitParams
}

func NewEnvVarFetcher(init *InitParams) (*EnvVarFetcher, error) {

	fetcher := EnvVarFetcher{}
	fetcher.init = init
	return &fetcher, nil
}

func (fetcher *EnvVarFetcher) FetchConfig() (*Config, error) {

	if fetcher.init.Aliases == nil {
		return nil, fmt.Errorf(
			"Aliases must be provided for envvar config type")
	}

	var cfgValues = make(Config)

	for alias, variablename := range *fetcher.init.Aliases {
		value := os.Getenv(variablename)
		cfgValues[alias] = value
	}
	return &cfgValues, nil
}

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
	"strings"
	"testing"
)

func TestConfigErrorHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		initPath   string
		wantSubstr []string
		run        func(t *testing.T, initPath string) error
	}{
		{
			name:     "malformed init file",
			initPath: "unittest-init-malformed.json",
			wantSubstr: []string{
				"unittest-init-malformed.json",
			},
			run: func(t *testing.T, initPath string) error {
				t.Helper()
				_, err := NewInitParams(&initPath)
				return err
			},
		},
		{
			name:     "malformed JSON config with aliases",
			initPath: "unittest-init-malformed-config-aliases.json",
			wantSubstr: []string{
				"unittest-config-malformed.json",
			},
			run: func(t *testing.T, initPath string) error {
				t.Helper()
				initObj, err := NewInitParams(&initPath)
				if err != nil {
					return err
				}
				fetcher, err := NewConfigFetcher(initObj)
				if err != nil {
					return err
				}
				_, err = fetcher.FetchConfig()
				return err
			},
		},
		{
			name:     "malformed JSON config without aliases",
			initPath: "unittest-init-malformed-config-no-aliases.json",
			wantSubstr: []string{
				"unittest-config-malformed.json",
			},
			run: func(t *testing.T, initPath string) error {
				t.Helper()
				initObj, err := NewInitParams(&initPath)
				if err != nil {
					return err
				}
				fetcher, err := NewConfigFetcher(initObj)
				if err != nil {
					return err
				}
				_, err = fetcher.FetchConfig()
				return err
			},
		},
		{
			name:     "envvar init with no Aliases",
			initPath: "unittest-init-envvar-no-aliases.json",
			wantSubstr: []string{
				"Aliases",
			},
			run: func(t *testing.T, initPath string) error {
				t.Helper()
				initObj, err := NewInitParams(&initPath)
				if err != nil {
					return err
				}
				fetcher, err := NewConfigFetcher(initObj)
				if err != nil {
					return err
				}
				_, err = fetcher.FetchConfig()
				return err
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.run(t, tt.initPath)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			errMsg := err.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(errMsg, substr) {
					t.Fatalf("error %q should contain %q", errMsg, substr)
				}
			}
		})
	}
}

func TestNewConfigFetcherEnvVarErrorMessage(t *testing.T) {
	initPath := "unittest-init-envvar-no-aliases.json"
	initObj, err := NewInitParams(&initPath)
	if err != nil {
		t.Fatal(err)
	}
	fetcher, err := NewConfigFetcher(initObj)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fetcher.FetchConfig()
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "local JSON") {
		t.Fatalf("envvar path should not mention local JSON: %v", err)
	}
}

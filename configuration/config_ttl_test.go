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
	"strings"
	"testing"
	"time"
)

func TestParseConfigTTL(t *testing.T) {
	t.Parallel()

	tenSeconds := 10 * time.Second
	fiveMinutes := 5 * time.Minute
	ninetyMinutes := 90 * time.Minute
	tenSecondsNanos := time.Duration(10_000_000_000)

	tests := []struct {
		name    string
		raw     string
		want    *time.Duration
		wantErr string
	}{
		{
			name: "duration string seconds",
			raw:  `"10s"`,
			want: &tenSeconds,
		},
		{
			name: "duration string compound",
			raw:  `"1h30m"`,
			want: &ninetyMinutes,
		},
		{
			name: "duration string minutes",
			raw:  `"5m"`,
			want: &fiveMinutes,
		},
		{
			name: "json number nanoseconds",
			raw:  `10000000000`,
			want: &tenSecondsNanos,
		},
		{
			name: "omitted null",
			raw:  `null`,
			want: nil,
		},
		{
			name: "zero duration string unset",
			raw:  `"0s"`,
			want: nil,
		},
		{
			name: "zero nanoseconds unset",
			raw:  `0`,
			want: nil,
		},
		{
			name:    "empty string",
			raw:     `""`,
			wantErr: "cannot be empty",
		},
		{
			name:    "garbage string",
			raw:     `"not-a-duration"`,
			wantErr: "invalid duration",
		},
		{
			name:    "negative duration string",
			raw:     `"-5s"`,
			wantErr: "must not be negative",
		},
		{
			name:    "negative nanoseconds",
			raw:     `-1000000000`,
			wantErr: "must not be negative",
		},
		{
			name:    "invalid json type",
			raw:     `{"bad":"object"}`,
			wantErr: "invalid value",
		},
		{
			name:    "overflow nanoseconds",
			raw:     `9223372036854775808`,
			wantErr: "overflows int64",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseConfigTTL(json.RawMessage(tt.raw))
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error %q should contain %q", err.Error(), tt.wantErr)
				}
				if !strings.Contains(err.Error(), configTTLField) {
					t.Fatalf("error %q should name %q", err.Error(), configTTLField)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertDurationEqual(t, got, tt.want)
		})
	}
}

func TestNewInitParamsConfigTTL(t *testing.T) {
	t.Parallel()

	t.Run("duration string in init file", func(t *testing.T) {
		t.Parallel()
		initPath := "unittest-init-config-ttl-duration.json"
		initObj, err := NewInitParams(&initPath)
		if err != nil {
			t.Fatal(err)
		}
		if initObj.ConfigTTL == nil {
			t.Fatal("expected ConfigTTL to be set")
		}
		if *initObj.ConfigTTL != 10*time.Second {
			t.Fatalf("got %v, want 10s", *initObj.ConfigTTL)
		}
	})

	t.Run("nanosecond number in init file", func(t *testing.T) {
		t.Parallel()
		initPath := "unittest-init-config-ttl-nanoseconds.json"
		initObj, err := NewInitParams(&initPath)
		if err != nil {
			t.Fatal(err)
		}
		if initObj.ConfigTTL == nil {
			t.Fatal("expected ConfigTTL to be set")
		}
		if *initObj.ConfigTTL != 10*time.Second {
			t.Fatalf("got %v, want 10s", *initObj.ConfigTTL)
		}
	})

	t.Run("invalid ConfigTTL returns error not panic", func(t *testing.T) {
		t.Parallel()
		initPath := "unittest-init-config-ttl-invalid.json"
		_, err := NewInitParams(&initPath)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errMsg := err.Error()
		if !strings.Contains(errMsg, "unittest-init-config-ttl-invalid.json") {
			t.Fatalf("error should mention init file: %q", errMsg)
		}
		if !strings.Contains(errMsg, configTTLField) {
			t.Fatalf("error should mention %s: %q", configTTLField, errMsg)
		}
	})

	t.Run("negative ConfigTTL returns error", func(t *testing.T) {
		t.Parallel()
		initPath := "unittest-init-config-ttl-negative.json"
		_, err := NewInitParams(&initPath)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), configTTLField) {
			t.Fatalf("error should mention %s: %v", configTTLField, err)
		}
	})
}

func TestConfigFileWithConfigTTL(t *testing.T) {
	initPath := "unittest-init.json"
	initObj, err := NewInitParams(&initPath)
	if err != nil {
		t.Fatal(err)
	}
	if initObj.ConfigTTL == nil || *initObj.ConfigTTL != 10*time.Second {
		t.Fatalf("ConfigTTL = %v, want 10s", initObj.ConfigTTL)
	}
}

func assertDurationEqual(t *testing.T, got, want *time.Duration) {
	t.Helper()
	if got == nil && want == nil {
		return
	}
	if got == nil || want == nil {
		t.Fatalf("got %v, want %v", got, want)
	}
	if *got != *want {
		t.Fatalf("got %v, want %v", *got, *want)
	}
}

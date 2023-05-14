// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"os"
	"testing"
)

func TestGetenv(t *testing.T) {
	key := "MY_ENV_VAR"
	fallback := "default_value"

	// Environment variable is set
	os.Setenv(key, "my_value")
	defer os.Unsetenv(key)
	if result := Getenv(key, fallback); result != "my_value" {
		t.Errorf("Getenv(%q, %q) = %q; expected %q", key, fallback, result, "my_value")
	}

	// Environment variable is not set, fallback value should be returned
	os.Unsetenv(key)
	defer os.Setenv(key, fallback)
	if result := Getenv(key, fallback); result != fallback {
		t.Errorf("Getenv(%q, %q) = %q; expected %q", key, fallback, result, fallback)
	}

	// Environment variable is set but empty, fallback value should be returned
	os.Setenv(key, "")
	defer os.Unsetenv(key)
	if result := Getenv(key, fallback); result != fallback {
		t.Errorf("Getenv(%q, %q) = %q; expected %q", key, fallback, result, fallback)
	}

	os.Clearenv()
}

func TestStringToInt64(t *testing.T) {
	s1 := "123"
	expected1 := int64(123)
	result1 := StringToInt64(s1)
	if result1 != expected1 {
		t.Errorf("StringToInt64(%q) = %d; expected %d", s1, result1, expected1)
	}
}

func TestSetUserAgent(t *testing.T) {
	// Test with version "1.0"
	result := SetUserAgent("1.0")
	expected := DefaultUserAgent + "1.0"
	if result != expected {
		t.Errorf("SetUserAgent returned %q, expected %q", result, expected)
	}

	// Test with version "2.0"
	result = SetUserAgent("2.0")
	expected = DefaultUserAgent + "2.0"
	if result != expected {
		t.Errorf("SetUserAgent returned %q, expected %q", result, expected)
	}
}

// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"testing"
)

func TestGetRandomWaitMessage(t *testing.T) {
	// Call the function and check that it returns a non-empty string
	result := GetRandomWaitMessage()
	if result == "" {
		t.Errorf("Expected a non-empty string, but got an empty string")
	}
}

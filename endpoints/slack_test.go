// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package endpoints

import "testing"

func TestCheckAfterKeyword(t *testing.T) {
	// Test with input that has text after the second word
	input1 := "/docs ask How do I change order?"
	keyword1 := "ask"
	err1 := checkAfterKeyword(input1, keyword1)
	if err1 != nil {
		t.Errorf("checkAfterKeyword(%q, %q) returned an unexpected error: %v", input1, keyword1, err1)
	}

	// Test with input that has no text after the second word
	input2 := "/docs ask"
	keyword2 := "ask"
	err2 := checkAfterKeyword(input2, keyword2)
	if err2 == nil {
		t.Errorf("checkAfterKeyword(%q, %q) did not return an error as expected", input2, keyword2)
	}

	// Test with input that has only whitespace characters after the second word
	input3 := "/docs ask      "
	keyword3 := "ask"
	err3 := checkAfterKeyword(input3, keyword3)
	if err3 == nil {
		t.Errorf("checkAfterKeyword(%q, %q) did not return an error as expected", input3, keyword3)
	}

	// Test with input that has multiple consecutive whitespace characters after the second word
	input4 := "/docs ask     How do I change order?"
	keyword4 := "ask"
	err4 := checkAfterKeyword(input4, keyword4)
	if err4 != nil {
		t.Errorf("checkAfterKeyword(%q, %q) returned an unexpected error: %v", input4, keyword4, err4)
	}
}

func TestDetermineCommand(t *testing.T) {
	// Test with a valid input command
	input1 := "ask"
	expected1 := Ask
	actual1, err1 := determineCommand(input1)
	if err1 != nil {
		t.Errorf("determineCommand(%q) returned an unexpected error: %v", input1, err1)
	}
	if actual1 != expected1 {
		t.Errorf("determineCommand(%q) returned an unexpected SlackCommands value: got %v, want %v", input1, actual1, expected1)
	}

	// Test with an invalid input command
	input2 := "unknown"
	_, err2 := determineCommand(input2)
	if err2 == nil {
		t.Errorf("determineCommand(%q) did not return an error as expected", input2)
	}

	// Test with an empty input command
	input3 := ""
	_, err3 := determineCommand(input3)
	if err3 == nil {
		t.Errorf("determineCommand(%q) did not return an error as expected", input3)
	}
}

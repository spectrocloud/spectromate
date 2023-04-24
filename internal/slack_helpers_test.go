// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestValidateSignature(t *testing.T) {
	signingTestSecret := "8f742231b10e8888abcd99yyyzzz85a5"
	slackSigningBaseString := "v0:123456789:command=/weather&text=94070"
	signatureTest := "v0=9d04a85ecf26e3f62a8e6f3dad8272130cd2845dd6309b81c669ca2d2338a562"

	ctx := context.Background()

	// Test case for a valid signature
	validSignature := validateSignature(ctx, signatureTest, signingTestSecret, slackSigningBaseString)
	if !validSignature {
		t.Errorf("Expected valid signature, but got invalid")
	}

	// Test case for an invalid signature
	invalidSignature := "v0=invalid_signature"
	validSignature = validateSignature(ctx, invalidSignature, signingTestSecret, slackSigningBaseString)
	if validSignature {
		t.Errorf("Expected invalid signature, but got valid")
	}
}

type TestStruct struct {
	Name     string
	IsActive bool
}

func TestDecodeForm(t *testing.T) {
	// Prepare form data
	formData := url.Values{
		"Name":     {"John Doe"},
		"IsActive": {"true"},
	}

	// Initialize the destination struct
	dst := &TestStruct{}

	// Test successful decoding
	err := DecodeForm(formData, dst)
	if err != nil {
		t.Errorf("DecodeForm returned an error: %v", err)
	}

	// Test expected field values
	if dst.Name != "John Doe" {
		t.Errorf("Expected Name to be 'John Doe', but got '%s'", dst.Name)
	}

	if !dst.IsActive {
		t.Errorf("Expected IsActive to be true, but got false")
	}

	// Test unsupported field type
	formData["Age"] = []string{"30"}
	err = DecodeForm(formData, dst)
	if err == nil {
		t.Errorf("Expected an error for unsupported field type, but got nil")
	}

	// Test invalid field name
	formData = url.Values{
		"InvalidField": {"Invalid"},
	}
	err = DecodeForm(formData, dst)
	if err == nil {
		t.Errorf("Expected an error for invalid field name, but got nil")
	}

	// Test invalid bool value
	formData = url.Values{
		"Name":     {"John Doe"},
		"IsActive": {"invalid_bool"},
	}
	err = DecodeForm(formData, dst)
	if err == nil {
		t.Errorf("Expected an error for invalid bool value, but got nil")
	}
}

func TestGetSlackEvent(t *testing.T) {
	// Prepare form data
	formData := url.Values{
		"token":                 {"test_token"},
		"team_id":               {"test_team_id"},
		"trigger_id":            {"test_trigger_id"},
		"text":                  {"test_text"},
		"team_domain":           {"test_team_domain"},
		"channel_id":            {"test_channel_id"},
		"channel_name":          {"test_channel_name"},
		"user_id":               {"test_user_id"},
		"user_name":             {"test_user_name"},
		"command":               {"/test_command"},
		"response_url":          {"http://example.com/response"},
		"api_app_id":            {"test_api_app_id"},
		"is_enterprise_install": {"false"},
	}

	// Create a fake HTTP request with form data
	requestBody := bytes.NewBufferString(formData.Encode())
	request, err := http.NewRequest("POST", "/slack", io.NopCloser(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Call GetSlackEvent
	event, err := GetSlackEvent(request)
	if err != nil {
		t.Errorf("GetSlackEvent returned an error: %v", err)
	}

	// Check if the returned event has the expected values
	if event.Token != "test_token" {
		t.Errorf("Expected Token to be 'test_token', but got '%s'", event.Token)
	}

	if event.TeamID != "test_team_id" {
		t.Errorf("Expected TeamID to be 'test_team_id', but got '%s'", event.TeamID)
	}

	if event.TriggerID != "test_trigger_id" {
		t.Errorf("Expected TriggerID to be 'test_trigger_id', but got '%s'", event.TriggerID)
	}

	if event.Text != "test_text" {
		t.Errorf("Expected Text to be 'test_text', but got '%s'", event.Text)
	}
}

package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestCreateNewConversation(t *testing.T) {
	// Define a mock HTTP server that returns a specific response for testing purposes
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the request has the expected HTTP method and Content-Type header
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Verify that the request body is correct
		expectedBody := MendableAPIRequest{
			ApiKey: "test-api-key",
		}
		var requestBody MendableAPIRequest
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}
		if !reflect.DeepEqual(requestBody, expectedBody) {
			t.Errorf("Expected request body %+v, got %+v", expectedBody, requestBody)
		}

		// Return a mock response with a specific conversation ID
		responseBody := MendableNewConversationResponse{
			ConversationID: 123,
		}
		jsonBody, err := json.Marshal(responseBody)
		if err != nil {
			t.Errorf("Error encoding response body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBody)
	}))
	defer ts.Close()

	// Call the function with the mock server's URL and API key
	conversationID, err := CreateNewConversation(context.Background(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify that the function returned the expected conversation ID
	expectedConversationID := int64(123)
	if conversationID != expectedConversationID {
		t.Errorf("Expected conversation ID %d, got %d", expectedConversationID, conversationID)
	}
}

func TestSendDocsQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := MendablePayload{
			MessageID: 1,
			Sources: []MendableSources{
				{
					ID:      1,
					Content: "Content 1",
					Score:   0.8,
					Link:    "https://example.com/1",
				},
				{
					ID:      2,
					Content: "Content 2",
					Score:   0.7,
					Link:    "https://docs.spectrocloud.com",
				},
				{
					ID:      3,
					Content: "Content 3",
					Score:   0.9,
					Link:    "https://example.com/1",
				},
			},
		}
		response.Answer.Text = "This is a test answer."
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	// Set up the test case
	testQuery := MendableRequestPayload{
		ApiKey:         "test_api_key",
		Question:       "test_question",
		History:        []HistoryItems{},
		ShouldStream:   false,
		ConversationID: 1,
	}

	expectedOutput := MendableQueryResponse{
		ConversationID: 1,
		Question:       "test_question",
		Answer:         "This is a test answer.",
		Links:          []string{"https://example.com/1"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Call the SendDocsQuery function
	result, err := SendDocsQuery(ctx, testQuery, ts.URL)
	if err != nil {
		t.Fatalf("SendDocsQuery returned an error: %v", err)
	}

	// Compare the result with the expected output
	if result.ConversationID != expectedOutput.ConversationID || result.Question != expectedOutput.Question ||
		result.Answer != expectedOutput.Answer || !compareLinks(result.Links, expectedOutput.Links) {
		t.Fatalf("SendDocsQuery returned incorrect result: %+v, expected: %+v", result, expectedOutput)
	}
}

func compareLinks(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

package slackCmds

import (
	"strings"
	"testing"
)

func TestLinksBuilderString(t *testing.T) {
	// Test case 1: Empty slice
	urls1 := []string{}
	expected1 := ":mag: Unable identify a specific documentation URL."

	result1 := linksBuilderString(urls1)
	if result1 != expected1 {
		t.Errorf("Expected: %s, got: %s", expected1, result1)
	}

	// Test case 2: Non-empty slice
	urls2 := []string{"https://example.com/doc1", "https://example.com/doc2"}

	result2 := linksBuilderString(urls2)
	for _, url := range urls2 {
		if !strings.Contains(result2, url) {
			t.Errorf("Expected result2 to contain: %s, but it did not", url)
		}
	}

	// Test case 3: URL with special characters
	urls3 := []string{"https://example.com/doc1?query=value&another=value"}

	result3 := linksBuilderString(urls3)
	for _, url := range urls3 {
		if !strings.Contains(result3, url) {
			t.Errorf("Expected result3 to contain: %s, but it did not", url)
		}
	}
}

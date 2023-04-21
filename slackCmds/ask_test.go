package slackCmds

import (
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"spectrocloud.com/spectromate/internal"
	"spectrocloud.com/spectromate/mock"
)

func TestLinksBuilderString(t *testing.T) {
	// Test case 1: Empty slice
	urls1 := []string{}
	expected1 := internal.DefaultNoSourcesIdentifiedMessage

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
func TestStoreUserEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockCache(ctrl)

	slackEvent := &internal.SlackEvent{
		UserID:    "U123456",
		ChannelID: "C123456",
	}

	history := []internal.HistoryItems{
		{
			Prompt:   "What is the capital of Germany?",
			Response: "The capital of Germany is Berlin.",
		},
	}

	// Test Case 1: Previous cache item is nil.
	mockCache.EXPECT().StoreHashMap(
		context.Background(),
		"docs_bot:user_id:channel_id:U123456:C123456",
		gomock.Any(),
	).Return(nil)

	mockCache.EXPECT().ExpireKey(
		context.Background(),
		"docs_bot:user_id:channel_id:U123456:C123456",
		internal.DefaultCacheExpirationPeriod,
	).Return(nil)

	mendableQueryResponse := &internal.MendableQueryResponse{
		ConversationID: 123456,
		Question:       "What is the capital of France?",
		Answer:         "The capital of France is Paris.",
		Links:          []string{"https://example.com/doc1", "https://example.com/doc2"},
	}

	cacheItem := &internal.CacheItem{
		UserID:         "U123456",
		ChannelID:      "C123456",
		ConversationID: "123456",
		Question:       "What is the capital of France?",
		Answer:         "The capital of France is Paris.",
		History:        history,
		Counter:        "1",
	}

	slackAskRequest := &SlackAskRequest{
		ctx:        context.Background(),
		slackEvent: slackEvent,
		cache:      mockCache,
	}

	err := storeUserEntry(context.Background(), slackAskRequest, *mendableQueryResponse, 1, nil)
	assert.NoError(t, err)

	// Test Case 2: Previous cache item is not nil.
	mockCache.EXPECT().StoreHashMap(
		context.Background(),
		"docs_bot:user_id:channel_id:U123456:C123456",
		gomock.Any(),
	).Return(nil)

	mockCache.EXPECT().ExpireKey(
		context.Background(),
		"docs_bot:user_id:channel_id:U123456:C123456",
		internal.DefaultCacheExpirationPeriod,
	).Return(nil)

	cacheItem.History = append(cacheItem.History, internal.HistoryItems{
		Prompt:   "What is the capital of France?",
		Response: "The capital of France is Paris.",
	})

	slackAskRequest = &SlackAskRequest{
		ctx:        context.Background(),
		slackEvent: slackEvent,
		cache:      mockCache,
	}

	err = storeUserEntry(context.Background(), slackAskRequest, *mendableQueryResponse, 2, cacheItem)
	assert.NoError(t, err)
}

func TestStoreUserEntryAppend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockCache(ctrl)

	slackEvent := &internal.SlackEvent{
		UserID:    "U123456",
		ChannelID: "C123456",
	}

	history := []internal.HistoryItems{
		{
			Prompt:   "What is the capital of Germany?",
			Response: "The capital of Germany is Berlin.",
		},
	}

	cacheItem := &internal.CacheItem{
		UserID:         "U123456",
		ChannelID:      "C123456",
		ConversationID: "123456",
		Question:       "What is the capital of Italy?",
		Answer:         "The capital of Italy is Rome.",
		History:        history,
		Counter:        "2",
	}

	// Expectations for StoreHashMap
	mockCache.EXPECT().StoreHashMap(
		context.Background(),
		"docs_bot:user_id:channel_id:U123456:C123456",
		gomock.Any(),
	).Return(nil).Times(1)

	// Expectations for ExpireKey
	mockCache.EXPECT().ExpireKey(
		context.Background(),
		"docs_bot:user_id:channel_id:U123456:C123456",
		internal.DefaultCacheExpirationPeriod,
	).Return(nil).Times(1)

	mendableQueryResponse := &internal.MendableQueryResponse{
		ConversationID: 123456,
		Question:       "What is the capital of Italy?",
		Answer:         "The capital of Italy is Rome.",
		Links:          []string{"https://example.com/doc1", "https://example.com/doc2"},
	}

	slackAskRequest := &SlackAskRequest{
		ctx:        context.Background(),
		slackEvent: slackEvent,
		cache:      mockCache,
	}

	err := storeUserEntry(context.Background(), slackAskRequest, *mendableQueryResponse, 2, cacheItem)
	assert.NoError(t, err)
}

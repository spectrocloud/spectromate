package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"spectrocloud.com/docs-slack-bot/internal"
)

// The ask command is used to ask a question about the docs.
// The user will be prompted to enter their question.
// Users can ask a question privately or publicly.
// If the user asks a question privately, the bot will respond privately.
// If the user asks a question publicly, the bot will respond publicly.
// Set the isPrivate bool to true to ask a question privately.
func AskCmd(s *SlackRoute, isPrivate bool) {

	var (
		conversationId   int64
		mendableResponse internal.MendableQueryResponse
		requestCounter   int
	)
	log.Debug().Msgf("The message is private: %v", isPrivate)
	// This will get the user's question.
	// Split the string on spaces
	words := strings.Split(s.SlackEvent.Text, " ")[1:]
	lastWord := words[len(words)-1]
	userQuery := strings.Join(words[:len(words)-1], " ") + " " + strings.TrimRight(lastWord, "\r\n")

	// Log the user query
	log.Debug().Msgf("User query: %v", userQuery)

	// Check if a conversation already exists for this user.
	isNewConversation, cacheItem, err := getUserCache(s.ctx, s)
	if err != nil {
		log.Debug().Err(err).Msgf("an error occured when checking for an exiting conversation: %+v", s.SlackEvent)
		return
	}

	switch isNewConversation {
	case true:
		// Create a new conversation.
		log.Debug().Msgf("Creating a new conversation for user: %v", s.SlackEvent.UserID)
		id, err := internal.CreateNewConversation(s.ctx, s.mendableApiKey, internal.MendandableNewConversationURL)
		if err != nil {
			log.Debug().Err(err).Msgf("Error creating new conversation: %+v", s.SlackEvent)
			return
		}

		conversationId = id
		questionRequestItem := internal.MendableRequestPayload{
			ApiKey:         s.mendableApiKey,
			Question:       userQuery,
			History:        []internal.HistoryItems{},
			ConversationID: conversationId,
			ShouldStream:   false,
		}

		mendableResponse, err = internal.SendDocsQuery(s.ctx, questionRequestItem, internal.MendableChatQueryURL)
		if err != nil {
			log.Debug().Err(err).Msgf("Error sending question to Mendable: %+v", s.SlackEvent)
			internal.LogError(err)
			return
		}

		// Set the question counter to 1.
		requestCounter = 1

	default:
		// Use the existing conversation.
		log.Debug().Msgf("Using existing conversation for user: %v", s.SlackEvent.UserID)

		// Parse the conversation ID.
		cID, err := strconv.ParseInt(cacheItem.ConversationID, 10, 64)
		if err != nil {
			log.Debug().Err(err).Msgf("Error parsing conversation ID: %+v", s.SlackEvent)
			internal.LogError(err)
			return
		}

		// Get the current question counter.
		cNew, err := strconv.ParseInt(cacheItem.Counter, 10, 64)
		if err != nil {
			log.Debug().Err(err).Msgf("Error parsing conversation ID: %+v", s.SlackEvent)
			internal.LogError(err)
			return
		}
		requestCounter = int(cNew)
		conversationId = cID
		questionRequestItem := internal.MendableRequestPayload{
			ApiKey:   s.mendableApiKey,
			Question: userQuery,
			History: []internal.HistoryItems{
				{
					Prompt:   cacheItem.Question,
					Response: cacheItem.Answer,
				},
			},
			ConversationID: conversationId,
			ShouldStream:   false,
		}

		mendableResponse, err = internal.SendDocsQuery(s.ctx, questionRequestItem, internal.MendableChatQueryURL)
		if err != nil {
			log.Debug().Err(err).Msgf("Error sending question to Mendable: %+v", s.SlackEvent)
			internal.LogError(err)
			return
		}

		requestCounter++

	}

	log.Debug().Msgf("ChacheItem: %v", cacheItem)

	linksString := linksBuilderString(mendableResponse.Links)
	markdownContent := fmt.Sprintf(`%v`, mendableResponse.Answer)

	err = storeUserEntry(s.ctx, s, mendableResponse, requestCounter)
	if err != nil {
		log.Debug().Err(err).Msgf("Error storing user entry: %+v", s.SlackEvent)
		return
	}

	slackReplyPayload, err := askMarkdownPayload(markdownContent, linksString, "Docs Answer", isPrivate)
	if err != nil {
		log.Info().Err(err).Msg("Error creating markdown payload.")
		return
	}

	err = internal.ReplyWithAnswer(s.SlackEvent.ResponseURL, slackReplyPayload, isPrivate)
	if err != nil {
		log.Info().Err(err).Msg("Error replying with answer.")
		internal.LogError(err)
		return
	}
}

// // createMarkdownPayload creates a Slack payload with a markdown block
func askMarkdownPayload(content, links, title string, isPrivate bool) ([]byte, error) {
	log.Info().Msgf("Incoming Message: %v", content)

	var responseType string
	if isPrivate {
		responseType = "ephemeral"
	} else {
		responseType = "in_channel"
	}

	payload := internal.SlackPayload{
		ResponseType: responseType,
		Blocks: []internal.SlackBlock{
			{
				Type: "header",
				Text: &internal.SlackTextObject{
					Type: "plain_text",
					Text: title,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: &internal.SlackTextObject{
					Type: "mrkdwn",
					Text: content,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: &internal.SlackTextObject{
					Type: "mrkdwn",
					Text: links,
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	return payloadBytes, nil
}

// storeUserEntry stores the user entry in the cache.
// This is used to track the user's conversation.
// The primary key is a combination of the user ID and channel ID.
// - docs_bot:user_id:channel_id
// Values stored in the cache are:
// - User ID
// - Channel ID
// - Conversation ID
// - Timestamp
func storeUserEntry(ctx context.Context, s *SlackRoute, response internal.MendableQueryResponse, counter int) error {

	tNow := time.Now()

	primaryKey := fmt.Sprintf("docs_bot:user_id:channel_id:%s:%s", s.SlackEvent.UserID, s.SlackEvent.ChannelID)

	cacheItem := map[string]interface{}{
		"UserID":         s.SlackEvent.UserID,
		"ChannelID":      s.SlackEvent.ChannelID,
		"ConversationID": response.ConversationID,
		"Question":       response.Question,
		"Answer":         response.Answer,
		"Timestamp":      &tNow,
		"Counter":        counter,
	}

	err := s.cache.HSet(ctx, primaryKey, cacheItem).Err()
	if err != nil {
		log.Error().Err(err).Msg("Error storing user entry in cache.")
		return err
	}
	log.Debug().Msgf("Stored user entry in cache: %v", cacheItem)
	// Set the expiration on the key.
	err = s.cache.Expire(ctx, primaryKey, time.Minute*5).Err()
	if err != nil {
		log.Error().Err(err).Msg("Error setting expiration on cache key.Exiting program.")
		panic(err)
	}

	return err
}

// getUserCache retrieves the entire cache item from the cache.
// If the cache item is not found, it returns false and a nil cache item.
// If an error occurs, it returns false and the error.
func getUserCache(ctx context.Context, s *SlackRoute) (bool, *internal.CacheItem, error) {
	primaryKey := fmt.Sprintf("docs_bot:user_id:channel_id:%s:%s", s.SlackEvent.UserID, s.SlackEvent.ChannelID)

	result, err := s.cache.HGetAll(ctx, primaryKey).Result()
	if err != nil {
		if err == redis.Nil {
			log.Debug().Msgf("User cache not found in cache: %s", primaryKey)
			return true, nil, nil
		}
		log.Error().Err(err).Msg("Error retrieving user cache from cache.")
		return false, nil, err
	}

	if len(result) == 0 {
		log.Debug().Msgf("User cache not found in cache: %s", primaryKey)
		return true, nil, nil
	}

	cacheItem := &internal.CacheItem{
		UserID:         result["UserID"],
		ChannelID:      result["ChannelID"],
		ConversationID: result["ConversationID"],
		Answer:         result["Answer"],
		Question:       result["Question"],
		Counter:        result["Counter"],
	}

	if timestamp, ok := result["Timestamp"]; ok {
		if t, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
			tNow := time.Unix(t, 0)
			cacheItem.Timestamp = &tNow
		}
	}

	log.Debug().Msgf("Retrieved user cache from cache: %v", cacheItem)

	return false, cacheItem, nil
}

// linksBuilderString builds a string of links.
func linksBuilderString(urls []string) string {
	var sb strings.Builder
	sb.WriteString("*Sources*:\n")
	for _, url := range urls {
		sb.WriteString(fmt.Sprintf("- [%s](%s)\n", url, url))
	}
	return sb.String()
}

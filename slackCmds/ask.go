// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package slackCmds

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"spectrocloud.com/spectromate/internal"
)

type SlackAskRequest struct {
	ctx            context.Context
	slackEvent     *internal.SlackEvent
	mendableAPIKey string
	cache          internal.Cache
	version        string
}

func NewSlackAskRequest(ctx context.Context, slackEvent *internal.SlackEvent, mendableAPIKey string, cache internal.Cache, version string) *SlackAskRequest {
	return &SlackAskRequest{ctx, slackEvent, mendableAPIKey, cache, version}
}

// The ask command is used to ask a question about the docs.
// The user will be prompted to enter their question.
// Users can ask a question privately or publicly.
// If the user asks a question privately, the bot will respond privately.
// If the user asks a question publicly, the bot will respond publicly.
// Set the isPrivate bool to true to ask a question privately.
func AskCmd(s *SlackAskRequest, isPrivate bool) {

	var (
		conversationId   int64
		mendableResponse internal.MendableQueryResponse
		requestCounter   int
		globalErr        *error
	)
	// This will run after the current function returns.
	// This will check if an error occurred and send an error message to the user.
	// This acts as a catch all for any errors that may occur and notifiy the user.
	// A pointer value is used to workaround defer's default behavior of evaulating the arguments immediately.
	defer func() {
		errorEval(s.ctx, globalErr, s, isPrivate)
		globalErr = nil
	}()

	// This will get the user's question.
	// Split the string on spaces

	words := strings.Split(s.slackEvent.Text, " ")[1:]
	lastWord := words[len(words)-1]
	userQuery := strings.Join(words[:len(words)-1], " ") + " " + strings.TrimRight(lastWord, "\r\n")

	// Log the user query
	log.Debug().Msgf("User query: %v", userQuery)

	// Check if a conversation already exists for this user.
	isExistingConversation, cacheItem, err := getUserCache(s.ctx, s)
	if err != nil {
		log.Debug().Err(err).Msgf("an error occured when checking for an exiting conversation: %+v", s.slackEvent)
		globalErr = &err
		return
	}

	// This catches any errors that may occur and returns an error message to the user.

	switch isExistingConversation {
	case false:
		// Create a new conversation.
		log.Debug().Msgf("Creating a new conversation for user: %v", s.slackEvent.UserID)
		id, err := internal.CreateNewConversation(s.ctx, s.mendableAPIKey, internal.MendandableNewConversationURL)
		if err != nil {
			log.Debug().Err(err).Msgf("Error creating new conversation: %+v", s.slackEvent)
			globalErr = &err
			return
		}

		conversationId = id
		questionRequestItem := internal.MendableRequestPayload{
			ApiKey:         s.mendableAPIKey,
			Question:       userQuery,
			History:        []internal.HistoryItems{},
			ConversationID: conversationId,
			ShouldStream:   false,
		}

		mendableResponse, err = internal.SendDocsQuery(s.ctx, questionRequestItem, internal.MendableChatQueryURL, s.version)
		if err != nil {
			internal.LogError(err)
			globalErr = &err
			log.Debug().Err(err).Msgf("Error sending question to Mendable: %+v", s.slackEvent)
			return
		}

		// Set the question counter to 1.
		requestCounter = 1

	default:
		// Use the existing conversation.
		log.Debug().Msgf("Using existing conversation for user: %v", s.slackEvent.UserID)

		// Parse the conversation ID.
		cID, err := strconv.ParseInt(cacheItem.ConversationID, 10, 64)
		if err != nil {
			log.Debug().Err(err).Msgf("Error parsing conversation ID: %+v", s.slackEvent)
			internal.LogError(err)
			globalErr = &err
			return
		}

		// Get the current question counter.
		cNew, err := strconv.ParseInt(cacheItem.Counter, 10, 64)
		if err != nil {
			log.Debug().Err(err).Msgf("Error parsing conversation ID: %+v", s.slackEvent)
			internal.LogError(err)
			globalErr = &err
			return
		}
		requestCounter = int(cNew)
		conversationId = cID
		questionRequestItem := internal.MendableRequestPayload{
			ApiKey:         s.mendableAPIKey,
			Question:       userQuery,
			History:        cacheItem.History,
			ConversationID: conversationId,
			ShouldStream:   false,
		}

		mendableResponse, err = internal.SendDocsQuery(s.ctx, questionRequestItem, internal.MendableChatQueryURL, s.version)
		if err != nil {
			log.Debug().Err(err).Msgf("Error sending question to Mendable: %+v", s.slackEvent)
			internal.LogError(err)
			globalErr = &err
			return
		}

		requestCounter++

	}

	log.Debug().Msgf("ChacheItem: %v", cacheItem)

	linksString := linksBuilderString(mendableResponse.Links)
	markdownContent := fmt.Sprintf(`%v`, mendableResponse.Answer)

	err = storeUserEntry(s.ctx, s, mendableResponse, requestCounter, cacheItem)
	if err != nil {
		log.Debug().Err(err).Msgf("Error storing user entry: %+v", s.slackEvent)
		globalErr = &err
		return
	}

	q := fmt.Sprintf(`:question: %v`, mendableResponse.Question)

	slackReplyPayload, err := askMarkdownPayload(markdownContent, q, linksString, "Docs Answer", mendableResponse.MessageID, isPrivate, mendableResponse.Confidence)
	if err != nil {
		log.Info().Err(err).Msg("Error creating markdown payload.")
		globalErr = &err
		return
	}

	err = internal.ReplyWithAnswer(s.slackEvent.ResponseURL, slackReplyPayload, isPrivate)
	if err != nil {
		log.Info().Err(err).Msg("Error when attempting to return the answer back to Slack.")
		internal.LogError(err)
		globalErr = &err
		// Waiting 5 seconds before returning the error to Slack.
	}
}

// // createMarkdownPayload creates a Slack payload with a markdown block
func askMarkdownPayload(content, question, links, title, messageId string, isPrivate bool, confidence string) ([]byte, error) {
	log.Debug().Msgf("Incoming Message: %v", content)

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
					Text: question,
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
			{
				Type: "divider",
			},
			{
				Type: "section",
				Fields: []internal.SlackTextObject{
					{
						Type: "mrkdwn",
						Text: "*Answer Confidence Level:* " + confidence + "%",
					},
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Fields: []internal.SlackTextObject{
					{
						Type: "mrkdwn",
						Text: "*Rate Answer:*",
					},
				},
			},
			{
				Type: "actions",
				Elements: []internal.SlackElements{
					{
						Type:  "button",
						Value: messageId,
						Text: internal.SlackElementText{
							Type:  "plain_text",
							Emoji: true,
							Text:  ":thumbsup:",
						},
						Style:    "primary",
						ActionID: internal.ActionsAskModelPositiveFeedbackID,
					},
					{
						Type:  "button",
						Value: messageId,
						Text: internal.SlackElementText{
							Type:  "plain_text",
							Emoji: true,
							Text:  ":thumbsdown:",
						},
						Style:    "danger",
						ActionID: internal.ActionsAskModelNegativeFeedbackID,
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	log.Debug().Msgf("Slack Answer Payload: %v", string(payloadBytes))
	return payloadBytes, nil
}

// storeUserEntry stores the user entry in the cache.
// This is used to track the user's conversation.
// The primary key is a combination of the user ID and channel ID.
// - docs_bot:user_id:channel_id
// Values stored in the cache are:
// - User ID
// - Channel ID
// - Conversation IDSlackCommandsIntf
// - Timestamp
func storeUserEntry(ctx context.Context, s *SlackAskRequest, response internal.MendableQueryResponse, counter int, previousCacheItem *internal.CacheItem) error {

	primaryKey := fmt.Sprintf("docs_bot:user_id:channel_id:%s:%s", s.slackEvent.UserID, s.slackEvent.ChannelID)

	if previousCacheItem == nil {
		log.Debug().Msg("Previous cache item is nil.")
		previousCacheItem = &internal.CacheItem{}
	}

	previousHistory := previousCacheItem.History

	newHistItem := []internal.HistoryItems{
		{
			Prompt:   response.Question,
			Response: response.Answer,
		},
	}

	newHistory := append(previousHistory, newHistItem...)

	// convert newHistory to a string
	newHistoryString, err := json.Marshal(newHistory)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling new history.")
		return err
	}

	cacheItem := map[string]interface{}{
		"UserID":         s.slackEvent.UserID,
		"ChannelID":      s.slackEvent.ChannelID,
		"ConversationID": response.ConversationID,
		"Question":       response.Question,
		"Answer":         response.Answer,
		"History":        newHistoryString,
		"Counter":        counter,
	}

	err = s.cache.StoreHashMap(ctx, primaryKey, cacheItem)
	if err != nil {
		log.Error().Err(err).Msg("Error storing user entry in cache.")
		return err
	}

	// log.Debug().Msgf("Stored user entry in cache: %v", cacheItem)

	err = s.cache.ExpireKey(ctx, primaryKey, internal.DefaultCacheExpirationPeriod)
	if err != nil {
		log.Error().Err(err).Msg("Error setting expiration on user entry in cache.")
		return err
	}

	return err
}

// getUserCache retrieves the entire cache item from the cache.
// If the cache item is not found, it returns true and a nil cache item.
// If an error occurs, it returns false and the error.
func getUserCache(ctx context.Context, s *SlackAskRequest) (bool, *internal.CacheItem, error) {
	primaryKey := fmt.Sprintf("docs_bot:user_id:channel_id:%s:%s", s.slackEvent.UserID, s.slackEvent.ChannelID)

	ok, result, err := s.cache.GetHashMap(ctx, primaryKey)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving user cache from cache.")
		return false, nil, err
	}

	if !ok {
		log.Debug().Msgf("User cache not found in cache: %s", primaryKey)
		return false, nil, nil
	}

	cacheItem := &internal.CacheItem{
		UserID:         result["UserID"],
		ChannelID:      result["ChannelID"],
		ConversationID: result["ConversationID"],
		Answer:         result["Answer"],
		Question:       result["Question"],
		Counter:        result["Counter"],
	}

	if history, ok := result["History"]; ok {
		if err := json.Unmarshal([]byte(history), &cacheItem.History); err != nil {
			log.Error().Err(err).Msg("Error unmarshalling history.")
		}
	}
	// log.Debug().Msgf("Retrieved user cache from cache: %v", cacheItem)
	return true, cacheItem, nil
}

// linksBuilderString builds a string of links.
func linksBuilderString(urls []string) string {

	if len(urls) == 0 {
		return internal.DefaultNoSourcesIdentifiedMessage
	}

	var sb strings.Builder
	sb.WriteString("*Sources*:\n")
	for _, url := range urls {
		sb.WriteString(fmt.Sprintf("- %s\n", url))
	}
	return sb.String()
}

func errorEval(ctx context.Context, e *error, s *SlackAskRequest, isPrivate bool) {
	if e != nil {
		err := internal.ReplyWithErrorMessage(s.slackEvent.ResponseURL, isPrivate)
		if err != nil {
			log.Error().Err(err).Msg("Error when attempting to return the answer back to Slack.")
			internal.LogError(err)
			return
		}
	}
}

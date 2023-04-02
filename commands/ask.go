package commands

import (
	"context"
	"encoding/json"
	"fmt"
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
// Set the private bool to true to ask a question privately.
func AskCmd(s *SlackRoute, private bool) ([]byte, error) {

	// Check if a conversation already exists for this user.
	_, err := getConversationID(s.ctx, s)
	if err != nil {
		log.Debug().Err(err).Msgf("an error occured when checking for an exiting conversation: %+v", s.SlackEvent)
		return nil, err
	}

	markdownContent := "Aspernatur ut est delectus molestias consequatur quo explicabo nesciunt impedit aut. Doloremque nesciunt fuga exercitationem architecto ut rerum porro soluta ducimus. Unde in eveniet aut magni qui assumenda laborum iusto hic consequatur ad debitis nostrum labore. Deserunt quaerat iusto neque sit labore totam similique corporis magni corrupti. Consequatur et omnis ducimus expedita beatae explicabo blanditiis voluptatem eius aliquam veritatis nulla autem id quia. Et eligendi sunt sit nesciunt architecto quas molestiae excepturi dolorem enim id et dolor cum. Cum consequatur quo nemo ut et ex et non et. Placeat cum esse ut eaque debitis quo quas non molestias accusantium fugit temporibus ut at sequi. Ipsum similique molestiae voluptas mollitia voluptates perferendis deserunt."

	err = storeUserEntry(s.ctx, s)
	if err != nil {
		log.Debug().Err(err).Msgf("Error storing user entry: %+v", s.SlackEvent)
		return nil, err
	}

	returnPayload, err := askMarkdownPayload(markdownContent, "Docs Answer")
	if err != nil {
		log.Info().Err(err).Msg("Error creating markdown payload.")
		return nil, err
	}

	return returnPayload, nil
}

// // createMarkdownPayload creates a Slack payload with a markdown block
func askMarkdownPayload(content, title string) ([]byte, error) {
	log.Info().Msgf("Incoming Message: %v", content)

	payload := internal.SlackPayload{
		ReponseType: "in_channel",
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
func storeUserEntry(ctx context.Context, s *SlackRoute) error {

	tNow := time.Now()

	primaryKey := fmt.Sprintf("docs_bot:user_id:channel_id:%s:%s", s.SlackEvent.UserID, s.SlackEvent.ChannelID)

	cacheItem := map[string]interface{}{
		"UserID":         s.SlackEvent.UserID,
		"ChannelID":      s.SlackEvent.ChannelID,
		"ConversationID": "",
		"Timestamp":      &tNow,
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

// getConversationID retrieves the conversation ID from the cache.
func getConversationID(ctx context.Context, s *SlackRoute) (string, error) {
	primaryKey := fmt.Sprintf("docs_bot:user_id:channel_id:%s:%s", s.SlackEvent.UserID, s.SlackEvent.ChannelID)

	conversationID, err := s.cache.HGet(ctx, primaryKey, "ConversationID").Result()
	if err != nil {
		if err == redis.Nil {
			log.Debug().Msgf("Conversation ID not found in cache: %s", primaryKey)
			return "", nil
		}
		log.Error().Err(err).Msg("Error retrieving conversation ID from cache.")
		return "", err
	}

	log.Debug().Msgf("Retrieved conversation ID from cache: %s", conversationID)

	return conversationID, nil
}

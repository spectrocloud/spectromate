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
// Set the private bool to true to ask a question privately.
func AskCmd(s *SlackRoute, private bool) ([]byte, error) {

	var conversationId int64

	// Get the user's question.
	userQuery := strings.Split(strings.Join(strings.Split(s.SlackEvent.Text, " ")[1:], " "), "?")
	log.Debug().Msgf("User query: %v", userQuery)

	// sleep for 5 seconds to simulate a long running process.
	time.Sleep(5 * time.Second)

	// Check if a conversation already exists for this user.
	isNewConversation, cacheItem, err := getUserCache(s.ctx, s)
	if err != nil {
		log.Debug().Err(err).Msgf("an error occured when checking for an exiting conversation: %+v", s.SlackEvent)
		return nil, err
	}

	switch isNewConversation {
	case true:
		// Create a new conversation.
		log.Debug().Msgf("Creating a new conversation for user: %v", s.SlackEvent.UserID)
		id, err := internal.CreateNewConversation(s.mendableApiKey, internal.MendandableNewConversationURL)
		if err != nil {
			log.Debug().Err(err).Msgf("Error creating new conversation: %+v", s.SlackEvent)
			return nil, err
		}
		conversationId = id
	default:
		// Use the existing conversation.
	}

	log.Debug().Msgf("ChacheItem: %v", cacheItem)

	markdownContent := "Aspernatur ut est delectus molestias consequatur quo explicabo nesciunt impedit aut. Doloremque nesciunt fuga exercitationem architecto ut rerum porro soluta ducimus. Unde in eveniet aut magni qui assumenda laborum iusto hic consequatur ad debitis nostrum labore. Deserunt quaerat iusto neque sit labore totam similique corporis magni corrupti. Consequatur et omnis ducimus expedita beatae explicabo blanditiis voluptatem eius aliquam veritatis nulla autem id quia. Et eligendi sunt sit nesciunt architecto quas molestiae excepturi dolorem enim id et dolor cum. Cum consequatur quo nemo ut et ex et non et. Placeat cum esse ut eaque debitis quo quas non molestias accusantium fugit temporibus ut at sequi. Ipsum similique molestiae voluptas mollitia voluptates perferendis deserunt."

	err = storeUserEntry(s.ctx, s, conversationId)
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
func storeUserEntry(ctx context.Context, s *SlackRoute, conversationId int64) error {

	tNow := time.Now()

	primaryKey := fmt.Sprintf("docs_bot:user_id:channel_id:%s:%s", s.SlackEvent.UserID, s.SlackEvent.ChannelID)

	cacheItem := map[string]interface{}{
		"UserID":         s.SlackEvent.UserID,
		"ChannelID":      s.SlackEvent.ChannelID,
		"ConversationID": conversationId,
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

// getUserCache retrieves the entire cache item from the cache.
// If the cache item is not found, it returns false and a nil cache item.
// If an error occurs, it returns false and the error.
func getUserCache(ctx context.Context, s *SlackRoute) (bool, *internal.CacheItem, error) {
	primaryKey := fmt.Sprintf("docs_bot:user_id:channel_id:%s:%s", s.SlackEvent.UserID, s.SlackEvent.ChannelID)

	result, err := s.cache.HGetAll(ctx, primaryKey).Result()
	if err != nil {
		if err == redis.Nil {
			log.Debug().Msgf("User cache not found in cache: %s", primaryKey)
			return false, nil, nil
		}
		log.Error().Err(err).Msg("Error retrieving user cache from cache.")
		return false, nil, err
	}

	if len(result) == 0 {
		log.Debug().Msgf("User cache not found in cache: %s", primaryKey)
		return false, nil, nil
	}

	cacheItem := &internal.CacheItem{
		UserID:         result["UserID"],
		ChannelID:      result["ChannelID"],
		ConversationID: result["ConversationID"],
	}

	if timestamp, ok := result["Timestamp"]; ok {
		if t, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
			tNow := time.Unix(t, 0)
			cacheItem.Timestamp = &tNow
		}
	}

	log.Debug().Msgf("Retrieved user cache from cache: %v", cacheItem)

	return true, cacheItem, nil
}

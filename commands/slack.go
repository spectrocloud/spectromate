package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"spectrocloud.com/docs-slack-bot/internal"
)

// NewHandlerContext returns a new CounterRoute with a database connection.
func NewHelpHandlerContext(ctx context.Context, signingSecret, mendableApiKey string, r *redis.Client) *SlackRoute {
	return &SlackRoute{ctx, signingSecret, mendableApiKey, &internal.SlackEvent{}, r}
}

func (slack *SlackRoute) SlackHTTPHandler(writer http.ResponseWriter, request *http.Request) {
	log.Debug().Msg("Slack request received.")
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Validate the request signature came from the Spectro Cloud Slack app.
	err := internal.SourceValidation(request.Context(), request, slack.signingSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("Error validating Slack request signature.")
	}

	var event internal.SlackEvent
	event, err = internal.GetSlackEvent(request)
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting slack event.")
	}

	// Set the slack event in the SlackRoute struct so it can be used by the other handlers.
	slack.SlackEvent = &event

	log.Debug().Msgf("User: %+v", event.UserName)
	log.Debug().Msgf("UserId: %+v", event.UserID)
	log.Debug().Msgf("Channel: %+v", event.ChannelName)
	log.Debug().Msgf("ChannelId: %+v", event.ChannelID)
	log.Debug().Msgf("Text: %+v", event.Text)
	log.Debug().Msgf("ResponseUrl: %+v", event.ResponseURL)

	value, err := slack.getHandler(writer, request)
	if err != nil {
		internal.LogError(err)
		log.Debug().Msg("an error occurred in the Slack getHandler.")
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(value))
	if err != nil {
		internal.LogError(err)
		log.Error().Err(err).Msg("Error writing response to the help endpoint.")
	}
}

// getHandler returns a health check response.
func (slack *SlackRoute) getHandler(writer http.ResponseWriter, r *http.Request) ([]byte, error) {

	var returnPayload []byte

	// Get the command from the Slack event.
	userCmd := strings.Split(slack.SlackEvent.Text, " ")[0]
	// Convert the command to a SlackCommands type.
	cmd, err := determineCommnad(userCmd)
	if err != nil {
		log.Debug().Msg("Error converting string to SlackCommands type.")
	}

	switch cmd {
	case Help:
		returnPayload, err = HelpCmd()
		if err != nil {
			internal.LogError(err)
			log.Info().Err(err).Msg(internal.SlackDefaultUserErrorMessage)
			return nil, err
		}
	case Ask:
		err := replyStatus200(slack.SlackEvent.ResponseURL, writer)
		if err != nil {
			log.Info().Err(err).Msg("failed to reply to slack with status 200.")
			return nil, err
		}
		go AskCmd(slack, false)
		if err != nil {
			log.Info().Err(err).Msg(internal.SlackDefaultUserErrorMessage)
			return nil, err
		}
	case PAsk:
		returnPayload, err = AskCmd(slack, true)
		if err != nil {
			log.Info().Err(err).Msg(internal.SlackDefaultUserErrorMessage)
			return nil, err
		}
	default:
		returnPayload, err = HelpCmd()
		if err != nil {
			log.Info().Err(err).Msg(internal.SlackDefaultUserErrorMessage)
			return nil, err
		}
	}

	return returnPayload, err
}

// determineCommnad converts a string to a SlackCommands type.
// An error is returned if the string does not match a SlackCommands type.
func determineCommnad(input string) (SlackCommands, error) {

	cmd, err := SlackCommandsFromString(input)
	if err != nil {
		internal.LogError(err)
		log.Info().Err(err).Msg("Error converting string to SlackCommands type.")
		return cmd, err
	}

	return cmd, nil
}

// replyStatus200 replies to the Slack event with a 200 status.
// This is required to prevent Slack from considering the request a failure.
// Slack requires a response within 3 seconds.
func replyStatus200(reponseURL string, writer http.ResponseWriter) error {

	markdownContent := "Hang tight while I review the docs..."

	_, err := waitMarkdownPayload("Docs Answer", markdownContent)
	if err != nil {
		internal.LogError(err)
		log.Error().Err(err).Msg("Error creating Slack 200 markdown payload.")
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(markdownContent))
	if err != nil {
		internal.LogError(err)
		log.Error().Err(err).Msg("Error writing 200 OK Wait Reply.")
	}

	log.Debug().Msg("Successfully replied to Slack with status 200.")
	return nil
}

func waitMarkdownPayload(title, content string) ([]byte, error) {

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

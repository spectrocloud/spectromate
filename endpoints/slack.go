package endpoints

import (
	"context"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"spectrocloud.com/docs-slack-bot/endpoints/slackcmds"
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

	log.Debug().Msgf("Event: %+v", event)

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
	cmd, err := determineCommand(userCmd)
	if err != nil {
		log.Debug().Msg("Error converting string to SlackCommands type.")
	}

	switch cmd {
	case Help:
		returnPayload, err = slackcmds.HelpCmd()
		if err != nil {
			internal.LogError(err)
			log.Info().Err(err).Msg(internal.SlackDefaultUserErrorMessage)
			return nil, err
		}
	case Ask:
		slackRequestInfo := slackcmds.NewSlackAskRequest(
			slack.ctx,
			slack.SlackEvent,
			slack.mendableApiKey,
			slack.cache,
		)

		reply200Payload, err := internal.ReplyStatus200(slack.SlackEvent.ResponseURL, writer, false)
		if err != nil {
			log.Info().Err(err).Msg("failed to reply to slack with status 200.")
			return nil, err
		}
		returnPayload = reply200Payload
		go slackcmds.AskCmd(slackRequestInfo, false)
	case PAsk:
		slackRequestInfo := slackcmds.NewSlackAskRequest(
			slack.ctx,
			slack.SlackEvent,
			slack.mendableApiKey,
			slack.cache,
		)
		reply200Payload, err := internal.ReplyStatus200(slack.SlackEvent.ResponseURL, writer, true)
		if err != nil {
			log.Info().Err(err).Msg("failed to reply to slack with status 200.")
			return nil, err
		}
		returnPayload = reply200Payload
		go slackcmds.AskCmd(slackRequestInfo, true)
	default:
		returnPayload, err = slackcmds.HelpCmd()
		if err != nil {
			log.Info().Err(err).Msg(internal.SlackDefaultUserErrorMessage)
			return nil, err
		}
	}

	return returnPayload, err
}

// determineCommand converts a string to a SlackCommands type.
// An error is returned if the string does not match a SlackCommands type.
func determineCommand(input string) (SlackCommands, error) {

	cmd, err := SlackCommandsFromString(input)
	if err != nil {
		internal.LogError(err)
		log.Info().Err(err).Msg("Error converting string to SlackCommands type.")
		return cmd, err
	}

	return cmd, nil
}

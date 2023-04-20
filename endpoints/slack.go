package endpoints

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"spectrocloud.com/spectromate/internal"
	"spectrocloud.com/spectromate/slackCmds"
)

// NewHandlerContext returns a new CounterRoute with a database connection.
// func NewSlackHandlerContext(ctx context.Context, signingSecret, mendableApiKey string, r *redis.Client) *SlackRoute {
// 	return &SlackRoute{ctx, signingSecret, mendableApiKey, &internal.SlackEvent{}, r}
// }

func NewSlackHandlerContext(ctx context.Context, signingSecret, mendableApiKey string, c internal.Cache) *SlackRoute {
	return &SlackRoute{ctx, signingSecret, mendableApiKey, &internal.SlackEvent{}, c}
}

func (slack *SlackRoute) SlackHTTPHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	if request.Method != http.MethodPost {
		log.Debug().Msg("Invalid method used for Slack endpoint.")
		writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := writer.Write([]byte("Invalid method used for Slack endpoint."))
		if err != nil {
			internal.LogError(err)
			log.Error().Err(err).Msg("Error writing response to the Slack endpoint.")
		}
		return
	}

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

	var (
		returnPayload []byte
		userCmd       string
	)

	parts := strings.Split(slack.SlackEvent.Text, " ")

	// Check if there are any parts
	if len(parts) > 0 {
		// Trim any leading or trailing whitespace from the first part
		userCmd = strings.Split(slack.SlackEvent.Text, " ")[0]
	} else {
		log.Debug().Msg("No parts found in the Slack event text.")
		userCmd = "help"
	}

	// Get the command from the Slack event.
	err := checkAfterKeyword(slack.SlackEvent.Text, userCmd)
	if err != nil {
		userCmd = "help"
	}

	// Convert the command to a SlackCommands type.
	cmd, err := determineCommand(userCmd)
	if err != nil {
		log.Debug().Msg("Error converting string to SlackCommands type.")
	}

	switch cmd {
	case Help:
		returnPayload, err = slackCmds.HelpCmd()
		if err != nil {
			internal.LogError(err)
			log.Info().Err(err).Msg(internal.SlackDefaultUserErrorMessage)
			return nil, err
		}
	case Ask:
		slackRequestInfo := slackCmds.NewSlackAskRequest(
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
		// Reply back to slack with a 200 status code to avoid the 3 second timeout.
		returnPayload = reply200Payload
		go slackCmds.AskCmd(slackRequestInfo, false)
	case PAsk:
		slackRequestInfo := slackCmds.NewSlackAskRequest(
			slack.ctx,
			slack.SlackEvent,
			slack.mendableApiKey,
			slack.cache,
		)
		// Reply back to slack with a 200 status code to avoid the 3 second timeout.
		reply200Payload, err := internal.ReplyStatus200(slack.SlackEvent.ResponseURL, writer, true)
		if err != nil {
			log.Info().Err(err).Msg("failed to reply to slack with status 200.")
			return nil, err
		}
		returnPayload = reply200Payload
		go slackCmds.AskCmd(slackRequestInfo, true)
	default:
		returnPayload, err = slackCmds.HelpCmd()
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

	if strings.TrimSpace(input) == "" {
		return -1, errors.New("empty command")
	}

	cmd, err := SlackCommandsFromString(input)
	if err != nil {
		internal.LogError(err)
		log.Info().Err(err).Msg("Error converting string to SlackCommands type.")
		return cmd, err
	}
	log.Debug().Msgf("Determined Command: %s", cmd)
	return cmd, nil
}

func checkAfterKeyword(input string, keyword string) error {
	// Split the input string by whitespace using strings.Fields
	parts := strings.Fields(input)

	// Find the index of the keyword in the parts
	keywordIndex := -1
	for i, part := range parts {
		if part == keyword {
			keywordIndex = i
			break
		}
	}

	// Check if there is anything after the keyword
	if keywordIndex < len(parts)-1 {
		// Check if anything after the keyword is empty
		if strings.TrimSpace(parts[keywordIndex+1]) == "" {
			return errors.New("nothing after keyword")
		}
	} else {
		return errors.New("nothing after keyword")
	}

	// No error
	return nil
}

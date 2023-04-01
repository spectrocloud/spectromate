package commands

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"spectrocloud.com/docs-slack-bot/internal"
)

// NewHandlerContext returns a new CounterRoute with a database connection.
func NewHelpHandlerContext(ctx context.Context, signingSecret string) *HelpRoute {
	return &HelpRoute{ctx, signingSecret}
}

func (help *HelpRoute) HelpHTTPHandler(writer http.ResponseWriter, request *http.Request) {
	log.Debug().Msg("Help request received.")
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Validate the request signature came from the Spectro Cloud Slack app.
	err := internal.SourceValidation(request.Context(), request, help.signingSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("Error validating Slack request signature.")
	}

	var event internal.SlackEvent
	event, err = internal.GetSlackEvent(request)
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting slack event.")
	}

	log.Debug().Msgf("User: %+v", event.UserName)
	log.Debug().Msgf("UserId: %+v", event.UserID)
	log.Debug().Msgf("Channel: %+v", event.ChannelName)
	log.Debug().Msgf("ChannelId: %+v", event.ChannelID)
	log.Debug().Msgf("Command: %+v", event.Command)

	value, err := help.getHandler(request)
	if err != nil {
		log.Debug().Msg("Error reply to help request.")
		http.Error(writer, "Error getting reply to help request.", http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(value))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response to the help endpoint.")
	}
}

// getHandler returns a health check response.
func (health *HelpRoute) getHandler(r *http.Request) ([]byte, error) {
	markdownContent := "*Commands*\n\nThe following commands are available:\n\n- *help* - A summary of all available commands.\n- *ask* - Ask a docs related question.\n- *pask* - Ask a docs related question privately and receive a response only you can view."

	returnPayload, err := internal.CreateMarkdownPayload(markdownContent, "Docs Answer")
	if err != nil {
		log.Info().Err(err).Msg("Error creating markdown payload.")
		return nil, err
	}

	return returnPayload, nil
}

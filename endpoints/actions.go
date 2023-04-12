package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"spectrocloud.com/docs-slack-bot/internal"
)

// NewHandlerContext returns a new CounterRoute with a database connection.
func NewActionsHandlerContext(ctx context.Context, signingSecret, mendableApiKey string) *ActionsRoute {
	return &ActionsRoute{ctx, signingSecret, mendableApiKey, &internal.SlackActionEvent{}}
}

func (actions *ActionsRoute) ActionsHTTPHandler(writer http.ResponseWriter, request *http.Request) {
	log.Debug().Msg("Health check request received.")
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	var payload []byte

	if request.Method != http.MethodPost {
		log.Debug().Msg("Invalid request method for /slack/actions.")
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var event internal.SlackActionEvent
	event, err := internal.GetSlackActionsEvent(request)
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting slack event.")
	}

	// Set the slack event in the SlackRoute struct so it can be used by the other handlers.
	actions.ActionsEvent = &event

	log.Debug().Msgf("Slack event: %+v", event)

	value, err := actions.getHandler(request)
	if err != nil {
		log.Debug().Msg("Error getting counter value.")
		http.Error(writer, "Error getting counter value.", http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusOK)
	payload = value
	_, err = writer.Write([]byte(payload))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response to the health endpoint.")
	}
}

// getHandler returns a health check response.
func (health *ActionsRoute) getHandler(r *http.Request) ([]byte, error) {
	return json.Marshal(map[string]string{"status": "OK"})
}

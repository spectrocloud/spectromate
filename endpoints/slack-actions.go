// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package endpoints

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"spectrocloud.com/spectromate/internal"
	"spectrocloud.com/spectromate/slackActions"
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
		log.Debug().Msg("invalid request method for /slack/actions.")
		http.Error(writer, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Validate the request signature came from the Spectro Cloud Slack app.
	err := internal.SourceValidation(actions.ctx, request, actions.signingSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("error validating Slack request signature.")
	}

	var event internal.SlackActionEvent
	event, err = internal.GetSlackActionsEvent(request)
	if err != nil {
		log.Fatal().Err(err).Msg("error getting slack event.")
	}

	// Set the slack event in the SlackRoute struct so it can be used by the other handlers.
	actions.ActionsEvent = &event
	log.Debug().Msgf("Slack event: %+v", event)

	value, err := actions.getHandler(actions, request, &event)
	if err != nil {
		log.Debug().Msg("error getting counter value.")
		http.Error(writer, "error getting counter value.", http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusOK)
	payload = value
	_, err = writer.Write([]byte(payload))
	if err != nil {
		log.Error().Err(err).Msg("error writing response to the health endpoint.")
	}
}

// getHandler invokes the modelFeedbackHandler function from the slackActions package
func (actions *ActionsRoute) getHandler(routeRequest *ActionsRoute, reqeust *http.Request, action *internal.SlackActionEvent) ([]byte, error) {
	var returnPayload []byte

	slackRequestInfo := slackActions.NewSlackActionFeedback(routeRequest.ctx, action, routeRequest.mendableApiKey)

	switch action.Actions[0].ActionID {

	case internal.ActionsAskModelPositiveFeedbackID:
		log.Debug().Msg("Positive feedback action triggered.")
		go slackActions.ModelFeedbackHandler(slackRequestInfo, internal.PositiveFeedbackScore)
	case internal.ActionsAskModelNegativeFeedbackID:
		log.Debug().Msg("Negative feedback action triggered.")
		go slackActions.ModelFeedbackHandler(slackRequestInfo, internal.NegativeFeedbackScore)
	default:
		log.Debug().Msg("Unknown action.")
	}

	return returnPayload, nil

}

// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package slackActions

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/rs/zerolog/log"
	"spectrocloud.com/spectromate/internal"
)

type SlackActionRequest struct {
	ctx            context.Context
	action         *internal.SlackActionEvent
	mendableAPIKey string
}

// SlackActionRequest returns a new SlackActionRequest.
func NewSlackActionRequest(ctx context.Context, action *internal.SlackActionEvent, mendableAPIKey string) *SlackActionRequest {
	return &SlackActionRequest{ctx, action, mendableAPIKey}
}

func ModelFeedbackHandler(action *SlackActionRequest, ratingScore internal.MendableRatingScore) {

	var (
		globalErr *error
		isPrivate bool
	)

	isPrivate = action.action.Container.IsEphemeral

	// This will run after the current function returns.
	// This will check if an error occurred and send an error message to the user.
	// This acts as a catch all for any errors that may occur and notifiy the user.
	// A pointer value is used to workaround defer's default behavior of evaulating the arguments immediately.
	defer func() {
		errorEval(action.ctx, globalErr, action, isPrivate)
		globalErr = nil
	}()

	messageIDRaw := action.action.Actions[0].Value

	// convert the messageID to an int
	messageID, err := strconv.Atoi(messageIDRaw)
	if err != nil {
		log.Debug().Err(err).Msg("error converting messageID to an int.")
		internal.LogError(err)
		globalErr = &err
		return
	}

	err = internal.SendModelRating(action.ctx, messageID, ratingScore, action.mendableAPIKey, internal.MendableRatingFeedbackURL)
	if err != nil {
		log.Debug().Err(err).Msg("error sending model feedback.")
		internal.LogError(err)
		globalErr = &err
		return
	}

	log.Debug().Interface("message", action)
	// prevent panic if the message is missing any of the expected blocks
	if len(action.action.Message.Blocks) < 7 {
		log.Debug().Interface("message", action.action.Message).Msg("message is missing blocks.")

		slackReplyPayload, err := replyWithEmptyMessage(isPrivate, ratingScore)
		if err != nil {
			log.Info().Err(err).Msg("Error creating the rating markdown payload.")
			globalErr = &err
			return
		}

		err = internal.ReplyWithAnswer(action.action.ResponseURL, slackReplyPayload, isPrivate)
		if err != nil {
			log.Info().Err(err).Msg("Error when attempting to return the answer back to Slack.")
			internal.LogError(err)
			globalErr = &err
			return
		}

		return
	}

	var originalAnswer, orginalQuestion, orignalSourceLinks string

	if len(action.action.Message.Blocks) > 4 && action.action.Message.Blocks[4].Text.Text != "" {
		originalAnswer = action.action.Message.Blocks[4].Text.Text
	}

	if len(action.action.Message.Blocks) > 2 && action.action.Message.Blocks[2].Text.Text != "" {
		orginalQuestion = action.action.Message.Blocks[2].Text.Text
	}

	if len(action.action.Message.Blocks) > 6 && action.action.Message.Blocks[6].Text.Text != "" {
		orignalSourceLinks = action.action.Message.Blocks[6].Text.Text
	}

	slackReplyPayload, err := rateFeedbackMarkdownPayload(originalAnswer, orginalQuestion, orignalSourceLinks, isPrivate, ratingScore)
	if err != nil {
		log.Info().Err(err).Msg("Error creating the rating markdown payload.")
		globalErr = &err
		return
	}

	err = internal.ReplyWithAnswer(action.action.ResponseURL, slackReplyPayload, isPrivate)
	if err != nil {
		log.Info().Err(err).Msg("Error when attempting to return the answer back to Slack.")
		internal.LogError(err)
		globalErr = &err
		// Waiting 5 seconds before returning the error to Slack.
		return
	}

	log.Info().Msg("Successfully sent the answer back to Slack.")

}

func replyWithEmptyMessage(isPrivate bool, rating internal.MendableRatingScore) ([]byte, error) {
	var (
		responseType    string
		responseMessage string
	)

	if isPrivate {
		responseType = "ephemeral"
	} else {
		responseType = "in_channel"
	}

	if rating == internal.PositiveFeedbackScore {
		responseMessage = internal.DefaultPositiveRatingMessage
	}

	if rating == internal.NegativeFeedbackScore {
		responseMessage = internal.DefaultNegativeRatingMessage
	}

	payload := internal.SlackPayload{
		ResponseType:    responseType,
		DeleteOriginal:  false,
		ReplaceOriginal: false,
		Blocks: []internal.SlackBlock{
			{
				Type: "header",
				Text: &internal.SlackTextObject{
					Type: "plain_text",
					Text: "Docs Answer",
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: &internal.SlackTextObject{
					Type: "mrkdwn",
					Text: responseMessage,
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

func rateFeedbackMarkdownPayload(content, question, links string, isPrivate bool, rating internal.MendableRatingScore) ([]byte, error) {
	log.Info().Msgf("Incoming Message: %v", content)

	var (
		responseType    string
		responseMessage string
	)
	if isPrivate {
		responseType = "ephemeral"
	} else {
		responseType = "in_channel"
	}

	if rating == internal.PositiveFeedbackScore {
		responseMessage = internal.DefaultPositiveRatingMessage
	}

	if rating == internal.NegativeFeedbackScore {
		responseMessage = internal.DefaultNegativeRatingMessage
	}

	payload := internal.SlackPayload{
		ResponseType: responseType,
		Blocks: []internal.SlackBlock{
			{
				Type: "header",
				Text: &internal.SlackTextObject{
					Type: "plain_text",
					Text: "Docs Answer",
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
						Text: "*Rate Answer:*",
					},
				},
			}, {
				Type: "section",
				Text: &internal.SlackTextObject{
					Type: "mrkdwn",
					Text: responseMessage,
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

func errorEval(ctx context.Context, e *error, a *SlackActionRequest, isPrivate bool) {
	if e != nil {
		err := internal.ReplyWithErrorMessage(a.action.ResponseURL, isPrivate)
		if err != nil {
			log.Error().Err(err).Msg("Error when attempting to return the model rating feedback answer back to Slack.")
			internal.LogError(err)
			return
		}
	}
}

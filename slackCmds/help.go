// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package slackCmds

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"spectrocloud.com/spectromate/internal"
)

// HelpCmd returns the help Slack command logic and payload.
func HelpCmd() ([]byte, error) {

	markdownContent := "*Commands*\n\nThe following commands are available:\n\n\n- `help` - A summary of all available commands.\n\n- `ask` - Ask a docs related question. Example: `/docs ask how do I enable Prometheus?`\n\n- `pask` - Same as `ask` but with private replies."
	returnPayload, err := helpMarkdownPayload(markdownContent, "Docs Answer")
	if err != nil {
		log.Info().Err(err).Msg("Error creating markdown payload.")
		return nil, err
	}

	return returnPayload, nil
}

// // createMarkdownPayload creates a Slack payload with a markdown block
func helpMarkdownPayload(content, title string) ([]byte, error) {
	log.Info().Msgf("Incoming Message: %v", content)

	payload := internal.SlackPayload{
		ResponseType: "ephemeral",
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

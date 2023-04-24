// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"math/rand"
	"time"
)

const (
	// ApiVersionV1 is the version for the v1 API.
	ApiVersionV1 string = "v1/"
	// ApiPath is the path for the API.
	ApiPath string = "/api/"
	// ApiPrefixV1 is the prefix for the v1 API.
	ApiPrefixV1 string = ApiPath + ApiVersionV1
	// SlackPostMessageURL is the URL for the Slack chat.postMessage API.
	SlackPostMessageURL string = "https://slack.com/api/chat.postMessage"
	// SlackDefaultUserErrorMessage is the default error message for the user.
	SlackDefaultUserErrorMessage string = "An error occured with the help command. Please reach out to `#docs` for assistance."
	// MendableNewConversationURL is the URL for the Mendable new conversation API.
	MendandableNewConversationURL string = "https://api.mendable.ai/v0/newConversation"
	// MendableChatQueryURL is the URL for the Mendable chat query API.
	MendableChatQueryURL string = "https://api.mendable.ai/v0/mendableChat"
	// MendableRatingFeedbackURL is the URL for the Mendable rating feedback API.
	MendableRatingFeedbackURL string = "https://api.mendable.ai/v0/rateMessage"
	// DefaultUserErrorMessage is the default error message for the user.
	DefaultUserErrorMessage string = `:warning: I'm sorry, I'm having technical issues. Notify the docs team @ #docs and please try again later.`
	// DefaultNotFoundResponse is the default response for when no answer is found.
	DefaultNotFoundResponse string = `I'm sorry, I couldn't find an answer to your question. Please provide me feedback and try rephrasing your question.`
	// ActionsAskModelPositiveFeedbackID is the ID for the positive feedback action.
	ActionsAskModelPositiveFeedbackID string = "ask_model_positive_feedback"
	// ActionsAskModelNegativeFeedbackID is the ID for the negative feedback action.
	ActionsAskModelNegativeFeedbackID string = "ask_model_negative_feedback"
	// DefaultCacheExpirationPeriod is the default expiration period for the cache.
	DefaultCacheExpirationPeriod time.Duration = 15 * time.Minute
	// DefaultMendableQueryTimeout is the default timeout for Mendable queries.
	DefaultMendableQueryTimeout time.Duration = 60 * time.Second
	// DefaultPositiveRatingMessage is the default message for positive feedback.
	DefaultPositiveRatingMessage string = `Thank you for providing the :thumbsup: feedback!`
	// DefaultNegativeRatingMessage is the default message for negative feedback.
	DefaultNegativeRatingMessage string = `Thank you for providing the :thumbsdown: feedback!`
	// DefaultNoSourcesIdentifiedMessage is the default message for when no sources are identified.
	DefaultNoSourcesIdentifiedMessage string = `:mag: Unable to identify a specific documentation URL.`
)

// GetRandomWaitMessage returns a random wait message.
func GetRandomWaitMessage() string {
	waitMessages := []string{
		":hourglass_flowing_sand: Hang tight while I review the docs...",
		":hourglass_flowing_sand: Just a moment while I explore the documentation rabbit hole...",
		":hourglass_flowing_sand: Hold tight while I decode this technical jargon for you!",
		":hourglass_flowing_sand: Please be patient as I consult the ancient scrolls of documentation...",
		":hourglass_flowing_sand: Just a little bit longer, we're digging through the treasure trove of documentation to find what you need!",
		":hourglass_flowing_sand: Hold on tight, we're flipping through pages to bring you the answer!",
		":hourglass_flowing_sand: Just a sec, we're diving deep into the ocean of documentation to find the hidden gems.",
		":hourglass_flowing_sand: Please hold while we decipher the clues in the documentation mystery!",
		":hourglass_flowing_sand: Hold tight while we navigate through the labyrinth of documentation and find the solution.",
		":hourglass_flowing_sand: Just a few more seconds... we're almost done unraveling the mysteries of the documentation.",
		":hourglass_flowing_sand: We're almost there, just a little bit longer until we decode the secrets of the documentation universe!",
	}

	return waitMessages[rand.Intn(len(waitMessages))]
}

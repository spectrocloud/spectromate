package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

// sourceValidation validates the source of the request is Slack.
// The Slack signature header value is validated against the signing secret received from Slack
func SourceValidation(ctx context.Context, e events.LambdaFunctionURLRequest, signingSecret string) error {
	log.Info().Msg("Checking the signature")
	payload, err := base64.StdEncoding.DecodeString(e.Body)
	if err != nil {
		return errors.New("internal Error: Failed to decode the payload DB")
	}
	slackTimestamp := e.Headers["x-slack-request-timestamp"]
	slackSignature := e.Headers["x-slack-signature"]
	slackSigningBaseString := fmt.Sprintf("v0:%s:%s", slackTimestamp, string(payload))

	if !validateSignature(ctx, slackSignature, signingSecret, slackSigningBaseString) {
		return errors.New("signature did not match")

	}
	log.Info().Msg("Signature is valid")
	return nil
}

// validateSignature validates the signature of the request
func validateSignature(ctx context.Context, signature, signingSecret, slackSigningBaseString string) bool {

	//calculate SHA256 of the slackSigningBaseString using signingSecret
	mac := hmac.New(sha256.New, []byte(signingSecret))
	mac.Write([]byte(slackSigningBaseString))

	//hex encode the SHA256
	calculatedSignature := fmt.Sprintf("v0=%s", hex.EncodeToString(mac.Sum(nil)))
	return hmac.Equal([]byte(signature), []byte(calculatedSignature))
}

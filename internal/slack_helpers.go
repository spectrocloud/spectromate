package internal

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

// SourceValidation validates the request signature by comparing the signing secret value.
func SourceValidation(ctx context.Context, r *http.Request, signingSecret string) error {
	log.Info().Msg("Checking the signature")
	body, err := getRequestBody(r)
	if err != nil {
		return errors.New("internal Error: Failed to decode the payload DB")
	}
	slackTimestamp := r.Header.Get("X-Slack-Request-Timestamp")
	slackSignature := r.Header.Get("X-Slack-Signature")
	slackSigningBaseString := fmt.Sprintf("v0:%s:%s", slackTimestamp, string(body))

	if !validateSignature(ctx, slackSignature, signingSecret, slackSigningBaseString) {
		return errors.New("signature did not match")
	}

	log.Debug().Msg("Signature is valid")
	return nil
}

// getRequestBody returns the request body as a byte slice.
func getRequestBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
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

// // createMarkdownPayload creates a Slack payload with a markdown block
func CreateMarkdownPayload(content, title string) ([]byte, error) {
	log.Info().Msgf("Incoming Message: %v", content)

	payload := SlackPayload{
		ReponseType: "in_channel",
		Blocks: []SlackBlock{
			{
				Type: "header",
				Text: &SlackTextObject{
					Type: "plain_text",
					Text: title,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: &SlackTextObject{
					Type: "mrkdwn",
					Text: content,
				},
			},
		},
	}

	log.Debug().Msgf("Blocks: %v", payload)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	return payloadBytes, nil
}

// DecodeForm decodes the form data into the struct.
func DecodeForm(values url.Values, dst interface{}) error {
	// Use the reflect package to loop through the fields of the struct
	// and set them based on the values in the form data.
	for fieldName, fieldValue := range values {
		field := reflect.ValueOf(dst).Elem().FieldByName(fieldName)
		if !field.IsValid() || !field.CanSet() {
			return fmt.Errorf("no such field %s", fieldName)
		}
		switch field.Type().Kind() {
		case reflect.String:
			field.SetString(fieldValue[0])
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(fieldValue[0])
			if err != nil {
				return err
			}
			field.SetBool(boolValue)
		// Add cases for other types as needed.
		default:
			return fmt.Errorf("unsupported field type %s", field.Type().Kind())
		}
	}
	return nil
}

// getSlackEvent parses the request and returns a SlackEvent from the request.
func GetSlackEvent(request *http.Request) (SlackEvent, error) {
	var event SlackEvent

	err := request.ParseForm()
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing form.")
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(&event, request.PostForm)
	if err != nil {
		log.Fatal().Err(err).Msg("Error decoding form.")
	}

	return event, nil
}

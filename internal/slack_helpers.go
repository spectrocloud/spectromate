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

// replyStatus200 replies to the Slack event with a 200 status.
// This is required to prevent Slack from considering the request a failure.
// Slack requires a response within 3 seconds.
func ReplyStatus200(reponseURL string, writer http.ResponseWriter, isPrivate bool) error {

	markdownContent := GetRandomWaitMessage()
	_, err := genericMarkdownPayload("Docs Answer", markdownContent, isPrivate)
	if err != nil {
		LogError(err)
		log.Error().Err(err).Msg("Error creating Slack 200 markdown payload.")
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(markdownContent))
	if err != nil {
		LogError(err)
		log.Error().Err(err).Msg("Error writing 200 OK Wait Reply.")
	}

	log.Debug().Msg("Successfully replied to Slack with status 200.")
	return nil
}

// replyWithAnswer replies to the Slack event using the response URL provided.
func ReplyWithAnswer(reponseURL string, payload []byte, isPrivate bool) error {

	client := &http.Client{}
	req, err := http.NewRequest("POST", reponseURL, bytes.NewBuffer(payload))

	if err != nil {
		log.Debug().Err(err).Msg("error creating the reply back HTTP request")
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("error encountered while sending the Slack anser reply back HTTP request")
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	return nil
}

func genericMarkdownPayload(title, content string, isPrivate bool) ([]byte, error) {

	var reponseType string
	if isPrivate {
		reponseType = "ephemeral"
	} else {
		reponseType = "in_channel"
	}

	payload := SlackPayload{
		ReponseType: reponseType,
		Blocks: []SlackBlock{
			{
				Type: "header",
				Text: &SlackTextObject{
					Type: "plain_text",
					Text: title,
				},
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

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	return payloadBytes, nil
}

// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

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
	"time"

	"github.com/avast/retry-go"
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
				LogError(err)
				return err
			}
			field.SetBool(boolValue)
		// Add cases for other types as needed.
		default:
			err := fmt.Errorf("unsupported field type %s", field.Type().Kind())
			LogError(err)
			return err
		}
	}
	return nil
}

func DecodeSlackActionEventForm(payload string) (SlackActionEvent, error) {
	var event SlackActionEvent
	err := json.Unmarshal([]byte(payload), &event)
	if err != nil {
		LogError(err)
		log.Debug().Err(err).Msg("error encountered while decoding the Slack action event form data")
		return event, err
	}
	return event, err
}

// getSlackEvent parses the request and returns a SlackEvent from the request.
func GetSlackActionsEvent(request *http.Request) (SlackActionEvent, error) {
	var actionsEvent SlackActionEvent

	err := request.ParseForm()
	if err != nil {
		LogError(err)
		log.Fatal().Err(err).Msg("Error parsing form.")
	}

	payload := request.PostFormValue("payload")

	actionsEvent, err = DecodeSlackActionEventForm(payload)
	if err != nil {
		LogError(err)
		log.Debug().Err(err).Msg("error encountered while decoding the Slack action event form data")
		return actionsEvent, err
	}

	return actionsEvent, nil
}

// getSlackEvent parses the request and returns a SlackEvent from the request.
func GetSlackEvent(request *http.Request) (SlackEvent, error) {
	var event SlackEvent

	err := request.ParseForm()
	if err != nil {
		LogError(err)
		log.Fatal().Err(err).Msg("Error parsing form.")
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(&event, request.PostForm)
	if err != nil {
		LogError(err)
		log.Fatal().Err(err).Msg("Error decoding form.")
	}

	return event, nil
}

// replyStatus200 replies to the Slack event with a 200 status.
// This is required to prevent Slack from considering the request a failure.
// Slack requires a response within 3 seconds.
func ReplyStatus200(responseURL string, writer http.ResponseWriter, isPrivate bool) ([]byte, error) {

	if responseURL == "" {
		err := errors.New("response URL is empty")
		log.Debug().Err(err).Msg("error encountered while sending the Slack 200 OK reply back HTTP request")
		LogError(err)
		return []byte{}, err
	}

	markdownContent := GetRandomWaitMessage()
	returnValue, err := waitMessagePayload("Docs Answer", markdownContent, isPrivate)
	if err != nil {
		LogError(err)
		log.Error().Err(err).Msg("Error creating Slack 200 markdown payload.")
	}

	log.Debug().Msg("Successfully replied to Slack with status 200.")
	return returnValue, nil
}

// replyWithAnswer replies to the Slack event using the response URL provided.
func ReplyWithAnswer(responseURL string, payload []byte, isPrivate bool) error {

	if responseURL == "" {
		err := errors.New("response URL is empty")
		log.Debug().Err(err).Msg("error encountered while sending the Slack 200 OK reply back HTTP request")
		LogError(err)
		return err
	}

	// Retry the request to Slack up to 3 times.
	// This is to prevent the function from failing if Slack is slow to respond.
	err := retry.Do(
		func() error {
			client := DefaultHTTPClient()
			req, err := http.NewRequest("POST", responseURL, bytes.NewBuffer(payload))
			if err != nil {
				log.Debug().Err(err).Msg("error creating the reply back HTTP request")
				LogError(err)
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				log.Debug().Err(err).Msg("error encountered while sending the Slack answer reply back HTTP request")
				LogError(err)
				return err
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusOK {
				rawError, err := io.ReadAll(res.Body)
				if err != nil {
					log.Debug().Err(err).Msg("unable to read error from the Slack answer reply back HTTP request in an error scenairo")
					LogError(err)
				}
				erroString := string(rawError)
				log.Debug().Msgf("error encountered while sending the Slack answer reply back HTTP request, status code: %d", res.StatusCode)
				err = errors.New(erroString)
				LogError(err)
				return err
			}
			return nil
		}, retry.Attempts(3), retry.Delay(3*time.Second), retry.LastErrorOnly(true),
	)
	if err != nil {
		log.Debug().Err(err).Msg("error encountered while sending the Slack answer reply back HTTP request")
		LogError(err)
		return err
	}

	return nil
}

func ReplyWithErrorMessage(responseURL string, isPrivate bool) error {
	if responseURL == "" {
		err := errors.New("response URL is empty")
		log.Debug().Err(err).Msg("error encountered while sending the Slack 200 OK reply back HTTP request")
		LogError(err)
		return err
	}

	clientMessage, err := errorMessagePayload(DefaultUserErrorMessage, isPrivate)
	if err != nil {
		LogError(err)
		log.Error().Err(err).Msg("error creating the user error message markdown payload.")
		return err
	}

	// Retry the request to Slack up to 3 times.
	// This is to prevent the function from failing if Slack is slow to respond.
	err = retry.Do(
		func() error {
			client := DefaultHTTPClient()
			req, err := http.NewRequest("POST", responseURL, bytes.NewBuffer(clientMessage))
			if err != nil {
				log.Debug().Err(err).Msg("error creating the reply back HTTP request")
				LogError(err)
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				log.Debug().Err(err).Msg("error encountered while sending the Slack answer reply back HTTP request")
				LogError(err)
				return err
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusOK {
				rawError, err := io.ReadAll(res.Body)
				if err != nil {
					log.Debug().Err(err).Msg("unable to read error from the Slack answer reply back HTTP request in an error scenairo")
					LogError(err)
				}
				erroString := string(rawError)
				log.Debug().Msgf("error encountered while sending the Slack answer reply back HTTP request, status code: %d", res.StatusCode)
				err = errors.New(erroString)
				LogError(err)
				return err
			}
			return nil
		}, retry.Attempts(5), retry.Delay(2*time.Second), retry.LastErrorOnly(true),
	)
	if err != nil {
		log.Debug().Err(err).Msg("error encountered while sending the Slack answer reply back HTTP request")
		LogError(err)
		return err
	}

	return nil
}

func waitMessagePayload(title, content string, isPrivate bool) ([]byte, error) {

	payload := SlackPayload{
		ResponseType: "ephemeral",
		Blocks: []SlackBlock{
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
		log.Debug().Err(err).Msg("error marshalling the Slack payload")
		LogError(err)
		return []byte{}, err
	}

	return payloadBytes, nil
}

func errorMessagePayload(content string, isPrivate bool) ([]byte, error) {

	var responseType string
	if isPrivate {
		responseType = "ephemeral"
	} else {
		responseType = "in_channel"
	}

	payload := SlackPayload{
		ResponseType: responseType,
		Blocks: []SlackBlock{
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
		log.Debug().Err(err).Msg("error marshalling the user error message payload")
		LogError(err)
		return []byte{}, err
	}

	return payloadBytes, nil
}

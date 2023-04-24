// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// CreateNewConversation creates a new conversation with Mendable.
func CreateNewConversation(ctx context.Context, apiKey, apiURL string) (int64, error) {

	requestBody := MendableAPIRequest{
		ApiKey: apiKey,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Debug().Err(err).Msg("Error creating new conversation body")
		return -1, err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Debug().Err(err).Msg("error creating new conversation request")
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("an error occured during the HTTP request")
		LogError(err)
		return -1, err
	}
	defer resp.Body.Close()

	var responseBody MendableNewConversationResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		log.Debug().Err(err).Msg("error decoding response body")
		LogError(err)
		return -1, err
	}

	log.Debug().Msgf("Conversation ID: %d", responseBody.ConversationID)

	return responseBody.ConversationID, nil
}

// SendModelRating sends the model rating to Mendable.
func SendModelRating(ctx context.Context, messageID int, rateScore MendableRatingScore, apiKey, apiURL string) error {

	requestBody := MendableModelRatingFeedback{
		ApiKey:    apiKey,
		MessageID: messageID,
		Rating:    rateScore,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Debug().Err(err).Msg("Error marshaling the model rating body")
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Debug().Err(err).Msg("error creating the model rating request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("an error occured during the HTTP request")
		LogError(err)
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Debug().Err(err).Msg("error reading model rating response body")
		LogError(err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Debug().Msgf("error while sending model rating: %s", responseBody)
		log.Debug().Msgf("status code: %d", resp.StatusCode)
		return fmt.Errorf("error while sending model rating: %s", responseBody)
	}

	return nil

}

// sendDocsQuery sends a query to Mendable and returns the response.
func SendDocsQuery(ctx context.Context, query MendableRequestPayload, queryURL string) (MendableQueryResponse, error) {

	var mendableResponse MendableQueryResponse

	log.Debug().Msgf("Query Question: %s", query.Question)

	payload := MendableRequestPayload{
		ApiKey:         query.ApiKey,
		Question:       query.Question,
		History:        query.History,
		ShouldStream:   false,
		ConversationID: query.ConversationID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Debug().Err(err).Msg("Error while marshalling payload:")
		return mendableResponse, err
	}

	client := &http.Client{
		Timeout: DefaultMendableQueryTimeout,
	}

	request, err := http.NewRequest("POST", queryURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Debug().Err(err).Msg("Error while creating POST request:")
		LogError(err)
		return mendableResponse, err
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Debug().Err(err).Msg("Error while making POST request:")
		LogError(err)
		return mendableResponse, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Error while reading response body:")
		LogError(err)
		return mendableResponse, err
	}

	var result MendablePayload
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Debug().Err(err).Msg("Error while unmarshalling response:")
		LogError(err)
		return mendableResponse, err
	}

	uniqueLinks := retrieveUniqueLinks(&result.Sources)

	mendableResponse = MendableQueryResponse{
		ConversationID: int64(query.ConversationID),
		MessageID:      fmt.Sprint(result.MessageID),
		Question:       query.Question,
		Answer:         result.Answer.Text,
		Links:          uniqueLinks,
	}

	if mendableResponse.Answer == "" {
		// Set default response if the Answer field is empty
		mendableResponse.Answer = DefaultNotFoundResponse
	}

	log.Debug().Msgf("Mendable Question response: %v", mendableResponse.Answer)

	return mendableResponse, nil

}

// retrieveUniqueLinks returns a slice of unique links from the Mendable sources.
func retrieveUniqueLinks(list *[]MendableSources) []string {
	uniqueLinks := make([]string, 0)
	reviewedLinksMap := make(map[string]bool)

	for _, source := range *list {

		link := source.Link

		// Skip links that are equal to "https://docs.spectrocloud.com"
		if link == PublicDocumentationURL {
			continue
		}

		// Check if the link already exists in the map
		if _, exists := reviewedLinksMap[link]; exists {
			continue
		}

		// Add the link to the reviewedLinksMap and to the unique links slice
		reviewedLinksMap[link] = true
		uniqueLinks = append(uniqueLinks, link)
	}

	return uniqueLinks

}

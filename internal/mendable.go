package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
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

// sendDocsQuery sends a query to Mendable and returns the response.
func SendDocsQuery(ctx context.Context, query MendableRequestPayload, queryURL string) (MendableQueryResponse, error) {

	var mendableResponse MendableQueryResponse

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
		Timeout: 15 * time.Second,
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

	body, err := ioutil.ReadAll(response.Body)
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
		Question: query.Question,
		Answer:   result.Answer.Text,
		Links:    uniqueLinks,
	}

	return mendableResponse, nil

}

// retrieveUniqueLinks returns a slice of unique links from the Mendable sources.
func retrieveUniqueLinks(list *[]MendableSources) []string {
	uniqueLinks := make([]string, 0)
	reviewedLinksMap := make(map[string]bool)

	for _, source := range *list {

		link := source.Link

		// Skip links that are equal to "https://docs.spectrocloud.com"
		if link == "https://docs.spectrocloud.com" {
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

package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func CreateNewConversation(apiKey, apiURL string) (int64, error) {
	requestBody := map[string]string{
		"api_key": apiKey,
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
	req.Header.Set("api_key", apiKey)

	client := &http.Client{
		Timeout: 10 * 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("an error occured during the HTTP request")
		return -1, err
	}
	defer resp.Body.Close()

	var responseBody MendableNewConversationResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		log.Debug().Err(err).Msg("error decoding response body")
	}

	return responseBody.ConversationID, nil
}

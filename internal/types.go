package internal

import "time"

/*
 * Mendable API types
 */
type MendableAPIRequest struct {
	ApiKey string `json:"api_key"`
}

type MendableRequestPayload struct {
	ApiKey         string         `json:"api_key"`
	Question       string         `json:"question"`
	History        []HistoryItems `json:"history"`
	ShouldStream   bool           `json:"shouldStream"`
	ConversationID int64          `json:"conversation_id"`
}

type HistoryItems struct {
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
}

type MendablePayload struct {
	Answer struct {
		Text string `json:"text"`
	} `json:"answer"`
	MessageID int               `json:"message_id"`
	Sources   []MendableSources `json:"sources"`
}

type MendableSources struct {
	ID      int     `json:"id"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
	Date    *string `json:"date"`
	Link    string  `json:"link"`
}

type MendableQueryResponse struct {
	ConversationID int64
	Question       string
	Answer         string
	Links          []string
}

/*
 * Slack API types
 */

type SlackEvent struct {
	Token               string `schema:"token"`
	TeamID              string `schema:"team_id"`
	TeamDomain          string `schema:"team_domain"`
	ChannelID           string `schema:"channel_id"`
	ChannelName         string `schema:"channel_name"`
	UserID              string `schema:"user_id"`
	UserName            string `schema:"user_name"`
	Command             string `schema:"command"`
	Text                string `schema:"text"`
	ResponseURL         string `schema:"response_url"`
	TriggerID           string `schema:"trigger_id"`
	APIAppID            string `schema:"api_app_id"`
	IsEnterpriseInstall bool   `schema:"is_enterprise_install"`
}

type SlackPayload struct {
	ReponseType string       `json:"response_type,omitempty"`
	Blocks      []SlackBlock `json:"blocks"`
}

type SlackBlock struct {
	Type   string            `json:"type"`
	Fields []SlackTextObject `json:"fields,omitempty"`
	Text   *SlackTextObject  `json:"text,omitempty"`
}

type SlackTextObject struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

/*

* Cache types

 */

type CacheItem struct {
	UserID         string
	ChannelID      string
	ConversationID string
	Question       string
	Answer         string
	Timestamp      *time.Time
	Counter        string
}

type MendableNewConversationResponse struct {
	ConversationID int64 `json:"conversation_id"`
}

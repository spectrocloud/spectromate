package internal

import "time"

type Response struct {
	Chunk          string `json:"chunk,omitempty"`
	Metadata       MendableMetadata
	FullDocs       bool   `json:"full_docs,omitempty"`
	AlertMessage   string `json:"alert_message,omitempty"`
	FullDocsSource string `json:"full_docs_source,omitempty"`
}

type MendableResponse struct {
	Chunk    string             `json:"chunk"`
	Metadata []MendableMetadata `json:"metadata"`
}

type MendableMetadata struct {
	ID      string    `json:"id"`
	Content string    `json:"content"`
	Score   float64   `json:"score"`
	Date    time.Time `json:"date"`
	Link    string    `json:"link"`
}

type MendableAPIRequest struct {
	ApiKey string `json:"api_key"`
}

type HistoryItems struct {
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
}

type MendableQueryPayload struct {
	ApiKey         string         `json:"api_key"`
	Question       string         `json:"question,omitempty"`
	History        []HistoryItems `json:"history,omitempty"`
	ConversationID int            `json:"conversation_id,omitempty"`
}

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

type CacheItem struct {
	UserID         string
	ChannelID      string
	ConversationID string
	Timestamp      *time.Time
}

type MendableNewConversationResponse struct {
	ConversationID int64 `json:"conversation_id"`
}

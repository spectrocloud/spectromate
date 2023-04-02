package internal

import "time"

type MendableResponse struct {
	Chunk    string             `json:"chunk"`
	Metadata []MendableMetadata `json:"metadata"`
}

type MendableMetadata struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Link    string `json:"link"`
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

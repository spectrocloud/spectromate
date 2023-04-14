package internal

import (
	"time"
)

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
	MessageID      string
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
	ResponseType    string       `json:"response_type,omitempty"`
	DeleteOriginal  bool         `json:"delete_original,omitempty"`
	ReplaceOriginal bool         `json:"replace_original,omitempty"`
	Blocks          []SlackBlock `json:"blocks"`
}

type SlackBlock struct {
	Type     string            `json:"type"`
	Fields   []SlackTextObject `json:"fields,omitempty"`
	Elements []SlackElements   `json:"elements,omitempty"`
	Text     *SlackTextObject  `json:"text,omitempty"`
}

type SlackTextObject struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SlackElementText struct {
	Type  string `json:"type"`
	Emoji bool   `json:"emoji"`
	Text  string `json:"text"`
}
type SlackElements struct {
	Type     string           `json:"type"`
	Text     SlackElementText `json:"text"`
	Style    string           `json:"style,omitempty"`
	Value    string           `json:"value"`
	ActionID string           `json:"action_id"`
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
	History        []HistoryItems
	Counter        string
}

type MendableNewConversationResponse struct {
	ConversationID int64 `json:"conversation_id"`
}

/*
* Slack Actions types

 */

type SlackActionEvent struct {
	Type                string       `json:"type"`
	User                SlackUser    `json:"user"`
	ApiAppID            string       `json:"api_app_id"`
	Token               string       `json:"token"`
	Container           Container    `json:"container"`
	TriggerID           string       `json:"trigger_id"`
	Team                Team         `json:"team"`
	Enterprise          interface{}  `json:"enterprise"`
	IsEnterpriseInstall bool         `json:"is_enterprise_install"`
	Channel             Channel      `json:"channel"`
	Message             SlackMessage `json:"message"`
	State               State        `json:"state"`
	ResponseURL         string       `json:"response_url"`
	Actions             []Action     `json:"actions"`
}

type SlackUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	TeamID   string `json:"team_id"`
}

type Container struct {
	Type        string `json:"type"`
	MessageTS   string `json:"message_ts"`
	ChannelID   string `json:"channel_id"`
	IsEphemeral bool   `json:"is_ephemeral"`
}

type Team struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackMessage struct {
	Type    string              `json:"type"`
	Subtype string              `json:"subtype"`
	Text    string              `json:"text"`
	Ts      string              `json:"ts"`
	BotID   string              `json:"bot_id"`
	Blocks  []SlackActionsBlock `json:"blocks"`
}

type SlackActionsBlock struct {
	Type     string      `json:"type"`
	BlockID  string      `json:"block_id"`
	Text     ActionsText `json:"text,omitempty"`
	Elements []Action    `json:"elements,omitempty"`
}

type ActionsText struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Emoji    bool   `json:"emoji,omitempty"`
	Verbatim bool   `json:"verbatim,omitempty"`
}

type State struct {
	Values map[string]interface{} `json:"values"`
}

type Action struct {
	ActionID string      `json:"action_id"`
	BlockID  string      `json:"block_id"`
	Type     string      `json:"type"`
	Text     ActionsText `json:"text,omitempty"`
	Value    string      `json:"value,omitempty"`
	Style    string      `json:"style,omitempty"`
	ActionTS string      `json:"action_ts,omitempty"`
}

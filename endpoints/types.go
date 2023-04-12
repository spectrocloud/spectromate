package endpoints

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"spectrocloud.com/docs-slack-bot/internal"
)

type HealthRoute struct {
	ctx context.Context
}

type SlackRoute struct {
	ctx            context.Context
	signingSecret  string
	mendableApiKey string
	SlackEvent     *internal.SlackEvent
	cache          *redis.Client
}

type ActionsRoute struct {
	ctx            context.Context
	signingSecret  string
	mendableApiKey string
	ActionsEvent   *internal.SlackActionEvent
}

type SlackCommands int

const (
	Help SlackCommands = iota
	Ask
	PAsk
)

// String converts a SlackCommands type to a string.
func (s SlackCommands) String() string {
	switch s {
	case Help:
		return "help"
	case Ask:
		return "ask"
	case PAsk:
		return "pask"
	default:
		return "unknown"
	}
}

// SlackCommandsFromString converts a string to a SlackCommands type.
func SlackCommandsFromString(s string) (SlackCommands, error) {
	switch s {
	case "help":
		return Help, nil
	case "ask":
		return Ask, nil
	case "pask":
		return PAsk, nil
	default:
		return -1, fmt.Errorf("unknown command: %s", s)
	}
}

// SlackCommandsIntf is an interface for the SlackCommands type.
// This is used to pass the SlackRoute struct to the SlackCommands functions.
// And avoid circular dependencies.
type SlackCommandsIntf interface {
	GetContext() context.Context
	GetSigningSecret() string
	GetMendableAPIKey() string
	GetSlackEvent() *internal.SlackEvent
	GetCache() *redis.Client
}

func (sr *SlackRoute) GetContext() context.Context {
	return sr.ctx
}

func (sr *SlackRoute) GetSigningSecret() string {
	return sr.signingSecret
}

func (sr *SlackRoute) GetMendableAPIKey() string {
	return sr.mendableApiKey
}

func (sr *SlackRoute) GetSlackEvent() *internal.SlackEvent {
	return sr.SlackEvent
}

func (sr *SlackRoute) GetCache() *redis.Client {
	return sr.cache
}

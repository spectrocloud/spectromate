package endpoints

import (
	"context"
	"fmt"

	"spectrocloud.com/spectromate/internal"
)

type HealthRoute struct {
	ctx context.Context
}

type SlackRoute struct {
	ctx            context.Context
	signingSecret  string
	mendableApiKey string
	SlackEvent     *internal.SlackEvent
	cache          internal.Cache
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

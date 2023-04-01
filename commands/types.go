package commands

import "context"

type HealthRoute struct {
	ctx context.Context
}

type HelpRoute struct {
	ctx           context.Context
	signingSecret string
}

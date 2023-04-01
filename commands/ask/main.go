package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, e events.LambdaFunctionURLRequest) (string, error) {

	signingSecret, ok := os.LookupEnv("SLACK_SIGNING_SECRET_ID")
	if !ok {
		log.Fatal().Msg(`could not find variable "SLACK_SIGNING_SECRET_ID"`)
	}
	err := SourceValidation(ctx, e, signingSecret)
	if err != nil {
		return err.Error(), err
	}

	return "OK", nil
}

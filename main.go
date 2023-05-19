// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	_ "go.uber.org/automaxprocs"
	"spectrocloud.com/spectromate/endpoints"
	"spectrocloud.com/spectromate/internal"
)

var (
	globalRedisPort      int64
	globalRedisURL       string
	globalRedisClient    internal.Cache
	globalRedisPassword  string
	globalRedisUser      string
	globalRedisTLS       string
	globalTraceLevel     string
	globalHost           string
	globalPort           string
	globalHostURL        string = globalHost + ":" + globalPort
	globalSigningSecret  string
	globalMendableAPIKey string
	Version              string
)

func init() {
	globalTraceLevel = strings.ToUpper(internal.Getenv("TRACE", "INFO"))
	internal.InitLogger(globalTraceLevel)
	globalSigningSecret = internal.Getenv("SLACK_SIGNING_SECRET", "")
	globalMendableAPIKey = internal.Getenv("MENDABLE_API_KEY", "")
	globalRedisTLS = strings.ToLower(internal.Getenv("REDIS_TLS", "false"))
	redisTLS := globalRedisTLS
	port := internal.Getenv("PORT", "3000")
	host := internal.Getenv("HOST", "0.0.0.0")
	globalHost = host
	globalPort = port
	globalHostURL = host + ":" + port
	globalRedisPort = internal.StringToInt64(internal.Getenv("REDIS_PORT", "6379"))
	globalRedisURL = internal.Getenv("REDIS_URL", "localhost")
	globalRedisPassword = internal.Getenv("REDIS_PASSWORD", "")
	globalRedisUser = internal.Getenv("REDIS_USER", "")
	reditConnectionString := fmt.Sprintf("%s:%d", globalRedisURL, globalRedisPort)

	if globalSigningSecret == "" {
		log.Fatal().Msg("The required environment variable SLACK_SIGNING_SECRET is not set. Exiting...")
	}

	if globalMendableAPIKey == "" {
		log.Fatal().Msg("The required environment variable MENDABLE_API_KEY is not set. Exiting...")
	}

	var tlsConfig *tls.Config

	if redisTLS == "true" {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := internal.NewCache(reditConnectionString,
		globalRedisUser,
		globalRedisPassword,
		tlsConfig)
	log.Debug().Msg("Checking database connection...")
	err := rdb.Ping()
	if err != nil {
		log.Debug().Msg("Redis is not available")
		log.Fatal().Err(err).Msg("Error connecting to redis")
	}
	globalRedisClient = rdb
}

func main() {
	ctx := context.Background()
	rdb := globalRedisClient
	healthRoute := endpoints.NewHealthHandlerContext(ctx, Version)
	slackRoute := endpoints.NewSlackHandlerContext(ctx, globalSigningSecret, globalMendableAPIKey, rdb, Version)
	slackActionsRoute := endpoints.NewActionsHandlerContext(ctx, globalSigningSecret, globalMendableAPIKey, Version)

	http.HandleFunc(internal.ApiPrefixV1+"health", healthRoute.HealthHTTPHandler)
	http.HandleFunc(internal.ApiPrefixV1+"slack", slackRoute.SlackHTTPHandler)
	http.HandleFunc(internal.ApiPrefixV1+"slack/actions", slackActionsRoute.ActionsHTTPHandler)

	log.Info().Msgf("Server is configured for port %s and listing on %s", globalPort, globalHostURL)
	log.Info().Msgf("API Server version:  %s", Version)
	log.Info().Msgf("Redis is configured for %s:%d", globalRedisURL, globalRedisPort)
	log.Info().Msgf("Trace level set to: %s", globalTraceLevel)
	log.Info().Msg("Starting server...")
	http.DefaultClient = internal.DefaultHTTPClient()
	err := http.ListenAndServe(globalHostURL, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("There's an error with the server")
	}

}

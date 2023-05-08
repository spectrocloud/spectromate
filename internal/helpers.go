// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"net/http"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

// getEnv returns the value of the environment variable or a default value
func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}

func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error converting %s string to int64", s)
	}
	return i
}

// DefaultHTTPClient returns a default HTTP client with proxy support.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
}

// SetUserAgent sets the User-Agent header on the request.
func SetUserAgent(version string) string {

	return DefaultUserAgent + version
}

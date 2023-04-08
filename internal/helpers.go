package internal

import (
	"errors"
	"os"
	"strconv"
	"strings"

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

// CheckRequiredEnvVars checks that the required environment variables are set.
func CheckRequiredEnvVars(envVars map[string]string) error {
	var missingVars []string

	for name, defaultValue := range envVars {
		value, exists := os.LookupEnv(name)
		if !exists || strings.TrimSpace(value) == "" {
			missingVars = append(missingVars, name)
		} else {
			envVars[name] = value
		}
		if defaultValue != "" && !exists {
			os.Setenv(name, defaultValue)
		}
	}

	if len(missingVars) > 0 {
		errorMessage := "Missing required environment variables: " + strings.Join(missingVars, ", ")
		return errors.New(errorMessage)
	}

	return nil
}

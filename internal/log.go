package internal

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Use Unix time in every log message.
// Only write logs with a level greater than or equal to logLevel.
func InitLogger(logLevel string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	switch logLevel {
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// If the error is not nil, panic.
// Only use this in init().
func PanicErr(err error) {
	if err != nil {
		log.Fatal().Send()
	}
}

// If the ok value is false, panic.
// If the error is not nil, panic.
// Only use this in init().
func PanicOkErr(ok bool, err error) {
	if err != nil || !ok {
		log.Fatal().Send()
	}
}

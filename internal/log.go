// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"runtime"

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

func LogError(err error) {
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		log.Err(err).
			Str("file", file).
			Int("line", line).
			Str("function", runtime.FuncForPC(pc).Name()).
			Send()
	} else {
		log.Err(err).Send()
	}
}

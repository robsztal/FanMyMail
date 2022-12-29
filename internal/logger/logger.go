package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HumanLogs bool   `env:"HUMAN_LOGS,default=true"`
	LogLevel  string `env:"LOG_LEVEL,default=debug"`
}

// Configure setups Zerolog to print pretty logs if HUMAN_LOGS is set to 1/T/true (default is false).
// Also caller (go file name) is added to log statements
func Configure(config Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldUnit = time.Millisecond

	if config.HumanLogs {
		output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		log.Logger = log.Output(output).With().Caller().Logger()
	} else {
		log.Logger = log.Logger.With().Caller().Logger()
	}

	logLevel, err := zerolog.ParseLevel(strings.ToLower(config.LogLevel))
	if err != nil {
		log.Error().Err(err).Msg("Error parsing provided log level")
	}
	zerolog.SetGlobalLevel(logLevel)
}

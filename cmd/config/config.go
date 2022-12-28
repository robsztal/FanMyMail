package config

import (
	"github.com/Netflix/go-env"
	"github.com/robsztal/FanMyMail/internal/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config of the service
type Config struct {
	Version      string `env:"COMMIT,default=unknown"`
	ServiceName  string `env:"SERVICE_NAME,default=FanMyMail"`
	LoggerConfig logger.Config
	// Add your configuration keys
}

// GetEnvConfig returns Config struct with all the parameters found or their default values
func GetEnvConfig() (Config, error) {
	var config Config
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Error().Err(err).Msg("Reading configuration from env returned an error")
		return Config{}, err
	}
	return config, nil
}

// MarshalZerologObject for logging config
func (conf Config) MarshalZerologObject(e *zerolog.Event) {
	e.Str("version", conf.Version)
	e.Str("serviceName", conf.ServiceName)
	e.Bool("logger.humanLogs", conf.LoggerConfig.HumanLogs)
	e.Str("logger.logLevel", conf.LoggerConfig.LogLevel)
}

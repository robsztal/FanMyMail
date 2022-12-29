package config

import (
	"os"

	"github.com/Netflix/go-env"
	"github.com/robsztal/FanMyMail/internal/logger"
	"github.com/robsztal/FanMyMail/internal/rest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/gmail/v1"
)

// Config of the service
type Config struct {
	Version      string `env:"COMMIT,default=unknown"`
	ServiceName  string `env:"SERVICE_NAME,default=FanMyMail"`
	LoggerConfig logger.Config
	OAuthCfg     *oauth2.Config
	REST         rest.Config
	Token        *oauth2.Token
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

func (conf *Config) InitGoogleConfig() {
	b, err := os.ReadFile("./credentials.json")
	if err != nil {
		log.Panic().Err(err).Msg("failed to load credentials")
	}
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatal().Msg("cannot load credentials")
	}
	conf.OAuthCfg = config
	conf.OAuthCfg.Scopes = []string{gmail.MailGoogleComScope}
	conf.OAuthCfg.RedirectURL = "http://localhost:8080/callback-gl"
}

// MarshalZerologObject for logging config
func (conf *Config) MarshalZerologObject(e *zerolog.Event) {
	e.Str("version", conf.Version)
	e.Str("serviceName", conf.ServiceName)
	e.Bool("logger.humanLogs", conf.LoggerConfig.HumanLogs)
	e.Str("logger.logLevel", conf.LoggerConfig.LogLevel)
}

package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/robsztal/FanMyMail/cmd/config"
	"github.com/robsztal/FanMyMail/internal/logger"
)

func main() {
	cfg, err := config.GetEnvConfig()
	if err != nil {
		log.Fatal().Msg("failed to load config")

	}
	logger.Configure(cfg.LoggerConfig)
	ctx, cancel := context.WithCancel(context.Background())
	ctx = log.Logger.WithContext(ctx)
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	log.Ctx(ctx).Info().Msg("Blowing Fan...")
}

package main

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/robsztal/FanMyMail/cmd/config"
	"github.com/robsztal/FanMyMail/internal/gmail"
	"github.com/robsztal/FanMyMail/internal/logger"
	"github.com/robsztal/FanMyMail/internal/rest"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.GetEnvConfig()
	if err != nil {
		log.Fatal().Msg("failed to load config")

	}
	cfg.InitGoogleConfig()
	logger.Configure(cfg.LoggerConfig)
	ctx, cancel := context.WithCancel(context.Background())
	ctx = log.Logger.WithContext(ctx)
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	router := mux.NewRouter()
	h := gmail.NewHandlers(cfg)
	handlers(router, h)

	server := rest.New(cfg.REST, router)
	server.Serve(ctx)

	svc, err := gmail.NewClient(ctx, cfg)
	if err != nil {
		log.Fatal().Msg("failed to create gmail client")
		return
	}
	fetch, err := svc.Fetch(ctx, "")
	if err != nil {
		return
	}
	log.Ctx(ctx).Info().Interface("fetch", fetch).Msg("Fetched")
	log.Ctx(ctx).Info().Msg("Blowing Fan...")
}

func handlers(router *mux.Router, handlers gmail.Handlers) {
	router.HandleFunc("/", handlers.HandleMain)
	router.HandleFunc("/login-gl", handlers.HandleGoogleLogin)
	router.HandleFunc("/callback-gl", handlers.HandleCallBackFromGoogle)
}

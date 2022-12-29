package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Server interface {
	// Serve runs HTTP server
	Serve(ctx context.Context)
	// Stop stops HTTP server
	Stop(ctx context.Context)
}

type Config struct {
	Port string `env:"HTTP_PORT,default=8080"`
	// DefaultReadTimeout applies to ReadTimeout and ReadHeaderTimeout
	DefaultReadTimeout time.Duration `env:"HTTP_DEFAULT_READ_TIMEOUT,default=10s"`
}

type serverContext struct {
	httpServer *http.Server
}

// New creates Server
func New(conf Config, handler http.Handler) Server {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", conf.Port),
		Handler:           handler,
		ReadTimeout:       conf.DefaultReadTimeout,
		ReadHeaderTimeout: conf.DefaultReadTimeout,
	}

	c := serverContext{
		httpServer: srv,
	}

	return &c
}

// Serve runs HTTP server
func (c *serverContext) Serve(ctx context.Context) {
	log.Ctx(ctx).Info().Msgf("Starting HTTP server at port %s", c.httpServer.Addr)
	go func() {
		err := c.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Ctx(ctx).Error().Err(err).Msg("HTTP server stopped.")
		}
		if errors.Is(err, http.ErrServerClosed) {
			log.Ctx(ctx).Info().Err(err).Msg("HTTP server stopped.")
		}
	}()
}

// Stop stops HTTP server
func (c *serverContext) Stop(ctx context.Context) {
	log.Ctx(ctx).Info().Msg("Stopping HTTP server...")
	ctxShutdown, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := c.httpServer.Shutdown(ctxShutdown); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("An error occurred while shutting down the HTTP server.")
	}
}

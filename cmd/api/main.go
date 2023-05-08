package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"bybarcode/internal/api"
	"bybarcode/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msgf(".env file load error: %s", err)
		return
	}

	cfg := config.NewApiConfig()
	apiApp := api.NewAppApi(cfg, logger)

	go func() {
		if err := apiApp.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Msgf("Api starting error: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Stopping api...")

	if err := apiApp.ShoutDown(ctx); err != nil {
		logger.Fatal().Msgf("Api stopping error: %s", err.Error())
	}

	logger.Info().Msg("Api stopped")
}

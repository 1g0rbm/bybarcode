package main

import (
	"bybarcode/internal/listener"
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
	"bybarcode/internal/db"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msgf(".env file load error: %s", err)
		return
	}

	cfg := config.NewApiConfig()

	conn, err := db.NewConnect("pgx", cfg.DBDsn)
	defer func(conn *db.Connect) {
		err = conn.Close()
	}(&conn)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	l := listener.NewEventListener(&conn)

	apiApp := api.NewAppApi(conn, cfg, logger, l)

	go func() {
		if err = apiApp.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Msgf("Api starting error: %s", err.Error())
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := l.Listen(ctx); err != nil {
			logger.Fatal().Msg(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Stopping api...")

	if err = apiApp.ShoutDown(ctx); err != nil {
		logger.Fatal().Msgf("Api stopping error: %s", err.Error())
	}

	logger.Info().Msg("Api stopped")
}

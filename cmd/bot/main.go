package main

import (
	"os"
	"os/signal"
	"syscall"

	"bybarcode/internal/bot"
	"bybarcode/internal/config"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msgf(".env file load error: %s", err)
		return
	}

	cfg := config.NewBotConfig()
	app := bot.NewAppBot(logger, cfg)

	go func() {
		if err := app.Run(); err != nil {
			logger.Fatal().Msgf("Bot starting error: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Stopping bot...")

	app.Shutdown()

	logger.Info().Msg("Bot was stopped.")
}

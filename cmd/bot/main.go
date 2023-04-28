package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"

	"bybarcode/internal/bot"
	"bybarcode/internal/config"
	"bybarcode/internal/db"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msgf(".env file load error: %s", err)
		return
	}

	cfg := config.NewBotConfig()

	conn, err := db.NewConnect("pgx", cfg.DBDsn)
	defer func(conn *db.Connect) {
		err = conn.Close()
	}(&conn)

	botApi, err := tgbotapi.NewBotAPI(cfg.Token)

	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	app := bot.NewAppBot(botApi, logger, cfg, conn)

	if err = app.Run(); err != nil {
		logger.Fatal().Msg(err.Error())
	}
}

package main

import (
	"bybarcode/internal/db"
	"bybarcode/internal/message"
	"context"
	"github.com/joho/godotenv"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"

	"bybarcode/internal/config"
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

	bot, err := tgbotapi.NewBotAPI(cfg.Token)

	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	for upd := range updates {
		if upd.Message == nil {
			continue
		}

		chat := upd.Message.Chat
		err := conn.CreateAccountIfNotExist(ctx, int(chat.ID), chat.UserName, chat.FirstName, chat.LastName)
		if err != nil {
			logger.Error().Msg(err.Error())
			continue
		}

		if upd.Message.IsCommand() {
			switch upd.Message.Command() {
			case "start":
				_, err := bot.Send(message.OnStartMessage(chat.ID))
				if err != nil {
					logger.Fatal().Msg(err.Error())
				}
			}
		}
	}
}

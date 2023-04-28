package bot

import (
	"context"
	"net/url"

	"github.com/rs/zerolog"
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"

	"bybarcode/internal/config"
	"bybarcode/internal/db"
	"bybarcode/internal/message"
)

type AppBot struct {
	bot    *tgbotapi.BotAPI
	logger zerolog.Logger
	cfg    *config.BotConfig
	db     db.Connect
}

func NewAppBot(bot *tgbotapi.BotAPI, logger zerolog.Logger, cfg *config.BotConfig, db db.Connect) *AppBot {
	return &AppBot{
		bot:    bot,
		logger: logger,
		cfg:    cfg,
		db:     db,
	}
}

func (ab AppBot) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updates := ab.bot.GetUpdatesChan(u)

	for upd := range updates {
		if upd.Message == nil {
			continue
		}

		chat := upd.Message.Chat
		err := ab.db.CreateAccountIfNotExist(ctx, int(chat.ID), chat.UserName, chat.FirstName, chat.LastName)
		if err != nil {
			ab.logger.Error().Msg(err.Error())
			continue
		}

		if upd.Message.IsCommand() {
			switch upd.Message.Command() {
			case "start":
				ab.errorHandler(upd.Message, ab.onStartHandler(upd.Message))
			case "open":
				ab.errorHandler(upd.Message, ab.onOpenHandler(upd.Message))
			}
		}
	}

	return nil
}

func (ab AppBot) onStartHandler(msg *tgbotapi.Message) error {
	response := tgbotapi.NewMessage(msg.Chat.ID, message.OnStartMessage())

	_, err := ab.bot.Send(response)

	return err
}

func (ab AppBot) onOpenHandler(msg *tgbotapi.Message) error {
	webAppURL, err := url.Parse(ab.cfg.TgWebAppUrl)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session, err := ab.db.CreateSession(ctx, msg.Chat.ID)
	if err != nil {
		return err
	}

	webAppURL.Query().Add("token", session.Token.String())
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp(
				"Web App",
				tgbotapi.WebAppInfo{
					URL: webAppURL.String(),
				}),
		),
	)

	response := tgbotapi.NewMessage(msg.Chat.ID, "Нажми на кнопку, чтобы открыть приложение.")
	response.ReplyMarkup = inlineKeyboard

	_, err = ab.bot.Send(response)

	return err
}

func (ab AppBot) errorHandler(msg *tgbotapi.Message, err error) {
	if err == nil {
		return
	}

	response := tgbotapi.NewMessage(
		msg.Chat.ID,
		"Упс! При обработке сообщения возникла ошибка, попробуй позжк :(",
	)

	ab.logger.Error().Msg(err.Error())

	_, err = ab.bot.Send(response)
	if err != nil {
		ab.logger.Fatal().Msg(err.Error())
	}
}
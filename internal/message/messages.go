package message

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func OnStartMessage(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Hi! Soon here will be an app to manage your shopping list.")
}

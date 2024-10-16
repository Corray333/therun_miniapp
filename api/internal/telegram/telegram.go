package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramClient struct {
	bot *tgbotapi.BotAPI
}

func (t *TelegramClient) GetBot() *tgbotapi.BotAPI {
	return t.bot
}

func NewClient(token string) *TelegramClient {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Failed to create bot: ", err)
	}

	bot.Debug = true

	return &TelegramClient{
		bot: nil,
	}
}

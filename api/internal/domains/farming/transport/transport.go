package transport

import (
	"github.com/go-chi/chi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FarmingTransport struct {
	router *chi.Mux
	bot    *tgbotapi.BotAPI
}

type service interface {
	ClaimTokens(userID int64) (pointsGot, pointsBalnce, farmingTime, maxPoints, farmingFrom int)
}

func New(router *chi.Mux, bot *tgbotapi.BotAPI, service service) *FarmingTransport {

	return &FarmingTransport{
		router: router,
		bot:    bot,
	}
}

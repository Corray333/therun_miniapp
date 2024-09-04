package farming

import (
	"github.com/Corray333/therun_miniapp/internal/domains/farming/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/farming/service"
	"github.com/Corray333/therun_miniapp/internal/domains/farming/transport"
)

type FarmingController struct {
	repo      repository.FarmingRepository
	service   service.FarmingService
	transport transport.FarmingTransport
}

// func NewFarmingController(router *chi.Mux, tg *tgbotapi.BotAPI, store storage.Storage) *FarmingController {
// 	repo := repository.New(store)
// 	service := service.New(repo)
// 	transport := transport.New(router, tg, service)

// 	return &FarmingController{
// 		repo:      *repo,
// 		service:   *service,
// 		transport: *transport,
// 	}
// }

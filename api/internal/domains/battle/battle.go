package battle

import (
	"github.com/Corray333/therun_miniapp/internal/domains/battle/external"
	"github.com/Corray333/therun_miniapp/internal/domains/battle/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/battle/service"
	"github.com/Corray333/therun_miniapp/internal/domains/battle/transport"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
)

type BattleController struct {
	repo      repository.BattleRepository
	service   service.BattleService
	transport transport.BattleTransport
}

func NewBattleController(router *chi.Mux, store *storage.Storage) *BattleController {
	repo := repository.New(store)
	ext := external.New()
	service := service.New(repo, ext)
	transport := transport.New(router, service)

	return &BattleController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *BattleController) Build() {
	c.transport.RegisterRoutes()
	go c.service.RunRounds()
}
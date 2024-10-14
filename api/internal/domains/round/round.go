package round

import (
	battle_service "github.com/Corray333/therun_miniapp/internal/domains/battle/service"
	"github.com/Corray333/therun_miniapp/internal/domains/round/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/round/service"
	"github.com/Corray333/therun_miniapp/internal/domains/round/transport"
	user_service "github.com/Corray333/therun_miniapp/internal/domains/user/service"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
)

type RoundController struct {
	repo      repository.RoundRepository
	service   service.RoundService
	transport transport.RoundTransport
}

func NewRoundController(router *chi.Mux, store *storage.Storage, userService *user_service.UserService, battleService *battle_service.BattleService) *RoundController {
	repo := repository.New(store)
	service := service.New(repo, userService, battleService)
	transport := transport.New(router, service)

	return &RoundController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *RoundController) Build() {
	c.transport.RegisterRoutes()
	go c.service.RunRounds()
}

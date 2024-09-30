package battle

import (
	"github.com/Corray333/therun_miniapp/internal/domains/battle/external"
	"github.com/Corray333/therun_miniapp/internal/domains/battle/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/battle/service"
	"github.com/Corray333/therun_miniapp/internal/domains/battle/transport"
	user_service "github.com/Corray333/therun_miniapp/internal/domains/user/service"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
)

type BattleController struct {
	repo      repository.BattleRepository
	service   service.BattleService
	transport transport.BattleTransport
}

func NewBattleController(router *chi.Mux, store *storage.Storage, userService *user_service.UserService) *BattleController {
	repo := repository.New(store)
	ext := external.New()
	service := service.New(repo, ext, userService)
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

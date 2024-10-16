package farming

import (
	"github.com/Corray333/therun_miniapp/internal/domains/farming/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/farming/service"
	"github.com/Corray333/therun_miniapp/internal/domains/farming/transport"
	user_service "github.com/Corray333/therun_miniapp/internal/domains/user/service"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
)

type FarmingController struct {
	repo      repository.FarmingRepository
	service   service.FarmingService
	transport transport.FarmingTransport
}

func NewFarmingController(router *chi.Mux, store *storage.Storage, userService *user_service.UserService) *FarmingController {
	repo := repository.New(store)
	service := service.New(repo, userService)
	transport := transport.New(router, service)

	return &FarmingController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *FarmingController) Build() {
	c.transport.RegisterRoutes()
}

func (c *FarmingController) Run() {}

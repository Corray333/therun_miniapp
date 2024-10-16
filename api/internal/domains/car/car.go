package car

import (
	"github.com/Corray333/therun_miniapp/internal/domains/car/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/car/service"
	"github.com/Corray333/therun_miniapp/internal/domains/car/transport"
	round_service "github.com/Corray333/therun_miniapp/internal/domains/round/service"
	user_service "github.com/Corray333/therun_miniapp/internal/domains/user/service"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
)

type CarController struct {
	repo      repository.CarRepository
	service   service.CarService
	transport transport.CarTransport
}

func NewCarController(router *chi.Mux, store *storage.Storage, userService *user_service.UserService, roundService *round_service.RoundService) *CarController {
	repo := repository.New(store)
	service := service.New(repo, userService, roundService)
	transport := transport.New(router, service)

	return &CarController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *CarController) Build() {
	c.transport.RegisterRoutes()
}

func (c *CarController) Run() {}

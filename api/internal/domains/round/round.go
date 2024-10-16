package round

import (
	"github.com/Corray333/therun_miniapp/internal/domains/round/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/round/service"
	"github.com/Corray333/therun_miniapp/internal/domains/round/transport"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
)

type RoundController struct {
	repo      repository.RoundRepository
	service   service.RoundService
	transport transport.RoundTransport
}

func NewRoundController(router *chi.Mux, store *storage.Storage) *RoundController {
	repo := repository.New(store)
	service := service.New(repo)
	transport := transport.New(router, service)

	return &RoundController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *RoundController) Build() {
	c.transport.RegisterRoutes()
}

func (c *RoundController) Run() {
	c.service.RunRounds()
}

func (c *RoundController) GetService() *service.RoundService {
	return &c.service
}

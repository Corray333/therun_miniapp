package city

import (
	"github.com/Corray333/therun_miniapp/internal/domains/city/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/city/service"
	"github.com/Corray333/therun_miniapp/internal/domains/city/transport"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
)

type CityController struct {
	repo      repository.CityRepository
	service   service.CityService
	transport transport.CityTransport
}

func NewCityController(router *chi.Mux, store *storage.Storage) *CityController {
	repo := repository.New(store)
	service := service.New(repo)
	transport := transport.New(router, service)

	return &CityController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *CityController) Build() {
	c.transport.RegisterRoutes()
}

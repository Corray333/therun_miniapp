package city

import (
	"context"

	"github.com/Corray333/therun_miniapp/internal/domains/city/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/city/service"
	"github.com/Corray333/therun_miniapp/internal/domains/city/transport"
	user_types "github.com/Corray333/therun_miniapp/internal/domains/user/types"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
)

type CityController struct {
	repo      repository.CityRepository
	service   service.CityService
	transport transport.CityTransport
}
type userRepository interface {
	ChangeBalances(ctx context.Context, userID int64, changes []user_types.BalanceChange) error
}

func NewCityController(router *chi.Mux, store *storage.Storage, userRepository userRepository) *CityController {
	repo := repository.New(store, userRepository)
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

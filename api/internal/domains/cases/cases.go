package cases

import (
	"github.com/Corray333/therun_miniapp/internal/domains/cases/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/cases/service"
	"github.com/Corray333/therun_miniapp/internal/domains/cases/transport"
	user_service "github.com/Corray333/therun_miniapp/internal/domains/user/service"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
)

type CaseController struct {
	repo      repository.CityRepository
	service   service.CityService
	transport transport.CityTransport
}

func NewCasesController(router *chi.Mux, store *storage.Storage, userService *user_service.UserService) *CaseController {
	repo := repository.New(store)
	service := service.New(repo, userService)
	transport := transport.New(router, service)

	return &CaseController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *CaseController) Build() {
	c.transport.RegisterRoutes()
}

func (c *CaseController) Run() {}

package user

import (
	"github.com/Corray333/therun_miniapp/internal/domains/user/external"
	"github.com/Corray333/therun_miniapp/internal/domains/user/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/user/service"
	"github.com/Corray333/therun_miniapp/internal/domains/user/transport"
	"github.com/Corray333/therun_miniapp/internal/files"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/Corray333/therun_miniapp/internal/telegram"
	"github.com/go-chi/chi"
)

type UserController struct {
	repo      *repository.UserRepository
	service   *service.UserService
	transport *transport.UserTransport
}

func (c *UserController) GetService() *service.UserService {
	return c.service
}

func NewUserController(router *chi.Mux, tg *telegram.TelegramClient, store *storage.Storage, fileManager *files.FileManager) *UserController {
	repo := repository.New(store)
	external := external.New(tg)
	service := service.New(repo, external, fileManager)
	transport := transport.New(router, service)

	return &UserController{
		repo:      repo,
		service:   service,
		transport: transport,
	}
}

func (c *UserController) Build() {
	c.transport.RegisterRoutes()
}

func (c *UserController) GetRepository() *repository.UserRepository {
	return c.repo
}

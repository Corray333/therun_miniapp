package user

import (
	"github.com/Corray333/therun_miniapp/internal/domains/user/external"
	"github.com/Corray333/therun_miniapp/internal/domains/user/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/user/service"
	"github.com/Corray333/therun_miniapp/internal/domains/user/transport"
	"github.com/Corray333/therun_miniapp/internal/files"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserController struct {
	repo      repository.UserRepository
	service   service.UserService
	transport transport.UserTransport
}

func NewUserController(router *chi.Mux, tg *tgbotapi.BotAPI, store *storage.Storage, fileManager *files.FileManager) *UserController {
	repo := repository.New(store)
	external := external.New(tg)
	service := service.New(repo, external, fileManager)
	transport := transport.New(router, tg, service)

	return &UserController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

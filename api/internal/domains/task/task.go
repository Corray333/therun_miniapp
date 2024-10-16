package task

import (
	"github.com/Corray333/therun_miniapp/internal/domains/task/external"
	"github.com/Corray333/therun_miniapp/internal/domains/task/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/task/service"
	"github.com/Corray333/therun_miniapp/internal/domains/task/transport"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/Corray333/therun_miniapp/internal/telegram"
	"github.com/go-chi/chi"
)

type TaskController struct {
	repo      repository.TaskRepository
	service   service.TaskService
	transport transport.TaskTransport
}

func NewTaskController(router *chi.Mux, store *storage.Storage, tg *telegram.TelegramClient) *TaskController {
	repo := repository.New(store)
	ext := external.New(tg)
	service := service.New(repo, ext)
	transport := transport.New(router, service)

	return &TaskController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *TaskController) Build() {
	c.transport.RegisterRoutes()
}

func (c *TaskController) Run() {}

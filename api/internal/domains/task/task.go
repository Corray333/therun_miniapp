package task

import (
	"github.com/Corray333/therun_miniapp/internal/domains/task/repository"
	"github.com/Corray333/therun_miniapp/internal/domains/task/service"
	"github.com/Corray333/therun_miniapp/internal/domains/task/transport"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/go-chi/chi"
)

type TaskController struct {
	repo      repository.TaskRepository
	service   service.TaskService
	transport transport.TaskTransport
}

func NewTaskController(router *chi.Mux, store *storage.Storage) *TaskController {
	repo := repository.New(store)
	service := service.New(repo)
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

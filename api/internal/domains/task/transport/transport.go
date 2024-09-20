package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Corray333/therun_miniapp/internal/domains/task/types"
	global_types "github.com/Corray333/therun_miniapp/internal/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type TaskTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	GetTasks(userID int64, lang string) ([]*types.Task, error)
}

func New(router *chi.Mux, service service) *TaskTransport {
	return &TaskTransport{
		router:  router,
		service: service,
	}
}

func (t *TaskTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())
		r.Get("/api/tasks", t.getTasks)
	})
}

func (t *TaskTransport) getTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		slog.Error("invalid user id")
		return
	}

	lang := r.URL.Query().Get("lang")

	tasks, err := t.service.GetTasks(userID, lang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to get tasks: " + err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to encode tasks: " + err.Error())
	}
}

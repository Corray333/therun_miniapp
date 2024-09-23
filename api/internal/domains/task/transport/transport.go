package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

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
	CheckTask(userID, taskID int64) (done bool, err error)
	Claim(userID, taskID int64) (done bool, pointsBalance, raceBalance, keysBalance int, err error)
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
		r.Post("/api/tasks/{taskID}/check", t.checkTask)
		r.Post("/api/tasks/{taskID}/claim", t.claimTask)
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

type CheckTaskResponse struct {
	Done bool `json:"done"`
}

func (t *TaskTransport) checkTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		slog.Error("invalid user id")
		return
	}

	taskIDStr := chi.URLParam(r, "taskID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		slog.Error("invalid task id")
		return
	}

	done, err := t.service.CheckTask(userID, int64(taskID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to check task: " + err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(CheckTaskResponse{
		Done: done,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to encode done: " + err.Error())
	}
}

type ClaimTaskResponse struct {
	Done         bool `json:"done"`
	PointBalance int  `json:"pointBalance"`
	RaceBalance  int  `json:"raceBalance"`
	KeyBalance   int  `json:"keyBalance"`
}

func (t *TaskTransport) claimTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		slog.Error("invalid user id")
		return
	}

	taskIDStr := chi.URLParam(r, "taskID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		slog.Error("invalid task id")
		return
	}

	done, pointsBalance, raceBalance, keysBalance, err := t.service.Claim(userID, int64(taskID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to claim task: " + err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(ClaimTaskResponse{
		Done:         done,
		PointBalance: pointsBalance,
		RaceBalance:  raceBalance,
		KeyBalance:   keysBalance,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to encode done: " + err.Error())
	}
}

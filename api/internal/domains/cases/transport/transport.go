package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Corray333/therun_miniapp/internal/domains/cases/types"
	global_types "github.com/Corray333/therun_miniapp/internal/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type CityTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	OpenCase(userID int64, caseType string) (*types.Reward, error)
	GetCases(userID int64) ([]types.Case, error)
}

func New(router *chi.Mux, service service) *CityTransport {
	return &CityTransport{
		router:  router,
		service: service,
	}
}

func (t *CityTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())

		r.Post("/api/cases/open", t.openCase)
		r.Get("/api/cases", t.getCases)
	})
}

type OpenCaseRequest struct {
	CaseType string `json:"caseType"`
}

func (t *CityTransport) openCase(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)

	if !ok {
		http.Error(w, "failed to get user id from context", http.StatusInternalServerError)
		slog.Error("failed to get user id from context")
		return
	}

	req := &OpenCaseRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		slog.Error("failed to decode request: " + err.Error())
		return
	}

	reward, err := t.service.OpenCase(userID, req.CaseType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to open case: " + err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(reward); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		slog.Error("failed to encode response: " + err.Error())
		return
	}

}

func (t *CityTransport) getCases(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)

	if !ok {
		http.Error(w, "failed to get user id from context", http.StatusInternalServerError)
		slog.Error("failed to get user id from context")
		return
	}

	cases, err := t.service.GetCases(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to get cases" + err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(cases); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		slog.Error("failed to encode response" + err.Error())
		return
	}
}

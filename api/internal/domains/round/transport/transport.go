package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Corray333/therun_miniapp/internal/domains/battle/types"
	global_types "github.com/Corray333/therun_miniapp/internal/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type RoundTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	GetRound(userID int64) *types.Round
}

func New(router *chi.Mux, service service) *RoundTransport {
	return &RoundTransport{
		router:  router,
		service: service,
	}
}

func (t *RoundTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())

		r.Get("/api/round", t.getRound)
	})
}

func (t *RoundTransport) getRound(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	round := t.service.GetRound(userID)

	if err := json.NewEncoder(w).Encode(round); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

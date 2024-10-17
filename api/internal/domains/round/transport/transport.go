package transport

import (
	"encoding/json"
	"net/http"

	"github.com/Corray333/therun_miniapp/internal/domains/round/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type RoundTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	GetRound() *types.Round
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
	round := t.service.GetRound()

	if err := json.NewEncoder(w).Encode(round); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

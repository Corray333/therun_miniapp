package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	battle_service "github.com/Corray333/therun_miniapp/internal/domains/battle/service"
	"github.com/Corray333/therun_miniapp/internal/domains/battle/types"
	global_types "github.com/Corray333/therun_miniapp/internal/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type BattleTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	GetRound(userID int64) *types.Round
	MakeBet(userID int64, battleID, pick int) error
}

func New(router *chi.Mux, service service) *BattleTransport {
	return &BattleTransport{
		router:  router,
		service: service,
	}
}

func (t *BattleTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())

		r.Get("/api/round", t.getRound)
		r.Post("/api/battles/{battle_id}/bet", t.makeBet)
	})
}

func (t *BattleTransport) getRound(w http.ResponseWriter, r *http.Request) {
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

func (t *BattleTransport) makeBet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	battleIDStr := chi.URLParam(r, "battle_id")
	if battleIDStr == "" {
		http.Error(w, "battle id not found in url", http.StatusBadRequest)
	}

	battleID, err := strconv.Atoi(battleIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var bet types.Bet
	if err := json.NewDecoder(r.Body).Decode(&bet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("error decoding bet: " + err.Error())
		return
	}

	if err := t.service.MakeBet(userID, battleID, bet.Pick); err != nil {
		if err == battle_service.ErrTooLate {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("error making bet: " + err.Error())
	}

	w.WriteHeader(http.StatusOK)
}

package transport

import (
	"encoding/json"
	"net/http"

	"github.com/Corray333/therun_miniapp/internal/domains/battle/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type BattleTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	GetRound() *types.Round

	// GetBattles()
	// Bet()
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
	})
}

func (t *BattleTransport) getRound(w http.ResponseWriter, r *http.Request) {
	round := t.service.GetRound()

	round.Battles = append(round.Battles, types.Battle{
		ID:      8982,
		RoundID: round.ID,
		User: types.User{
			ID:       7850,
			Username: "Monokoleso",
			Photo:    "https://therunapp.com/public/images/wPHSWUibXv.png",
			Miles:    777209.8347201593,
		},
		Opponent: types.User{
			ID:       2339,
			Username: "aleksandr28",
			Photo:    "https://therunapp.com/public/images/HrNgYAMmiU.jpg",
			Miles:    2301443.028976749,
		},
	})

	if err := json.NewEncoder(w).Encode(round); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"

	farm_service "github.com/Corray333/therun_miniapp/internal/domains/farming/service"
	"github.com/Corray333/therun_miniapp/internal/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type FarmingTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	ClaimTokens(userID int64) (pointsGot, pointsBalnce, farmingTime, maxPoints int, farmingFrom int64, err error)
	StartFarming(userID int64) (farmingFrom int64, err error)
}

func New(router *chi.Mux, service service) *FarmingTransport {
	return &FarmingTransport{
		router:  router,
		service: service,
	}
}

func (t *FarmingTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())
		r.Post("/api/farming/claim", t.claimTokens)
		r.Post("/api/farming/start", t.startFarming)
	})
}

type claimTokensResponse struct {
	PointsGot    int   `json:"pointsGot"`
	PointBalance int   `json:"pointBalance"`
	FarmingTime  int   `json:"farmingTime"`
	MaxPoints    int   `json:"maxPoints"`
	LastClaim    int64 `json:"lastClaim"`
}

func (t *FarmingTransport) claimTokens(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(types.ContextID).(int64)
	pointsGot, pointsBalance, farmingTime, maxPoints, farmingFrom, err := t.service.ClaimTokens(userID)

	if err != nil {
		if err == farm_service.ErrClaimTooEarly {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Error("failed to claim tokens" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(claimTokensResponse{
		PointsGot:    pointsGot,
		PointBalance: pointsBalance,
		FarmingTime:  farmingTime,
		MaxPoints:    maxPoints,
		LastClaim:    farmingFrom,
	}); err != nil {
		slog.Error("failed to encode response" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type startFarmingResponse struct {
	FarmingFrom int64 `json:"farmingFrom"`
}

func (t *FarmingTransport) startFarming(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(types.ContextID).(int64)

	if !ok {
		http.Error(w, "failed to get user id from context", http.StatusInternalServerError)
		slog.Error("failed to get user id from context")
		return
	}

	farmingFrom, err := t.service.StartFarming(userID)
	if err != nil {
		slog.Error("failed to start farming" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(startFarmingResponse{
		FarmingFrom: farmingFrom,
	}); err != nil {
		slog.Error("failed to encode response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

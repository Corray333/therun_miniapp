package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Corray333/therun_miniapp/internal/domains/city/types"
	global_types "github.com/Corray333/therun_miniapp/internal/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type CityTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	GetCity(ctx context.Context, userID int64) (map[types.BuildingType]types.BuildingPublic, error)
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

		r.Get("/api/city", t.getCity)
	})
}

func (t *CityTransport) getCity(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		slog.Error("invalid user id")
		return
	}

	buildings, err := t.service.GetCity(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to get city: " + err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(buildings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to encode buildings: " + err.Error())
		return
	}
}

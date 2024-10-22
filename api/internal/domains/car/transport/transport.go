package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Corray333/therun_miniapp/internal/domains/car/types"
	round_types "github.com/Corray333/therun_miniapp/internal/domains/round/types"
	global_types "github.com/Corray333/therun_miniapp/internal/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type CarTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	GetAllCars(ctx context.Context) []types.Car
	GetOwnedCars(ctx context.Context, userID int64) ([]types.Car, error)
	GetMainCar(ctx context.Context, userID int64) (*types.Car, error)
	BuyCar(ctx context.Context, element round_types.Element, userID int64) error
	GetCarByID(ctx context.Context, carID int64) (*types.Car, error)
	PickCar(ctx context.Context, carID int64, userID int64) error

	GetRace(ctx context.Context, userID int64) (race *types.RaceState, err error)
	StartRace(ctx context.Context, userID int64) (race *types.RaceState, err error)
	EndRace(ctx context.Context, userID int64) (race *types.RaceState, err error)

	GetModulesOfUser(ctx context.Context, carID int64) ([]types.Module, error)

	BuyFuel(ctx context.Context, userID int64) error
	BuyHealth(ctx context.Context, userID int64) error
}

func New(router *chi.Mux, service service) *CarTransport {
	return &CarTransport{
		router:  router,
		service: service,
	}
}

func (t *CarTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())

		r.Get("/api/cars/all", t.getAllCars)      // All cars
		r.Get("/api/cars/main", t.getMainCar)     // Current car
		r.Get("/api/cars/owned", t.getOwnedCars)  // Get all available cars
		r.Get("/api/cars/{car_id}", t.getCarByID) // Get car by id
		r.Post("/api/buy-car", t.buyCar)          // Choose start car
		r.Post("/api/pick-car", t.pickCar)        // Choose start car
		r.Get("/api/race", t.getRace)             // Get race state
		r.Post("/api/start-race", t.startRace)    // Start moving
		r.Post("/api/end-race", t.endRace)        // End moving

		r.Get("/api/cars/modules", t.getCarModules)

		r.Post("/api/buy-fuel", t.buyFuel)
		r.Post("/api/buy-health", t.buyHealth)
	})
}

func (t *CarTransport) getAllCars(w http.ResponseWriter, r *http.Request) {
	cars := t.service.GetAllCars(r.Context())

	if err := json.NewEncoder(w).Encode(cars); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *CarTransport) buyCar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	element := r.URL.Query().Get("element")
	if element == "" {
		http.Error(w, "element not found in query", http.StatusBadRequest)
		slog.Error("element not found in query")
		return
	}

	if err := t.service.BuyCar(r.Context(), round_types.Element(element), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *CarTransport) pickCar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	carIDStr := r.URL.Query().Get("car_id")
	if carIDStr == "" {
		http.Error(w, "car id not found in query", http.StatusBadRequest)
		slog.Error("car id not found in query")
		return
	}

	carID, err := strconv.ParseInt(carIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := t.service.PickCar(r.Context(), carID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *CarTransport) getMainCar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	car, err := t.service.GetMainCar(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(car); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *CarTransport) getCarByID(w http.ResponseWriter, r *http.Request) {

	carIDStr := chi.URLParam(r, "car_id")
	if carIDStr == "" {
		http.Error(w, "car id not found in query", http.StatusBadRequest)
		slog.Error("car id not found in query")
		return
	}

	carID, err := strconv.ParseInt(carIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	car, err := t.service.GetCarByID(r.Context(), carID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(car); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *CarTransport) getOwnedCars(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	cars, err := t.service.GetOwnedCars(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cars); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *CarTransport) getRace(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	race, err := t.service.GetRace(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(race); err != nil {
		slog.Error("failed to encode an info about race: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *CarTransport) startRace(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	race, err := t.service.StartRace(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(race); err != nil {
		slog.Error("failed to encode an info about race: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *CarTransport) endRace(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	race, err := t.service.EndRace(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(race); err != nil {
		slog.Error("failed to encode an info about race: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *CarTransport) getCarModules(w http.ResponseWriter, r *http.Request) {
	carIDStr := r.URL.Query().Get("car_id")
	if carIDStr == "" {
		http.Error(w, "car id not found in query", http.StatusBadRequest)
		slog.Error("car id not found in query")
		return
	}

	carID, err := strconv.ParseInt(carIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	modules, err := t.service.GetModulesOfUser(r.Context(), carID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(modules); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *CarTransport) buyFuel(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	if err := t.service.BuyFuel(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *CarTransport) buyHealth(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	if err := t.service.BuyHealth(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

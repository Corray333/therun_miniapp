package transport

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Corray333/therun_miniapp/internal/domains/user/types"
	global_types "github.com/Corray333/therun_miniapp/internal/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type UserTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	GetUser(userID int64) (*types.User, error)
	AuthUser(initData, refCode string) (accessToken string, isNew bool, err error)
	ListReferals(userID int64) ([]types.Referal, error)
	CountReferals(userID int64) (count, level, nextLevelCount, previousLevelCount int, err error)
}

func New(router *chi.Mux, service service) *UserTransport {
	return &UserTransport{
		router:  router,
		service: service,
	}
}

func (t *UserTransport) RegisterRoutes() {
	t.router.Post("/api/users/auth", t.auth)
	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())
		r.Get("/api/users/{userID}", t.getUser)
		r.Get("/api/users/{userID}/referals", t.listReferals)
		r.Get("/api/users/{userID}/referals/info", t.countReferals)
	})
}

type authRequest struct {
	InitData string `json:"initData"`
	RefCode  string `json:"refCode"`
}

type authResponse struct {
	AccessToken string `json:"accessToken"`
	IsNew       bool   `json:"isNew"`
}

func (t *UserTransport) auth(w http.ResponseWriter, r *http.Request) {
	req := &authRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("\nRequest: %+v\n", req)
	token, isNew, err := t.service.AuthUser(req.InitData, req.RefCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		slog.Error(err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(&authResponse{AccessToken: token, IsNew: isNew}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *UserTransport) getUser(w http.ResponseWriter, r *http.Request) {
	userIDRaw := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDRaw, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := t.service.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *UserTransport) listReferals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	referals, err := t.service.ListReferals(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to list referals" + err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(referals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to encode referals" + err.Error())
		return
	}
}

type countReferalsResponse struct {
	Count              int `json:"count"`
	Level              int `json:"level"`
	NextLevelCount     int `json:"nextLevelCount"`
	PreviousLevelCount int `json:"previousLevelCount"`
}

func (t *UserTransport) countReferals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	count, level, nextLevelCount, previousLevelCount, err := t.service.CountReferals(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to count referals" + err.Error())
		return
	}

	resp := &countReferalsResponse{
		Count:              count,
		Level:              level,
		NextLevelCount:     nextLevelCount,
		PreviousLevelCount: previousLevelCount,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to encode referals" + err.Error())
		return
	}
}

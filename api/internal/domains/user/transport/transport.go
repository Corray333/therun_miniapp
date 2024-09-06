package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Corray333/therun_miniapp/internal/domains/user/types"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type UserTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	GetUser(userID int64) (*types.User, error)
	AuthUser(initData, refCode string) (accessToken string, err error)
	ListReferals(userID int64) ([]types.Referal, error)
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
	})
}

type authRequest struct {
	InitData string `json:"initData"`
	RefCode  string `json:"refCode"`
}

type authResponse struct {
	AccessToken string `json:"accessToken"`
}

func (t *UserTransport) auth(w http.ResponseWriter, r *http.Request) {
	req := &authRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := t.service.AuthUser(req.InitData, req.RefCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewEncoder(w).Encode(&authResponse{AccessToken: token}); err != nil {
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
	userIDRaw := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDRaw, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	referals, err := t.service.ListReferals(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(referals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

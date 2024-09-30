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
	AuthUser(initData, refCode string) (accessToken string, isNew bool, isPremium bool, err error)
	ListActivatedReferals(userID int64) ([]types.Referal, error)
	ListNotActivatedReferals(userID int64) ([]types.Referal, error)
	CountReferals(userID int64) (refsActivated, refsFrozen, rewardsFrozen, rewardsAvailible int, err error)
	DailyCheck(userID int64) (dailyCheckStreak int, dailyCheckLast int64, err error)
	ClaimRefs(userID int64) (rewardsGot int, err error)

	SetPremium(userID int64, isPremium bool) error
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
		// r.Post("/api/users/{userID}/daily_check", t.dailyCheck)
		r.Post("/api/users/{userID}/referals/claim", t.claimRefs)
	})
}

type authRequest struct {
	InitData  string `json:"initData"`
	RefCode   string `json:"refCode"`
	IsPremium bool   `json:"isPremium"`
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

	// TODO: Put isPremium in token payload
	token, isNew, isPremium, err := t.service.AuthUser(req.InitData, req.RefCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		slog.Error(err.Error())
		return
	}

	creds, err := auth.ExtractCredentials(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to extract credentials: " + err.Error())
		return
	}

	if err := t.service.SetPremium(creds.ID, isPremium); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to set premium: " + err.Error())
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

	isActivated := r.URL.Query().Get("activated")

	var referals []types.Referal
	var err error

	if isActivated == "true" {
		referals, err = t.service.ListActivatedReferals(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.Error("failed to list referals" + err.Error())
			return
		}
	} else {
		referals, err = t.service.ListNotActivatedReferals(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.Error("failed to list referals" + err.Error())
			return
		}
	}

	if err := json.NewEncoder(w).Encode(referals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to encode referals" + err.Error())
		return
	}
}

type countReferalsResponse struct {
	RefsActivated    int `json:"refsActivated"`
	RefsFrozen       int `json:"refsFrozen"`
	RewardsFrozen    int `json:"rewardsFrozen"`
	RewardsAvailible int `json:"rewardsAvailible"`
}

func (t *UserTransport) countReferals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	refsActivated, refsFrozen, rewardsFrozen, rewardsAvailible, err := t.service.CountReferals(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to count referals" + err.Error())
		return
	}

	resp := &countReferalsResponse{
		RefsActivated:    refsActivated,
		RefsFrozen:       refsFrozen,
		RewardsFrozen:    rewardsFrozen,
		RewardsAvailible: rewardsAvailible,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to encode referals" + err.Error())
		return
	}
}

type dailyCheckResponse struct {
	DailyCheckStreak int   `json:"dailyCheckStreak"`
	DailyCheckLast   int64 `json:"dailyCheckLast"`
}

func (t *UserTransport) dailyCheck(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	dailyCheckStreak, dailyCheckLast, err := t.service.DailyCheck(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to daily check: " + err.Error())
		return
	}

	resp := &dailyCheckResponse{
		DailyCheckStreak: dailyCheckStreak,
		DailyCheckLast:   dailyCheckLast,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to encode daily check: " + err.Error())
		return
	}
}

type claimRefsResponse struct {
	RewardsGot int `json:"rewardsGot"`
}

func (t *UserTransport) claimRefs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(global_types.ContextID).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	rewardsGot, err := t.service.ClaimRefs(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to claim referals" + err.Error())
		return
	}

	resp := &claimRefsResponse{
		RewardsGot: rewardsGot,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to encode referals" + err.Error())
		return
	}
}

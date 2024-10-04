package transport

import (
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type CityTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
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

	})
}

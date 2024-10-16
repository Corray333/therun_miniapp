package transport

import (
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
	"github.com/go-chi/chi"
)

type CarTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
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

		// r.Get("/race/all-cars", t.getAllCars) // All cars
		// r.Get("/race/car", t.getCar) // Current car
		// r.Get("/race/cars", t.getCars) // Get all available cars
		// r.Post("/race/buy-car", t.buyCar) // Buy car
		// r.Get("/race/state", t.getRace) // Get race state
	})
}

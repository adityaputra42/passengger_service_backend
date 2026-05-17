package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func AirportRoutes(r chi.Router, h *handler.AirportHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
	adminMiddleware := middleware.RequireAdminArea(deps.RBACService)

	r.Route("/airport", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Get("/code/{code}", h.GetByCode)

		r.Group(func(r chi.Router) {
			r.Use(adminMiddleware)

			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})
}

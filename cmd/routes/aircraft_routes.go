package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func AirchaftRoutes(r chi.Router, h *handler.AircraftHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/aircraft", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Get("/{id}/seats", h.GetWithSeats)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAdminArea(deps.RBACService))

			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
			r.Post("/{id}/seats/generate", h.GenerateSeats)

		})
	})
}

package routes

import (
	"passenger_service_backend/internal/handler"

	"github.com/go-chi/chi/v5"
)

func FlightRoutes(r chi.Router, h *handler.FlightHandler, deps Dependencies) {
	// authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/flights", func(r chi.Router) {
		// r.Use(authMiddleware)
		r.Get("/search", h.Search)
		r.Get("/{id}", h.Get)
		r.Get("/{id}/seat-map", h.SeatMap)

		r.Group(func(r chi.Router) {
			// r.Use(middleware.RequireAdminArea(deps.RBACService))
			r.Post("/generate", h.Generate)
			r.Patch("/{id}/status", h.UpdateStatus)

		})
	})
}

package routes

import (
	"passenger_service_backend/internal/handler"

	"github.com/go-chi/chi/v5"
)

func BaggageRoutes(r chi.Router, h *handler.BaggageHandler, deps Dependencies) {
	// authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/baggage", func(r chi.Router) {
		// r.Use(authMiddleware)
		r.Post("/", h.Add)
		r.Get("/passenger/{passengerID}", h.GetByPassenger)

		r.Group(func(r chi.Router) {
			// r.Use(middleware.RequireAdminArea(deps.RBACService))
			r.Put("/{id}/status", h.UpdateStatus)

		})
	})
}

package routes

// ============================================================
// baggage_routes.go
// ============================================================
//
// Baggage routes:
//   POST /baggage                           → authenticated (agent+ to add baggage)
//   GET  /baggage/passenger/{passengerID}   → authenticated
//   PUT  /baggage/{id}/status               → agent+ (airport staff updating status)

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func BaggageRoutes(r chi.Router, h *handler.BaggageHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
	agentMiddleware := middleware.RequireAgentOrAbove(deps.RBACService)

	r.Route("/baggage", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/", h.Add)
		r.Get("/passenger/{passengerID}", h.GetByPassenger)

		r.Group(func(r chi.Router) {
			r.Use(agentMiddleware)
			r.Put("/{id}/status", h.UpdateStatus)
		})
	})
}

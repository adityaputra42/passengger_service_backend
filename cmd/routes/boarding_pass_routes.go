package routes

import (
	"passenger_service_backend/internal/handler"

	"github.com/go-chi/chi/v5"
)

func BoardingPassRoutes(r chi.Router, h *handler.BoardingPassHandler, deps Dependencies) {
	// authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/boarding_passes", func(r chi.Router) {
		// r.Use(authMiddleware)
		r.Post("/", h.Issue)
		r.Get("/passenger/{passengerID}/segment/{segmentID}", h.Get)
		r.Get("/segment/{segmentID}", h.GetBySegment)
	})
}

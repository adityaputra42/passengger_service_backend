package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// BoardingPass routes:
//   POST /boarding_passes                              → agent+ (issue a boarding pass)
//   GET  /boarding_passes/passenger/{id}/segment/{id} → authenticated (own boarding pass)
//   GET  /boarding_passes/segment/{segmentID}          → agent+ (view all for a segment)

func BoardingPassRoutes(r chi.Router, h *handler.BoardingPassHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
	agentMiddleware := middleware.RequireAgentOrAbove(deps.RBACService)

	r.Route("/boarding_passes", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/passenger/{passengerID}/segment/{segmentID}", h.Get)

		r.Group(func(r chi.Router) {
			r.Use(agentMiddleware)
			r.Post("/", h.Issue)
			r.Get("/segment/{segmentID}", h.GetBySegment)
		})
	})
}

package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// Aircraft routes:
//   GET  /aircraft          → public (anyone can view aircraft list)
//   GET  /aircraft/{id}     → public
//   GET  /aircraft/{id}/seats → public
//   POST /aircraft          → admin+ only
//   PUT  /aircraft/{id}     → admin+ only
//   DELETE /aircraft/{id}   → admin+ only
//   POST /aircraft/{id}/seats/generate → admin+ only

func AirchaftRoutes(r chi.Router, h *handler.AircraftHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
	adminMiddleware := middleware.RequireAdminArea(deps.RBACService)

	r.Route("/aircraft", func(r chi.Router) {

		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Get("/{id}/seats", h.GetWithSeats)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Use(adminMiddleware)

			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
			r.Post("/{id}/seats/generate", h.GenerateSeats)
		})
	})
}

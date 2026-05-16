package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// ============================================================
// flight_routes.go
// ============================================================
// Flight routes:
//   GET  /flights/search        → public (anyone can search flights)
//   GET  /flights/{id}          → public
//   GET  /flights/{id}/seat-map → authenticated (need to log in to see prices/availability)
//   POST /flights/generate      → admin+ only
//   PATCH /flights/{id}/status  → admin+ only

func FlightRoutes(r chi.Router, h *handler.FlightHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
	adminMiddleware := middleware.RequireAdminArea(deps.RBACService)

	r.Route("/flights", func(r chi.Router) {
		r.Get("/search", h.Search)
		r.Get("/{id}", h.Get)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/{id}/seat-map", h.SeatMap)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Use(adminMiddleware)
			r.Get("/", h.GetAll)
			r.Post("/generate", h.Generate)
			r.Patch("/{id}/status", h.UpdateStatus)
		})
	})
}

package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// ============================================================
// flight_schedule_routes.go
// ============================================================
// FlightSchedule routes:
//   GET  /schedules        → public (flight schedules are public info)
//   GET  /schedules/{id}   → public
//   POST /schedules        → admin+ only
//   PUT  /schedules/{id}   → admin+ only
//   DELETE /schedules/{id} → admin+ only

func FlightScheduleRoutes(r chi.Router, h *handler.FlightScheduleHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
	adminMiddleware := middleware.RequireAdminArea(deps.RBACService)

	r.Route("/schedules", func(r chi.Router) {
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Use(adminMiddleware)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})
}

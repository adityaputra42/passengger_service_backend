package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func DashboardRoutes(
	r chi.Router,
	h *handler.DashboardHandler,
	deps Dependencies) {

	r.Route("/dashboard", func(r chi.Router) {
		authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
		adminMiddleware := middleware.RequireAdminArea(deps.RBACService)

		r.Use(authMiddleware, adminMiddleware)

		r.Get("/summary", h.Summary)

		r.Get("/revenue-trend", h.RevenueTrend)

		r.Get("/booking-status", h.BookingStatus)

		r.Get("/today-flights", h.TodayFlights)

		r.Get("/recent-bookings", h.RecentBookings)
	})
}

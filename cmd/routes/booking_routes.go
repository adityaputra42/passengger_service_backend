package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// Booking (PNR) routes:
//   POST /bookings                          → authenticated (any logged-in user can book)
//   GET  /bookings                          → admin+ only (view ALL bookings)
//   GET  /bookings/{id}                     → authenticated (own booking or admin)
//   GET  /bookings/locator/{locator}        → authenticated
//   PUT  /bookings/{id}/contact             → authenticated
//   DELETE /bookings/{id}                   → authenticated (cancel own booking)
//   POST /bookings/passengers/{id}/ssr      → authenticated
//   DELETE /bookings/ssr/{ssrID}            → authenticated
//   POST /bookings/passengers/{id}/meal     → authenticated
//   DELETE /bookings/meal/{mealID}          → authenticated

func BookingRoutes(r chi.Router, h *handler.BookingHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
	adminMiddleware := middleware.RequireAdminArea(deps.RBACService)

	r.Route("/bookings", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByID)
		r.Get("/locator/{locator}", h.GetByLocator)
		r.Put("/{id}/contact", h.UpdateContact)
		r.Delete("/{id}", h.Cancel)
		r.Post("/passengers/{passengerID}/ssr", h.AddSSR)
		r.Delete("/ssr/{ssrID}", h.RemoveSSR)
		r.Post("/passengers/{passengerID}/meal", h.AddMeal)
		r.Delete("/meal/{mealID}", h.RemoveMeal)

		// Admin-only: view all bookings
		r.Group(func(r chi.Router) {
			r.Use(adminMiddleware)
			r.Get("/", h.List)
		})
	})
}

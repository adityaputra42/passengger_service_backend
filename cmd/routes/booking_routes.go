package routes

import (
	"passenger_service_backend/internal/handler"

	"github.com/go-chi/chi/v5"
)

func BookingRoutes(r chi.Router, h *handler.BookingHandler, deps Dependencies) {
	// authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/bookings", func(r chi.Router) {
		// r.Use(authMiddleware)
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByID)
		r.Get("/locator/{locator}", h.GetByLocator)
		r.Put("/{id}/contact", h.UpdateContact)
		r.Delete("/{id}", h.Cancel)
		r.Post("/passengers/{passengerID}/ssr", h.AddSSR)
		r.Delete("/ssr/{ssrID}", h.RemoveSSR)
		r.Post("/passengers/{passengerID}/meal", h.AddMeal)
		r.Delete("/meal/{mealID}", h.RemoveMeal)

		r.Group(func(r chi.Router) {
			// r.Use(middleware.RequireAdminArea(deps.RBACService))
			r.Get("/", h.List)

		})
	})
}

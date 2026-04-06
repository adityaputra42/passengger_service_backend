package routes

import (
	"passenger_service_backend/internal/handler"

	"github.com/go-chi/chi/v5"
)

func PaymentRoutes(r chi.Router, h *handler.PaymentHandler, deps Dependencies) {
	// authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/payments", func(r chi.Router) {
		// r.Use(authMiddleware)
		r.Post("/", h.Initiate)
		r.Get("/", h.ListByPNR)

		r.Group(func(r chi.Router) {
			// r.Use(middleware.RequireAdminArea(deps.RBACService))
			r.Post("/{id}/confirm", h.Confirm)
			r.Post("/{id}/refund", h.Refund)

		})
	})
}

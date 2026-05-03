package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// ============================================================
// payment_routes.go
// ============================================================
// Payment routes:
//   POST /payments              → authenticated (initiate payment for own booking)
//   GET  /payments?pnr_id=      → authenticated (view own payments)
//   POST /payments/{id}/confirm → admin+ (payment gateway callback / internal webhook)
//   POST /payments/{id}/refund  → admin+ (process refund)

func PaymentRoutes(r chi.Router, h *handler.PaymentHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
	adminMiddleware := middleware.RequireAdminArea(deps.RBACService)

	r.Route("/payments", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/", h.Initiate)
		r.Get("/", h.ListByPNR)

		r.Group(func(r chi.Router) {
			r.Use(adminMiddleware)
			r.Post("/{id}/confirm", h.Confirm)
			r.Post("/{id}/refund", h.Refund)
		})
	})
}

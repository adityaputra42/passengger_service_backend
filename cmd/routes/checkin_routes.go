package routes

import (
	"passenger_service_backend/internal/handler"

	"github.com/go-chi/chi/v5"
)

func CheckInRoutes(r chi.Router, h *handler.CheckinHandler, deps Dependencies) {
	// authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/checkin", func(r chi.Router) {
		// r.Use(authMiddleware)
		r.Post("/", h.Checkin)
		r.Get("/passenger/{passengerID} ", h.GetByPassenger)
	})
}

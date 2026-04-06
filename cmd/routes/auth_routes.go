package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router, h *handler.AuthHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.SignIn)
		r.Post("/admin/login", h.AdminLogin)
		r.Post("/refresh", h.Refresh)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/me", h.Me)
			r.Put("/change-password", h.ChangePassword)
			r.Post("/logout", h.Logout)

		})

	})
}

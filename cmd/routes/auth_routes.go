package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router, h *handler.AuthHandler, deps Dependencies) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.SignIn)
		r.Post("/admin/login", h.AdminLogin)
		r.Post("/logout", h.Logout)
		r.With(middleware.AuthMiddleware(deps.UserService, deps.JWTService)).Get("/me", h.Me)
		r.With(middleware.AuthMiddleware(deps.UserService, deps.JWTService)).Post("/refresh", h.Refresh)
		r.With(middleware.AuthMiddleware(deps.UserService, deps.JWTService)).Put("/change-password", h.ChangePassword)

	})
}

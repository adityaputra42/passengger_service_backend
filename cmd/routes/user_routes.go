package routes

import (
	"passenger_service_backend/internal/handler"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router, h *handler.UserHandler, deps Dependencies) {
	// authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/users", func(r chi.Router) {
		// r.Use(authMiddleware)
		// r.Use(middleware.RequireAdminArea(deps.RBACService))

		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{uid}", h.Get)
		r.Put("/{uid}", h.Update)
		r.Delete("/{uid}", h.Delete)
		r.Put("/me/profile", h.UpdateProfile)

	})
}

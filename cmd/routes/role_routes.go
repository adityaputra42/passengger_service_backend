package routes

import (
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RoleRoutes(r chi.Router, h *handler.RoleHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)
	adminMiddleware := middleware.RequireAdminArea(deps.RBACService)
	superAdminMiddleware := middleware.RequireSuperAdmin(deps.RBACService)

	r.Route("/roles", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(adminMiddleware)
			r.Get("/", h.List)
			r.Get("/{id}", h.Get)
		})

		r.Group(func(r chi.Router) {
			r.Use(superAdminMiddleware)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
			r.Post("/{id}/permissions", h.AssignPermissions)
			r.Put("/{id}/permissions", h.ReplacePermissions)
		})
	})
}

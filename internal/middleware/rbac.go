package middleware

import (
	"net/http"
	"passenger_service_backend/internal/services"
	"strconv"
	"strings"
)

func RequirePermission(
	rbac services.RBACService,
	resource string,
	action string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r)
			if userID == nil {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
				return
			}

			ok, err := rbac.CheckPermission(*userID, resource, action)
			if err != nil || !ok {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireRole(
	rbac services.RBACService,
	roleName string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r)
			if userID == nil {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
				return
			}

			ok, err := rbac.HasRole(*userID, roleName)
			if err != nil || !ok {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient role privileges")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequirePermissionOrOwn(
	rbac services.RBACService,
	resource string,
	action string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r)
			if userID == nil {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
				return
			}

			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			idStr := parts[len(parts)-1]

			resourceID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				sendError(w, http.StatusBadRequest, "bad_request", "Invalid resource ID")
				return
			}

			ok, err := rbac.CheckPermissionOrOwn(*userID, resource, action, uint(resourceID))
			if err != nil || !ok {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdminArea requires level >= 3 (admin or super_admin).
// FIX: Previously checked role IDs 1 and 2 (inverted). Role levels in the system:
//   - super_admin = level 4
//   - admin       = level 3
//   - agent       = level 2
//   - customer    = level 1
func RequireAdminArea(
	rbac services.RBACService,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r)
			if userID == nil {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
				return
			}

			role, err := rbac.GetUserRole(*userID)
			if err != nil {
				sendError(w, http.StatusForbidden, "forbidden", "Could not verify role")
				return
			}

			if role.Level < 3 {
				sendError(w, http.StatusForbidden, "forbidden", "Admin access only")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireSuperAdmin(
	rbac services.RBACService,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r)
			if userID == nil {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
				return
			}

			role, err := rbac.GetUserRole(*userID)
			if err != nil {
				sendError(w, http.StatusForbidden, "forbidden", "Could not verify role")
				return
			}

			if role.Level < 4 {
				sendError(w, http.StatusForbidden, "forbidden", "Super admin access only")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAgentOrAbove(
	rbac services.RBACService,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r)
			if userID == nil {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
				return
			}

			role, err := rbac.GetUserRole(*userID)
			if err != nil {
				sendError(w, http.StatusForbidden, "forbidden", "Could not verify role")
				return
			}

			if role.Level < 2 {
				sendError(w, http.StatusForbidden, "forbidden", "Agent access or above required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"net/http"
	"passenger_service_backend/internal/services"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// =======================
// REQUIRE PERMISSION
// =======================

func RequirePermission(
	rbac services.RBACService,
	resource string,
	action string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := GetUserIDFromContext(r)
			if userID == &uuid.Nil {
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

// =======================
// REQUIRE ROLE (HIERARCHY)
// =======================

func RequireRole(
	rbac services.RBACService,
	roleName string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := GetUserIDFromContext(r)

			ok, err := rbac.HasRole(*userID, roleName)
			if err != nil || !ok {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient role privileges")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// =======================
// PERMISSION OR OWN
// =======================

func RequirePermissionOrOwn(
	rbac services.RBACService,
	resource string,
	action string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := GetUserIDFromContext(r)

			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			idStr := parts[len(parts)-1]

			resourceID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				sendError(w, http.StatusBadRequest, "bad_request", "Invalid resource ID")
				return
			}

			ok, err := rbac.CheckPermissionOrOwn(
				*userID,
				resource,
				action,
				uint(resourceID),
			)

			if err != nil || !ok {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAdminArea(
	rbac services.RBACService,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := GetUserIDFromContext(r)

			isAdmin, _ := rbac.HasExactRole(*userID, 2)
			isSuper, _ := rbac.HasExactRole(*userID, 1)

			if !isAdmin && !isSuper {
				sendError(w, 403, "forbidden", "Admin access only")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

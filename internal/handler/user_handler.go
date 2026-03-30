package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/middleware"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	svc services.UserService
}

func NewUserHandler(svc services.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// POST /users  [admin]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.svc.Create(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", user)
}

// GET /users  [admin]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := utils.PageLimit(r)
	users, err := h.svc.GetAll(r.Context(), page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", users)
}

// GET /users/{uid}  [admin]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	uid, ok := parseUID(w, r, "uid")
	if !ok {
		return
	}
	user, err := h.svc.GetByUID(r.Context(), uid)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", user)
}

// PUT /users/{uid}  [admin]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	uid, ok := parseUID(w, r, "uid")
	if !ok {
		return
	}
	var req dto.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	user, err := h.svc.Update(r.Context(), uid, req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", user)
}

// DELETE /users/{uid}  [admin]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, ok := parseUID(w, r, "uid")
	if !ok {
		return
	}
	err := h.svc.Delete(r.Context(), uid)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", nil)
}

// PUT /users/me/profile  [AuthRequired]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserIDFromContext(r)
	if uid == &uuid.Nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}
	var req dto.UpdateProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	user, err := h.svc.UpdateProfile(r.Context(), *uid, req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", user)
}

// ─────────────────────────────────────────────
// Shared param helpers
// ─────────────────────────────────────────────

func parseUID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	return utils.UUIDParam(w, r, key)
}

func parseUint(w http.ResponseWriter, r *http.Request, key string) (uint, bool) {
	raw := chi.URLParam(r, key)
	n, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, key+" harus berupa angka yang valid", fmt.Errorf("%s", key+" harus berupa angka yang valid"))
		return 0, false
	}
	return uint(n), true
}

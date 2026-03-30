package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/middleware"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"

	"github.com/google/uuid"
)

type AuthHandler struct {
	ctx         context.Context
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService, ctx context.Context) *AuthHandler {
	return &AuthHandler{ctx: ctx,
		authService: authService,
	}
}

// SignIn - POST /api/auth/login
// @Summary SignIn user
// @Description Login with email and password to get access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} utils.Response{data=models.TokenResponse} "Login successful"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := h.authService.Login(h.ctx, req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Login successful", resp)
}

// SignIn - POST /api/auth/admin/login
// @Summary SignIn Admin
// @Description Login with email and password to get access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} utils.Response{data=models.TokenResponse} "Login successful"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := h.authService.LoginAdmin(h.ctx, req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Login successful", resp)
}

// POST /auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	if body.RefreshToken == "" {
		utils.WriteError(w, http.StatusBadRequest, "refresh_token diperlukan", utils.ErrTokenInvalid)
		return
	}
	result, err := h.authService.RefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Login successful", result)
}

// POST /auth/logout  [AuthRequired]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	uid := middleware.GetUserIDFromContext(r)
	if uid == &uuid.Nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}
	err := h.authService.Logout(r.Context(), *uid)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Logout successful", nil)
}

// GET /auth/me  [AuthRequired]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserIDFromContext(r)
	if uid == &uuid.Nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}
	user, err := h.authService.Me(r.Context(), *uid)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Logout successful", user)
}

// PUT /auth/change-password  [AuthRequired]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserIDFromContext(r)
	if uid == &uuid.Nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}
	var req dto.ChangePasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	err := h.authService.ChangePassword(r.Context(), *uid, req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Change password successful", nil)

}

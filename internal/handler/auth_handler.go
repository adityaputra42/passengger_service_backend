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
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService, ctx context.Context) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SignIn godoc
// @Summary      Login penumpang/customer
// @Description  Login dengan email dan password. Mengembalikan access token dan refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Kredensial login"
// @Success      200      {object}  utils.Response{data=dto.AuthResponseDTO}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response  "Email atau password salah"
// @Router       /auth/login [post]
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := h.authService.Login(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Login successful", resp)
}

// AdminLogin godoc
// @Summary      Login admin
// @Description  Login khusus untuk admin (level >= 3). Customer tidak bisa login di sini.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Kredensial login admin"
// @Success      200      {object}  utils.Response{data=dto.AuthResponseDTO}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response  "Bukan admin atau kredensial salah"
// @Router       /auth/admin/login [post]
func (h *AuthHandler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := h.authService.LoginAdmin(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Login successful", resp)
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Mendapatkan access token baru menggunakan refresh token yang masih valid.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body  object{refresh_token=string}  true  "Refresh token"
// @Success      200  {object}  utils.Response{data=dto.AuthResponseDTO}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response  "Refresh token tidak valid atau expired"
// @Router       /auth/refresh [post]
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
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Login successful", result)
}

// Logout godoc
// @Summary      Logout
// @Description  Logout dari sesi yang sedang aktif.
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Security     BearerAuth
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	uid := middleware.GetUserIDFromContext(r)
	if uid == &uuid.Nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}
	err := h.authService.Logout(r.Context(), *uid)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Logout successful", nil)
}

// Me godoc
// @Summary      Profil user yang sedang login
// @Description  Mengambil data user berdasarkan token yang dikirimkan.
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  utils.Response{data=dto.UserResponse}
// @Failure      401  {object}  utils.Response
// @Security     BearerAuth
// @Router       /auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserIDFromContext(r)
	if uid == &uuid.Nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}
	user, err := h.authService.Me(r.Context(), *uid)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Logout successful", user)
}

// ChangePassword godoc
// @Summary      Ganti password
// @Description  Mengganti password user yang sedang login.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ChangePasswordRequest  true  "Password lama dan baru"
// @Success      200      {object}  utils.Response
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      422      {object}  utils.Response  "Password lama tidak sesuai"
// @Security     BearerAuth
// @Router       /auth/change-password [put]
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
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Change password successful", nil)

}

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

// Create godoc
// @Summary      Daftarkan user baru
// @Description  Mendaftarkan user baru (customer self-register atau admin membuat akun).
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateUserRequest  true  "Data user"
// @Success      200      {object}  utils.Response{data=models.User}
// @Failure      400      {object}  utils.Response
// @Failure      409      {object}  utils.Response  "Email sudah terdaftar"
// @Router       /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.svc.Create(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", user)
}

// List godoc
// @Summary      Daftar semua user
// @Description  Mengambil seluruh daftar user dengan pagination. Hanya admin ke atas.
// @Tags         User
// @Produce      json
// @Param        page   query     int  false  "Halaman (default: 1)"
// @Param        limit  query     int  false  "Jumlah per halaman (default: 10)"
// @Success      200    {object}  utils.Response{data=dto.UserListResponse}
// @Failure      401    {object}  utils.Response
// @Failure      403    {object}  utils.Response
// @Security     BearerAuth
// @Router       /users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := utils.PageLimit(r)

	fmt.Printf("Page final %d", page)
	fmt.Printf("Limit final %d", limit)

	users, err := h.svc.GetAll(r.Context(), page, limit)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", users)
}

// Get godoc
// @Summary      Detail user
// @Description  Mengambil detail satu user berdasarkan UID. Hanya admin ke atas.
// @Tags         User
// @Produce      json
// @Param        uid  path      string  true  "User UUID"
// @Success      200  {object}  utils.Response{data=models.User}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     BearerAuth
// @Router       /users/{uid} [get]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	uid, ok := parseUID(w, r, "uid")
	if !ok {
		return
	}
	user, err := h.svc.GetByUID(r.Context(), uid)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", user)
}

// Update godoc
// @Summary      Update user
// @Description  Memperbarui data user (nama, role). Hanya admin ke atas.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        uid      path      string                 true  "User UUID"
// @Param        request  body      dto.UpdateUserRequest  true  "Data update"
// @Success      200      {object}  utils.Response{data=models.User}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /users/{uid} [put]
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
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", user)
}

// Delete godoc
// @Summary      Hapus user
// @Description  Menghapus user. Hanya admin ke atas.
// @Tags         User
// @Produce      json
// @Param        uid  path      string  true  "User UUID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     BearerAuth
// @Router       /users/{uid} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, ok := parseUID(w, r, "uid")
	if !ok {
		return
	}
	err := h.svc.Delete(r.Context(), uid)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", nil)
}

// UpdateProfile godoc
// @Summary      Update profil sendiri
// @Description  User yang sedang login memperbarui profil dirinya sendiri (nama lengkap).
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UpdateProfileRequest  true  "Data profil"
// @Success      200      {object}  utils.Response{data=models.User}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Security     BearerAuth
// @Router       /users/me/profile [put]
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
		utils.WriteServiceError(w, err)
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

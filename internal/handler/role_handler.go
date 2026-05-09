package handler

import (
	"encoding/json"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
)

type RoleHandler struct {
	svc services.RoleService
}

func NewRoleHandler(svc services.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

// Create godoc
// @Summary      Buat role baru
// @Description  Membuat role baru beserta permission. Hanya super_admin.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RoleInput  true  "Data role"
// @Success      200      {object}  utils.Response{data=models.Role}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Security     BearerAuth
// @Router       /roles [post]
func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.RoleInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	role, err := h.svc.CreateRole(r.Context(), &req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create role successful", role)
}

// List godoc
// @Summary      Daftar semua role
// @Description  Mengambil semua role beserta permission masing-masing. Hanya admin ke atas.
// @Tags         Role
// @Produce      json
// @Success      200  {object}  utils.Response{data=[]models.Role}
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Security     BearerAuth
// @Router       /roles [get]
func (h *RoleHandler) List(w http.ResponseWriter, r *http.Request) {
	roles, err := h.svc.FindAllRole(r.Context())
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", roles)
}

// Get godoc
// @Summary      Detail role
// @Description  Mengambil detail role beserta seluruh permission-nya. Hanya admin ke atas.
// @Tags         Role
// @Produce      json
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  utils.Response{data=models.Role}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     BearerAuth
// @Router       /roles/{id} [get]
func (h *RoleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUint(w, r, "id")
	if !ok {
		return
	}
	role, err := h.svc.FindById(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", role)
}

// Update godoc
// @Summary      Update role
// @Description  Memperbarui informasi role. Hanya super_admin.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id       path      int           true  "Role ID"
// @Param        request  body      dto.RoleInput  true  "Data update"
// @Success      200      {object}  utils.Response{data=models.Role}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /roles/{id} [put]
func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUint(w, r, "id")
	if !ok {
		return
	}
	var req dto.RoleInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	role, err := h.svc.UpdateRole(r.Context(), id, &req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", role)
}

// Delete godoc
// @Summary      Hapus role
// @Description  Menghapus role. System role dan role yang masih digunakan user tidak bisa dihapus. Hanya super_admin.
// @Tags         Role
// @Produce      json
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     BearerAuth
// @Router       /roles/{id} [delete]
func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUint(w, r, "id")
	if !ok {
		return
	}
	err := h.svc.DeleteRole(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", nil)
}

// AssignPermissions godoc
// @Summary      Tambah permission ke role
// @Description  Menambahkan satu atau lebih permission ke role. Hanya super_admin.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "Role ID"
// @Param        request  body      dto.RolePermissionInput  true  "Daftar permission ID"
// @Success      200      {object}  utils.Response
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Security     BearerAuth
// @Router       /roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUint(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	err := h.svc.AssignPermissions(r.Context(), id, body.PermissionIDs)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", nil)
}

// ReplacePermissions godoc
// @Summary      Ganti permission role
// @Description  Mengganti seluruh permission role dengan daftar baru. Hanya super_admin.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "Role ID"
// @Param        request  body      dto.RolePermissionInput  true  "Daftar permission ID baru"
// @Success      200      {object}  utils.Response
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Security     BearerAuth
// @Router       /roles/{id}/permissions [put]
func (h *RoleHandler) ReplacePermissions(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUint(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	err := h.svc.AssignPermissions(r.Context(), id, body.PermissionIDs)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", nil)
}

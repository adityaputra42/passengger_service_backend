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

// POST /roles  [super_admin]
func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.RoleInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	role, err := h.svc.CreateRole(r.Context(), &req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create role successful", role)
}

// GET /roles  [admin]
func (h *RoleHandler) List(w http.ResponseWriter, r *http.Request) {
	roles, err := h.svc.FindAllRole(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", roles)
}

// GET /roles/{id}  [admin]
func (h *RoleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUint(w, r, "id")
	if !ok {
		return
	}
	role, err := h.svc.FindById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", role)
}

// PUT /roles/{id}  [super_admin]
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
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", role)
}

// DELETE /roles/{id}  [super_admin]
func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUint(w, r, "id")
	if !ok {
		return
	}
	err := h.svc.DeleteRole(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", nil)
}

// POST /roles/{id}/permissions  [super_admin]
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
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", nil)
}

// PUT /roles/{id}/permissions  [super_admin]
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
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Create user successful", nil)
}

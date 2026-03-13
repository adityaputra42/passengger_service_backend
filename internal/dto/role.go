package dto

import (
	"passenger_service_backend/internal/models"
	"time"
)

type CreateRoleRequest struct {
	Name         string `json:"name"           validate:"required,min=2,max=50"`
	Description  string `json:"description"    validate:"omitempty,max=255"`
	Level        int    `json:"level"          validate:"required,min=1,max=10"`
	IsSystemRole bool   `json:"is_system_role"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"        validate:"omitempty,min=2,max=50"`
	Description string `json:"description" validate:"omitempty,max=255"`
	Level       int    `json:"level"       validate:"omitempty,min=1,max=10"`
}

type RoleInput struct {
	Name          string `json:"name" validate:"required,min=2,max=50"`
	Description   string `json:"description" validate:"max=500"`
	Level         int    `json:"level" validate:"required,min=1"`
	IsSystemRole  bool   `json:"is_system_role"`
	PermissionIDs []uint `json:"permission_ids"`
}

type RoleWithPermissions struct {
	models.Role
	PermissionIDs []uint `json:"permission_ids"`
}

type RolePermissionInput struct {
	PermissionIDs []uint `json:"permission_ids" validate:"required,min=1"`
}

type PermissionResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToPermissionResponse(p *models.Permission) *PermissionResponse {
	if p == nil {
		return nil
	}
	return &PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

func ToPermissionResponseList(perms []*models.Permission) []PermissionResponse {
	out := make([]PermissionResponse, 0, len(perms))
	for _, p := range perms {
		if r := ToPermissionResponse(p); r != nil {
			out = append(out, *r)
		}
	}
	return out
}

type RoleResponse struct {
	ID           uint                 `json:"id"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Level        int                  `json:"level"`
	IsSystemRole bool                 `json:"is_system_role"`
	Permissions  []PermissionResponse `json:"permissions"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

type RoleListResponse struct {
	Roles []RoleResponse `json:"roles"`
	PaginationMeta
}

func ToRoleResponse(r *models.Role) *RoleResponse {
	if r == nil {
		return nil
	}
	return &RoleResponse{
		ID:           r.ID,
		Name:         r.Name,
		Description:  r.Description,
		Level:        r.Level,
		IsSystemRole: r.IsSystemRole,
		Permissions:  ToPermissionResponseList(r.Permissions),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

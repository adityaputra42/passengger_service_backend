package dto

import (
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"omitempty,min=2,max=255"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	FullName string `json:"full_name" validate:"omitempty,min=1,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	RoleID   uint   `json:"role_id" validate:"required,min=1"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name" validate:"omitempty,min=2,max=255"`
	RoleID   uint   `json:"role_id"   validate:"omitempty"`
}

type PasswordUpdateInput struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type UserListRequest struct {
	UserId *uuid.UUID
	Limit  int
	Page   int
	SortBy string
}

type UserResponse struct {
	UID      uuid.UUID `json:"uid"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	RoleID   uint      `json:"role_id"`
	// Role        RoleResponse `json:"role"`
	// Permissions []string  `json:"permissions"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	PaginationMeta
}

func ToUserResponse(u *models.User) *UserResponse {
	if u == nil {
		return nil
	}
	permissions := make([]string, 0)
	for _, p := range u.Role.Permissions {
		permissions = append(permissions, p.Name)
	}
	// role := RoleResponse{}
	// if r := ToRoleResponse(&u.Role); r != nil {
	// 	role = *r
	// }
	return &UserResponse{
		UID:      u.UID,
		Email:    u.Email,
		FullName: u.FullName,
		RoleID:   u.RoleID,
		// Role:        role,
		// Permissions: permissions,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToUserListResponse(users []models.User, total int64, page, limit int) *UserListResponse {
	out := make([]UserResponse, 0, len(users))
	for i := range users {
		if r := ToUserResponse(&users[i]); r != nil {
			out = append(out, *r)
		}
	}
	return &UserListResponse{
		Users:          out,
		PaginationMeta: newPagination(total, page, limit),
	}
}

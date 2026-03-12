package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type User struct {
	UID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"uid"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"           json:"email"          validate:"required,email,max=255"`
	FullName     string    `gorm:"type:varchar(255);not null"                       json:"full_name"      validate:"required,min=2,max=255"`
	PasswordHash string    `gorm:"type:varchar(255);not null"                       json:"-"`
	RoleID       uint      `gorm:"not null;index"                                   json:"role_id"        validate:"required"`
	CreatedAt    time.Time `gorm:"autoCreateTime"                                   json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"                                   json:"updated_at"`
	Role         Role      `gorm:"foreignKey:RoleID"                                json:"role,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UID == uuid.Nil {
		u.UID = uuid.New()
	}
	return nil
}


type UserInput struct {
	Email     string `json:"email" validate:"omitempty,email"`
	FullName string `json:"full_name" validate:"omitempty,min=1,max=50"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
	RoleID    uint   `json:"role_id" validate:"required,min=1"`
}

type UserUpdateInput struct {
	Email     string `json:"email" validate:"omitempty,email"`
	FullName string `json:"full_name" validate:"omitempty,min=1,max=50"`
	RoleID    uint   `json:"role_id" validate:"omitempty,min=1"`
}

type PasswordUpdateInput struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type UserResponse struct {
	UID             uuid.UUID  `json:"uid"`
	Email           string     `json:"email"`
	FullName       string     `json:"full_name"`
	RoleID          uint       `json:"role_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Role            Role       `json:"role"`
	Permissions     []string   `json:"permissions"`
}

type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type UserListRequest struct {
	UserId *uuid.UUID
	Limit  int
	Page   int
	SortBy string
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToResponse() *UserResponse {
	permissions := []string{}
	if u.Role.Permissions != nil {
		for _, p := range u.Role.Permissions {
			permissions = append(permissions, p.Name)
		}
	}

	return &UserResponse{
		UID:              u.UID,
		Email:           u.Email,
		FullName:       u.FullName,
		RoleID:          u.RoleID,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		Role:            u.Role,
		Permissions:     permissions,
	}
}

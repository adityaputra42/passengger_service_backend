package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// User
// ─────────────────────────────────────────────



type User struct {
	UID          uuid.UUID `gorm:"type:uuid;primaryKey;column:uid" json:"uid"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null;column:email" json:"email" validate:"required,email,max=255"`
	FullName     string    `gorm:"type:varchar(255);not null;column:full_name" json:"full_name" validate:"required,min=2,max=255"`
	PasswordHash string    `gorm:"type:varchar(255);not null;column:password_hash" json:"-"`
	RoleID       uint      `json:"role_id" gorm:"not null;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	Role         Role      `json:"role" gorm:"foreignKey:RoleID"`

	// Relations
	Bookings []Booking `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
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
	UserId *uint
	Limit  int
	Page   int
	SortBy string
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToResponse() *UserResponse {
	permissions := []string{}
	// Ensure Role and Permissions are preloaded
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

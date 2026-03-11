package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// User
// ─────────────────────────────────────────────

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
)

type User struct {
	UID          uuid.UUID `gorm:"type:uuid;primaryKey;column:uid" json:"uid"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null;column:email" json:"email" validate:"required,email,max=255"`
	FullName     string    `gorm:"type:varchar(255);not null;column:full_name" json:"full_name" validate:"required,min=2,max=255"`
	PasswordHash string    `gorm:"type:varchar(255);not null;column:password_hash" json:"-"`
	Role         UserRole  `gorm:"type:varchar(255);not null;column:role" json:"role" validate:"required,oneof=admin customer"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`

	// Relations
	Bookings []Booking `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UID == uuid.Nil {
		u.UID = uuid.New()
	}
	return nil
}


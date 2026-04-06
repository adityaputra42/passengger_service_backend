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

func (User) TableName() string {
	return "users"
}

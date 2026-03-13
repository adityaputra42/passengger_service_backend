package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"type:varchar(50);uniqueIndex;not null" validate:"required,min=2,max=50"`
	Description  string         `json:"description" gorm:"type:text"`
	Level        int            `json:"level" gorm:"not null;default:1;index" validate:"required,min=1"`
	IsSystemRole bool           `json:"is_system_role" gorm:"not null;default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	Users        []User         `json:"users,omitempty" gorm:"foreignKey:RoleID"`
	Permissions  []*Permission  `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}


func (Role) TableName() string {
	return "roles"
}

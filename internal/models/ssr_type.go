package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type SSRType struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code string    `gorm:"type:varchar(10)"                                 json:"code" validate:"max=10"`
	Name string    `gorm:"type:varchar(100)"                                json:"name" validate:"max=100"`
}

func (s *SSRType) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

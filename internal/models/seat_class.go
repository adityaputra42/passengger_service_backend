package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type SeatClass struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code string    `gorm:"type:varchar(10)"                                 json:"code" validate:"max=10"`
	Name string    `gorm:"type:varchar(50)"                                 json:"name" validate:"max=50"`
}

func (s *SeatClass) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

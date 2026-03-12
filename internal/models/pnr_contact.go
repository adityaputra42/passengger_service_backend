package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type PNRContact struct {
	ID    uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PNRID *uuid.UUID `gorm:"type:uuid;index"                                  json:"pnr_id"`
	Name  string     `gorm:"type:varchar(255)"                                json:"name"  validate:"max=255"`
	Email string     `gorm:"type:varchar(255)"                                json:"email" validate:"omitempty,email,max=255"`
	Phone string     `gorm:"type:varchar(50)"                                 json:"phone" validate:"max=50"`
	PNR   *PNR       `gorm:"foreignKey:PNRID"                                 json:"pnr,omitempty"`
}

func (p *PNRContact) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

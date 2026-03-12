package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type PNRSegment struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PNRID        *uuid.UUID `gorm:"type:uuid;index"                                  json:"pnr_id"`
	FlightID     *uuid.UUID `gorm:"type:uuid"                                        json:"flight_id"`
	SegmentOrder int        `gorm:"default:1"                                        json:"segment_order"`
	PNR          *PNR       `gorm:"foreignKey:PNRID"                                 json:"pnr,omitempty"`
	Flight       *Flight    `gorm:"foreignKey:FlightID"                              json:"flight,omitempty"`
}

func (p *PNRSegment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PassengerSSR struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PassengerID *uuid.UUID    `gorm:"type:uuid;index"                                  json:"passenger_id"`
	SegmentID   *uuid.UUID    `gorm:"type:uuid"                                        json:"segment_id"`
	SSRTypeID   *uuid.UUID    `gorm:"type:uuid"                                        json:"ssr_type_id"`
	Passenger   *PNRPassenger `gorm:"foreignKey:PassengerID"                           json:"passenger,omitempty"`
	Segment     *PNRSegment   `gorm:"foreignKey:SegmentID"                             json:"segment,omitempty"`
	SSRType     *SSRType      `gorm:"foreignKey:SSRTypeID"                             json:"ssr_type,omitempty"`
}

func (p *PassengerSSR) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

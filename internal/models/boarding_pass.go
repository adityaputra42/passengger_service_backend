package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardingPass struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PassengerID   *uuid.UUID    `gorm:"type:uuid;index"                                  json:"passenger_id"`
	SegmentID     *uuid.UUID    `gorm:"type:uuid"                                        json:"segment_id"`
	BoardingGroup string        `gorm:"type:varchar(10)"                                 json:"boarding_group" validate:"max=10"`
	Gate          string        `gorm:"type:varchar(10)"                                 json:"gate"           validate:"max=10"`
	BoardingTime  *time.Time    `gorm:"type:timestamptz"                                 json:"boarding_time"`
	QRCode        string        `gorm:"type:text"                                        json:"qr_code"`
	Passenger     *PNRPassenger `gorm:"foreignKey:PassengerID"                           json:"passenger,omitempty"`
	Segment       *PNRSegment   `gorm:"foreignKey:SegmentID"                             json:"segment,omitempty"`
}

func (b *BoardingPass) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

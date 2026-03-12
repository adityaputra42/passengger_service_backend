package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Checkin struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PassengerID *uuid.UUID    `gorm:"type:uuid;index"                                  json:"passenger_id"`
	SegmentID   *uuid.UUID    `gorm:"type:uuid"                                        json:"segment_id"`
	CheckinTime *time.Time    `gorm:"type:timestamptz"                                 json:"checkin_time"`
	Passenger   *PNRPassenger `gorm:"foreignKey:PassengerID"                           json:"passenger,omitempty"`
	Segment     *PNRSegment   `gorm:"foreignKey:SegmentID"                             json:"segment,omitempty"`
}

func (c *Checkin) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

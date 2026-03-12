package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaggageStatus string

const (
	BaggageCheckedIn BaggageStatus = "checked_in"
	BaggageLoaded    BaggageStatus = "loaded"
	BaggageDelivered BaggageStatus = "delivered"
	BaggageLost      BaggageStatus = "lost"
)

type Baggage struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PassengerID *uuid.UUID    `gorm:"type:uuid;index"                                  json:"passenger_id"`
	SegmentID   *uuid.UUID    `gorm:"type:uuid"                                        json:"segment_id"`
	Weight      float64       `gorm:"type:numeric(5,2)"                                json:"weight"      validate:"min=0"`
	TagNumber   string        `gorm:"type:varchar(50)"                                 json:"tag_number"  validate:"max=50"`
	Status      BaggageStatus `gorm:"type:varchar(20)"                                 json:"status"      validate:"oneof=checked_in loaded delivered lost"`
	Passenger   *PNRPassenger `gorm:"foreignKey:PassengerID"                           json:"passenger,omitempty"`
	Segment     *PNRSegment   `gorm:"foreignKey:SegmentID"                             json:"segment,omitempty"`
}

func (b *Baggage) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

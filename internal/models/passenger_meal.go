package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type PassengerMeal struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PassengerID *uuid.UUID    `gorm:"type:uuid;index"                                  json:"passenger_id"`
	SegmentID   *uuid.UUID    `gorm:"type:uuid"                                        json:"segment_id"`
	MealID      *uuid.UUID    `gorm:"type:uuid"                                        json:"meal_id"`
	Passenger   *PNRPassenger `gorm:"foreignKey:PassengerID"                           json:"passenger,omitempty"`
	Segment     *PNRSegment   `gorm:"foreignKey:SegmentID"                             json:"segment,omitempty"`
	Meal        *Meal         `gorm:"foreignKey:MealID"                                json:"meal,omitempty"`
}

func (p *PassengerMeal) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

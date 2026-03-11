package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// FlightSeat
// ─────────────────────────────────────────────

type FlightSeat struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	FlightID uuid.UUID `gorm:"type:uuid;not null;column:flight_id" json:"flight_id" validate:"required"`
	SeatID   uuid.UUID `gorm:"type:uuid;not null;column:seat_id" json:"seat_id" validate:"required"`
	Price    float64   `gorm:"not null;column:price" json:"price" validate:"required,gt=0"`
	// true = available, false = booked
	Status bool `gorm:"not null;column:status" json:"status"`

	// Relations
	Flight Flight `gorm:"foreignKey:FlightID" json:"flight,omitempty"`
	Seat   Seat   `gorm:"foreignKey:SeatID" json:"seat,omitempty"`
}

func (fs *FlightSeat) BeforeCreate(tx *gorm.DB) error {
	if fs.ID == uuid.Nil {
		fs.ID = uuid.New()
	}
	return nil
}

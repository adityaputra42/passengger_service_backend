package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// Seat
// ─────────────────────────────────────────────


type Seat struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	AircraftID uuid.UUID `gorm:"type:uuid;not null;column:aircraft_id" json:"aircraft_id" validate:"required"`
	SeatNumber string    `gorm:"type:varchar(255);not null;column:seat_number" json:"seat_number" validate:"required,max=10"`
	SeatClass  SeatClass `gorm:"type:varchar(255);not null;column:seat_class" json:"seat_class" validate:"required,oneof=economy business first"`

	// Relations
	Aircraft    Aircraft     `gorm:"foreignKey:AircraftID" json:"aircraft,omitempty"`
	FlightSeats []FlightSeat `gorm:"foreignKey:SeatID" json:"flight_seats,omitempty"`
}

func (s *Seat) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

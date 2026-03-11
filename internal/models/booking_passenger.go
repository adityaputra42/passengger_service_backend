package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// BookingPassenger  (join table)
// ─────────────────────────────────────────────

type BookingPassenger struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	BookingID    uuid.UUID `gorm:"type:uuid;not null;column:booking_id" json:"booking_id" validate:"required"`
	PassengerID  uuid.UUID `gorm:"type:uuid;not null;column:passenger_id" json:"passenger_id" validate:"required"`
	FlightSeatID uuid.UUID `gorm:"type:uuid;not null;column:flight_seat_id" json:"flight_seat_id" validate:"required"`

	// Relations
	Booking    Booking    `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Passenger  Passenger  `gorm:"foreignKey:PassengerID" json:"passenger,omitempty"`
	FlightSeat FlightSeat `gorm:"foreignKey:FlightSeatID" json:"flight_seat,omitempty"`
}

func (bp *BookingPassenger) BeforeCreate(tx *gorm.DB) error {
	if bp.ID == uuid.Nil {
		bp.ID = uuid.New()
	}
	return nil
}

func (BookingPassenger) TableName() string { return "Booking Passengger" }

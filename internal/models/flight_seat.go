package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// FlightSeat
// ─────────────────────────────────────────────

type FlightSeatStatus string

const (
	FlightSeatAvailable FlightSeatStatus = "available"
	FlightSeatLocked    FlightSeatStatus = "locked"
	FlightSeatBooked    FlightSeatStatus = "booked"
)

type FlightSeat struct {
	ID             uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	FlightID       *uuid.UUID       `gorm:"type:uuid;index"                                  json:"flight_id"`
	AircraftSeatID *uuid.UUID       `gorm:"type:uuid"                                        json:"aircraft_seat_id"`
	Price          float64          `gorm:"type:numeric(10,2)"                               json:"price"   validate:"min=0"`
	Status         FlightSeatStatus `gorm:"type:varchar(20);default:available"               json:"status"  validate:"oneof=available locked booked"`
	Flight         *Flight          `gorm:"foreignKey:FlightID"                              json:"flight,omitempty"`
	AircraftSeat   *AircraftSeat    `gorm:"foreignKey:AircraftSeatID"                        json:"aircraft_seat,omitempty"`
}

func (fs *FlightSeat) BeforeCreate(tx *gorm.DB) error {
	if fs.ID == uuid.Nil {
		fs.ID = uuid.New()
	}
	return nil
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// Flight
// ─────────────────────────────────────────────

type FlightStatus string

const (
	FlightScheduled FlightStatus = "scheduled"
	FlightBoarding  FlightStatus = "boarding"
	FlightDeparted  FlightStatus = "departed"
	FlightArrived   FlightStatus = "arrived"
	FlightCancelled FlightStatus = "cancelled"
	FlightDelayed   FlightStatus = "delayed"
)

type Flight struct {
	ID                 uuid.UUID    `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	FlightNumber       string       `gorm:"type:varchar(20);not null;column:flight_number" json:"flight_number" validate:"required,max=20"`
	AircraftID         uuid.UUID    `gorm:"type:uuid;not null;column:aircraft_id" json:"aircraft_id" validate:"required"`
	DepartureAirportID uuid.UUID    `gorm:"type:uuid;not null;column:departure_airport_id" json:"departure_airport_id" validate:"required"`
	ArrivalAirportID   uuid.UUID    `gorm:"type:uuid;not null;column:arrival_airport_id" json:"arrival_airport_id" validate:"required"`
	DepartureTime      time.Time    `gorm:"not null;column:departure_time" json:"departure_time" validate:"required"`
	ArrivalTime        time.Time    `gorm:"not null;column:arrival_time" json:"arrival_time" validate:"required,gtfield=DepartureTime"`
	BasePrice          float64      `gorm:"not null;column:base_price" json:"base_price" validate:"required,gt=0"`
	Status             FlightStatus `gorm:"type:varchar(20);not null;column:status" json:"status" validate:"required,oneof=scheduled boarding departed arrived cancelled delayed"`

	// Relations
	Aircraft         Aircraft    `gorm:"foreignKey:AircraftID" json:"aircraft,omitempty"`
	DepartureAirport Airport     `gorm:"foreignKey:DepartureAirportID" json:"departure_airport,omitempty"`
	ArrivalAirport   Airport     `gorm:"foreignKey:ArrivalAirportID" json:"arrival_airport,omitempty"`
	FlightSeats      []FlightSeat `gorm:"foreignKey:FlightID" json:"flight_seats,omitempty"`
	Bookings         []Booking   `gorm:"foreignKey:FlightID" json:"bookings,omitempty"`
}

func (f *Flight) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

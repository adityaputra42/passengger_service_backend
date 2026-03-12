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
	FlightStatusScheduled FlightStatus = "scheduled"
	FlightStatusBoarding  FlightStatus = "boarding"
	FlightStatusDeparted  FlightStatus = "departed"
	FlightStatusArrived   FlightStatus = "arrived"
	FlightStatusCancelled FlightStatus = "cancelled"
	FlightStatusDelayed   FlightStatus = "delayed"
)

type Flight struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ScheduleID    *uuid.UUID      `gorm:"type:uuid"                                        json:"schedule_id"`
	AircraftID    *uuid.UUID      `gorm:"type:uuid"                                        json:"aircraft_id"`
	DepartureTime *time.Time      `gorm:"type:timestamptz;index"                           json:"departure_time"`
	ArrivalTime   *time.Time      `gorm:"type:timestamptz"                                 json:"arrival_time"`
	Status        FlightStatus    `gorm:"type:varchar(20)"                                 json:"status" validate:"oneof=scheduled boarding departed arrived cancelled delayed"`
	Schedule      *FlightSchedule `gorm:"foreignKey:ScheduleID"                            json:"schedule,omitempty"`
	Aircraft      *Aircraft       `gorm:"foreignKey:AircraftID"                            json:"aircraft,omitempty"`
	FlightSeats   []FlightSeat    `gorm:"foreignKey:FlightID"                              json:"flight_seats,omitempty"`
}

func (f *Flight) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

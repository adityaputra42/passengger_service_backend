package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatAssignment struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                    json:"id"`
	PassengerID  *uuid.UUID  `gorm:"type:uuid"                                                           json:"passenger_id"`
	SegmentID    *uuid.UUID  `gorm:"type:uuid;uniqueIndex:unique_seat_segment"                           json:"segment_id"`
	FlightSeatID *uuid.UUID  `gorm:"type:uuid;uniqueIndex:unique_seat_segment"                           json:"flight_seat_id"`
	AssignedAt   *time.Time  `gorm:"type:timestamptz"                                                    json:"assigned_at"`
	Passenger    *PNRPassenger `gorm:"foreignKey:PassengerID"                                           json:"passenger,omitempty"`
	Segment      *PNRSegment `gorm:"foreignKey:SegmentID"                                               json:"segment,omitempty"`
	FlightSeat   *FlightSeat `gorm:"foreignKey:FlightSeatID"                                            json:"flight_seat,omitempty"`
}

func (s *SeatAssignment) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type SeatLock struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	FlightSeatID *uuid.UUID  `gorm:"type:uuid;index"                                  json:"flight_seat_id"`
	PNRID        *uuid.UUID  `gorm:"type:uuid"                                        json:"pnr_id"`
	LockedAt     *time.Time  `gorm:"type:timestamptz"                                 json:"locked_at"`
	ExpiresAt    *time.Time  `gorm:"type:timestamptz;index"                           json:"expires_at"`
	FlightSeat   *FlightSeat `gorm:"foreignKey:FlightSeatID"                          json:"flight_seat,omitempty"`
	PNR          *PNR        `gorm:"foreignKey:PNRID"                                 json:"pnr,omitempty"`
}

func (s *SeatLock) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

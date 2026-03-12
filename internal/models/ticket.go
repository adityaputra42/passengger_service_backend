package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Ticket struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PassengerID  *uuid.UUID      `gorm:"type:uuid"                                        json:"passenger_id"`
	TicketNumber string          `gorm:"type:varchar(20);uniqueIndex"                     json:"ticket_number" validate:"max=20"`
	IssuedAt     *time.Time      `gorm:"type:timestamptz"                                 json:"issued_at"`
	Passenger    *PNRPassenger   `gorm:"foreignKey:PassengerID"                           json:"passenger,omitempty"`
	Segments     []TicketSegment `gorm:"foreignKey:TicketID"                              json:"segments,omitempty"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

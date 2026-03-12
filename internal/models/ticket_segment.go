package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type TicketSegment struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TicketID  *uuid.UUID  `gorm:"type:uuid;index"                                  json:"ticket_id"`
	SegmentID *uuid.UUID  `gorm:"type:uuid"                                        json:"segment_id"`
	Ticket    *Ticket     `gorm:"foreignKey:TicketID"                              json:"ticket,omitempty"`
	Segment   *PNRSegment `gorm:"foreignKey:SegmentID"                             json:"segment,omitempty"`
}

func (t *TicketSegment) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

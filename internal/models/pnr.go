package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type PNRStatus string

const (
	PNRStatusHold      PNRStatus = "hold"
	PNRStatusConfirmed PNRStatus = "confirmed"
	PNRStatusCancelled PNRStatus = "cancelled"
	PNRStatusTicketed  PNRStatus = "ticketed"
)

type PNR struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	RecordLocator string         `gorm:"type:varchar(6);uniqueIndex;not null"             json:"record_locator" validate:"required,len=6"`
	Status        PNRStatus      `gorm:"type:varchar(20)"                                 json:"status"         validate:"oneof=hold confirmed cancelled ticketed"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"                                   json:"created_at"`
	TTL           *time.Time     `gorm:"type:timestamptz"                                 json:"ttl"`
	Contact       *PNRContact    `gorm:"foreignKey:PNRID"                                 json:"contact,omitempty"`
	Passengers    []PNRPassenger `gorm:"foreignKey:PNRID"                                 json:"passengers,omitempty"`
	Segments      []PNRSegment   `gorm:"foreignKey:PNRID"                                 json:"segments,omitempty"`
	Payments      []Payment      `gorm:"foreignKey:PNRID"                                 json:"payments,omitempty"`
}

func (p *PNR) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (PNR) TableName() string { return "pnrs" }

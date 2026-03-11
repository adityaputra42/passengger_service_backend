package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// Airport
// ─────────────────────────────────────────────

type Airport struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	Code      string    `gorm:"type:varchar(3);not null;column:code" json:"code" validate:"required,len=3,uppercase"`
	Name      string    `gorm:"type:varchar(255);not null;column:name" json:"name" validate:"required,max=255"`
	City      string    `gorm:"type:varchar(255);not null;column:city" json:"city" validate:"required,max=255"`
	Country   string    `gorm:"type:varchar(255);not null;column:country" json:"country" validate:"required,max=255"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`

	// Relations
	DepartureFlights []Flight `gorm:"foreignKey:DepartureAirportID" json:"departure_flights,omitempty"`
	ArrivalFlights   []Flight `gorm:"foreignKey:ArrivalAirportID" json:"arrival_flights,omitempty"`
}

func (a *Airport) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

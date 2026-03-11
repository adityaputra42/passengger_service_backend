package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// Aircraft
// ─────────────────────────────────────────────

type Aircraft struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	Model          string    `gorm:"type:varchar(255);not null;column:model" json:"model" validate:"required,max=255"`
	RegistrationNo string    `gorm:"type:varchar(255);not null;column:registration_no" json:"registration_no" validate:"required,max=255"`
	TotalSeats     int64     `gorm:"not null;column:total_seats" json:"total_seats" validate:"required,min=1"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`

	// Relations
	Seats   []Seat   `gorm:"foreignKey:AircraftID" json:"seats,omitempty"`
	Flights []Flight `gorm:"foreignKey:AircraftID" json:"flights,omitempty"`
}

func (a *Aircraft) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

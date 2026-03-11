package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// Passenger  (fixed typo: "Passengger" → "Passenger")
// ─────────────────────────────────────────────

type Passenger struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	FirstName      string    `gorm:"type:varchar(255);not null;column:first_name" json:"first_name" validate:"required,max=255"`
	LastName       string    `gorm:"type:varchar(255);not null;column:last_name" json:"last_name" validate:"required,max=255"`
	DateOfBirth    time.Time `gorm:"type:date;not null;column:date_of_birth" json:"date_of_birth" validate:"required"`
	PassportNumber string    `gorm:"type:varchar(255);not null;column:passport_number" json:"passport_number" validate:"required,max=50,alphanum"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`

	// Relations
	BookingPassengers []BookingPassenger `gorm:"foreignKey:PassengerID" json:"booking_passengers,omitempty"`
}

func (p *Passenger) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName maps to the original (typo'd) table name in the DB
func (Passenger) TableName() string { return "Passengger" }

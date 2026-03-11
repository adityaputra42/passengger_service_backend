package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// Booking
// ─────────────────────────────────────────────

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	UserID      uuid.UUID     `gorm:"type:uuid;not null;column:user_id" json:"user_id" validate:"required"`
	FlightID    uuid.UUID     `gorm:"type:uuid;not null;column:flight_id" json:"flight_id" validate:"required"`
	BookingCode string        `gorm:"type:varchar(255);uniqueIndex;not null;column:booking_code" json:"booking_code" validate:"required,max=255"`
	TotalAmount float64       `gorm:"not null;column:total_amount" json:"total_amount" validate:"required,gt=0"`
	Status      BookingStatus `gorm:"type:varchar(255);not null;column:status" json:"status" validate:"required,oneof=pending confirmed cancelled"`
	CreatedAt   time.Time     `gorm:"autoCreateTime;column:created_at" json:"created_at"`

	// Relations
	User              User               `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Flight            Flight             `gorm:"foreignKey:FlightID" json:"flight,omitempty"`
	Payments          []Payment          `gorm:"foreignKey:BookingID" json:"payments,omitempty"`
	BookingPassengers []BookingPassenger `gorm:"foreignKey:BookingID" json:"booking_passengers,omitempty"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (Booking) TableName() string { return "booking" }

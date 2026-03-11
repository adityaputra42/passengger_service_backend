package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string
type PaymentMethod string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentSuccess  PaymentStatus = "success"
	PaymentFailed   PaymentStatus = "failed"
	PaymentRefunded PaymentStatus = "refunded"
	MethodCreditCard   PaymentMethod = "credit_card"
	MethodDebitCard    PaymentMethod = "debit_card"
	MethodBankTransfer PaymentMethod = "bank_transfer"
	MethodEWallet      PaymentMethod = "e_wallet"
)

type Payment struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	BookingID uuid.UUID     `gorm:"type:uuid;not null;column:booking_id" json:"booking_id" validate:"required"`
	Amount    float64       `gorm:"not null;column:amount" json:"amount" validate:"required,gt=0"`
	Status    PaymentStatus `gorm:"type:varchar(255);not null;column:status" json:"status" validate:"required,oneof=pending success failed refunded"`
	Method    PaymentMethod `gorm:"type:varchar(255);not null;column:method" json:"method" validate:"required,oneof=credit_card debit_card bank_transfer e_wallet"`
	PaidAt    time.Time     `gorm:"column:paid_at" json:"paid_at"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (Payment) TableName() string { return "payment" }

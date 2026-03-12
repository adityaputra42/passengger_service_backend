package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string
type PaymentMethod string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusSuccess  PaymentStatus = "success"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"

	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodDebitCard    PaymentMethod = "debit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodEWallet      PaymentMethod = "e_wallet"
)

type Payment struct {
	ID     uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PNRID  *uuid.UUID    `gorm:"type:uuid;index"                                  json:"pnr_id"`
	Amount float64       `gorm:"type:numeric(10,2)"                               json:"amount" validate:"min=0"`
	Method PaymentMethod `gorm:"type:varchar(50)"                                 json:"method" validate:"oneof=credit_card debit_card bank_transfer e_wallet"`
	Status PaymentStatus `gorm:"type:varchar(20)"                                 json:"status" validate:"oneof=pending success failed refunded"`
	PaidAt *time.Time    `gorm:"type:timestamptz"                                 json:"paid_at"`
	PNR    *PNR          `gorm:"foreignKey:PNRID"                                 json:"pnr,omitempty"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

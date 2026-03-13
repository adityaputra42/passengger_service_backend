package dto

import (
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type InitiatePaymentRequest struct {
	PNRID  uuid.UUID            `json:"pnr_id" validate:"required"`
	Method models.PaymentMethod `json:"method" validate:"required,oneof=credit_card debit_card bank_transfer e_wallet"`
}

type PaymentResponse struct {
	ID     uuid.UUID            `json:"id"`
	PNRID  *uuid.UUID           `json:"pnr_id"`
	Amount float64              `json:"amount"`
	Method models.PaymentMethod `json:"method"`
	Status models.PaymentStatus `json:"status"`
	PaidAt *time.Time           `json:"paid_at"`
}

func ToPaymentResponse(p *models.Payment) *PaymentResponse {
	if p == nil {
		return nil
	}
	return &PaymentResponse{
		ID:     p.ID,
		PNRID:  p.PNRID,
		Amount: p.Amount,
		Method: p.Method,
		Status: p.Status,
		PaidAt: p.PaidAt,
	}
}

func ToPaymentResponseList(payments []models.Payment) []PaymentResponse {
	out := make([]PaymentResponse, 0, len(payments))
	for i := range payments {
		if r := ToPaymentResponse(&payments[i]); r != nil {
			out = append(out, *r)
		}
	}
	return out
}

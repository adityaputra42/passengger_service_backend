package dto

import (
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type AddBaggageRequest struct {
	PassengerID uuid.UUID `json:"passenger_id" validate:"required"`
	SegmentID   uuid.UUID `json:"segment_id"   validate:"required"`
	Weight      float64   `json:"weight"       validate:"required,gt=0"`
}

type BaggageResponse struct {
	ID          uuid.UUID            `json:"id"`
	PassengerID *uuid.UUID           `json:"passenger_id"`
	SegmentID   *uuid.UUID           `json:"segment_id"`
	Weight      float64              `json:"weight"`
	TagNumber   string               `json:"tag_number"`
	Status      models.BaggageStatus `json:"status"`
}

func ToBaggageResponse(b *models.Baggage) *BaggageResponse {
	if b == nil {
		return nil
	}
	return &BaggageResponse{
		ID:          b.ID,
		PassengerID: b.PassengerID,
		SegmentID:   b.SegmentID,
		Weight:      b.Weight,
		TagNumber:   b.TagNumber,
		Status:      b.Status,
	}
}

func ToBaggageResponseList(bags []models.Baggage) []BaggageResponse {
	out := make([]BaggageResponse, 0, len(bags))
	for i := range bags {
		if r := ToBaggageResponse(&bags[i]); r != nil {
			out = append(out, *r)
		}
	}
	return out
}

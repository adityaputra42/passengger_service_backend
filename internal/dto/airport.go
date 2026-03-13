package dto

import (
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type CreateAirportRequest struct {
	Code     string `json:"code"     validate:"required,len=3"`
	Name     string `json:"name"     validate:"required,max=255"`
	City     string `json:"city"     validate:"required,max=255"`
	Country  string `json:"country"  validate:"required,max=255"`
	Timezone string `json:"timezone" validate:"required,max=100"`
}

type UpdateAirportRequest struct {
	Name     string `json:"name"     validate:"omitempty,max=255"`
	City     string `json:"city"     validate:"omitempty,max=255"`
	Country  string `json:"country"  validate:"omitempty,max=255"`
	Timezone string `json:"timezone" validate:"omitempty,max=100"`
}

type AirportResponse struct {
	ID       uuid.UUID `json:"id"`
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	City     string    `json:"city"`
	Country  string    `json:"country"`
	Timezone string    `json:"timezone"`
}

func ToAirportResponse(a *models.Airport) *AirportResponse {
	if a == nil {
		return nil
	}
	return &AirportResponse{
		ID:       a.ID,
		Code:     a.Code,
		Name:     a.Name,
		City:     a.City,
		Country:  a.Country,
		Timezone: a.Timezone,
	}
}

func ToAirportResponseList(airports []models.Airport) []AirportResponse {
	out := make([]AirportResponse, 0, len(airports))
	for i := range airports {
		if r := ToAirportResponse(&airports[i]); r != nil {
			out = append(out, *r)
		}
	}
	return out
}

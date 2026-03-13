package dto

import (
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)


type CreateFlightScheduleRequest struct {
	FlightNumber         string `json:"flight_number"          validate:"required,max=10"`
	DepartureAirportCode string `json:"departure_airport_code" validate:"required,len=3"`
	ArrivalAirportCode   string `json:"arrival_airport_code"   validate:"required,len=3"`
	DepartureTime        string `json:"departure_time"         validate:"required"`
	ArrivalTime          string `json:"arrival_time"           validate:"required"`
	OperatingDays        string `json:"operating_days"         validate:"required"`
}

type UpdateFlightScheduleRequest struct {
	DepartureTime string `json:"departure_time" validate:"omitempty"`
	ArrivalTime   string `json:"arrival_time"   validate:"omitempty"`
	OperatingDays string `json:"operating_days" validate:"omitempty"`
}

type FlightScheduleResponse struct {
	ID               uuid.UUID        `json:"id"`
	FlightNumber     string           `json:"flight_number"`
	DepartureAirport *AirportResponse `json:"departure_airport"`
	ArrivalAirport   *AirportResponse `json:"arrival_airport"`
	DepartureTime    string           `json:"departure_time"`
	ArrivalTime      string           `json:"arrival_time"`
	OperatingDays    string           `json:"operating_days"`
}

func ToFlightScheduleResponse(s *models.FlightSchedule) *FlightScheduleResponse {
	if s == nil {
		return nil
	}
	return &FlightScheduleResponse{
		ID:               s.ID,
		FlightNumber:     s.FlightNumber,
		DepartureAirport: ToAirportResponse(&s.DepartureAirport),
		ArrivalAirport:   ToAirportResponse(&s.ArrivalAirport),
		DepartureTime:    s.DepartureTime,
		ArrivalTime:      s.ArrivalTime,
		OperatingDays:    s.OperatingDays,
	}
}

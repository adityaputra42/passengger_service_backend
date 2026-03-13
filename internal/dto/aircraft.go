package dto

import (
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type CreateAircraftRequest struct {
	Model        string `json:"model"        validate:"required,max=100"`
	Manufacturer string `json:"manufacturer" validate:"required,max=100"`
}

type UpdateAircraftRequest struct {
	Model        string `json:"model"        validate:"omitempty,max=100"`
	Manufacturer string `json:"manufacturer" validate:"omitempty,max=100"`
}

type AircraftSeatResponse struct {
	ID         uuid.UUID          `json:"id"`
	SeatNumber string             `json:"seat_number"`
	RowNumber  int                `json:"row_number"`
	SeatLetter string             `json:"seat_letter"`
	XPosition  int                `json:"x_position"`
	YPosition  int                `json:"y_position"`
	SeatType   string             `json:"seat_type"`
	IsExitRow  bool               `json:"is_exit_row"`
	SeatClass  *SeatClassResponse `json:"seat_class"`
}

func ToAircraftSeatResponse(s *models.AircraftSeat) *AircraftSeatResponse {
	if s == nil {
		return nil
	}
	return &AircraftSeatResponse{
		ID:         s.ID,
		SeatNumber: s.SeatNumber,
		RowNumber:  s.RowNumber,
		SeatLetter: s.SeatLetter,
		XPosition:  s.XPosition,
		YPosition:  s.YPosition,
		SeatType:   s.SeatType,
		IsExitRow:  s.IsExitRow,
		SeatClass:  ToSeatClassResponse(s.SeatClass),
	}
}

type AircraftResponse struct {
	ID           uuid.UUID              `json:"id"`
	Model        string                 `json:"model"`
	Manufacturer string                 `json:"manufacturer"`
	TotalSeats   int                    `json:"total_seats"`
	Seats        []AircraftSeatResponse `json:"seats,omitempty"`
}

func ToAircraftResponse(a *models.Aircraft) *AircraftResponse {
	if a == nil {
		return nil
	}
	seats := make([]AircraftSeatResponse, 0, len(a.Seats))
	for i := range a.Seats {
		if r := ToAircraftSeatResponse(&a.Seats[i]); r != nil {
			seats = append(seats, *r)
		}
	}
	return &AircraftResponse{
		ID:           a.ID,
		Model:        a.Model,
		Manufacturer: a.Manufacturer,
		TotalSeats:   a.TotalSeats,
		Seats:        seats,
	}
}

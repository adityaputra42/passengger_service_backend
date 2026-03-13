package dto

import (
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type SeatClassConfig struct {
	ClassCode   string   `json:"class_code"   validate:"required,oneof=F C Y"`
	Rows        int      `json:"rows"         validate:"required,min=1"`
	Letters     []string `json:"letters"      validate:"required,min=1"`
	ExitRowNums []int    `json:"exit_row_nums"`
}

type GenerateSeatsRequest struct {
	Classes []SeatClassConfig `json:"classes" validate:"required,min=1"`
}

type SeatClassResponse struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

func ToSeatClassResponse(sc *models.SeatClass) *SeatClassResponse {
	if sc == nil {
		return nil
	}
	return &SeatClassResponse{ID: sc.ID, Code: sc.Code, Name: sc.Name}
}

package dto

import (
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)
type CheckinRequest struct {
	PassengerID uuid.UUID `json:"passenger_id" validate:"required"`
	SegmentID   uuid.UUID `json:"segment_id"   validate:"required"`
}

type CheckinResult struct {
	Checkin      *models.Checkin      `json:"checkin"`
	BoardingPass *models.BoardingPass `json:"boarding_pass,omitempty"`
}

type CheckinResponse struct {
	ID          uuid.UUID  `json:"id"`
	PassengerID *uuid.UUID `json:"passenger_id"`
	SegmentID   *uuid.UUID `json:"segment_id"`
	CheckinTime *time.Time `json:"checkin_time"`
}

func ToCheckinResponse(c *models.Checkin) *CheckinResponse {
	if c == nil {
		return nil
	}
	return &CheckinResponse{
		ID:          c.ID,
		PassengerID: c.PassengerID,
		SegmentID:   c.SegmentID,
		CheckinTime: c.CheckinTime,
	}
}

type CheckinResultResponse struct {
	Checkin      *CheckinResponse      `json:"checkin"`
	BoardingPass *BoardingPassResponse `json:"boarding_pass"`
}

func ToCheckinResultResponse(r *CheckinResult) *CheckinResultResponse {
	if r == nil {
		return nil
	}
	return &CheckinResultResponse{
		Checkin:      ToCheckinResponse(r.Checkin),
		BoardingPass: ToBoardingPassResponse(r.BoardingPass),
	}
}

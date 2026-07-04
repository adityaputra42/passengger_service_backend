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

type CheckinGetRequest struct {
	PassengerID uuid.UUID `json:"passenger_id"`
	SegmentID   uuid.UUID `json:"segment_id"`
}

type CheckinResult struct {
	Checkin      *models.Checkin      `json:"checkin"`
	BoardingPass *models.BoardingPass `json:"boarding_pass,omitempty"`
}

type CheckinResponse struct {
	ID          uuid.UUID            `json:"id"`
	PassengerID *models.PNRPassenger `json:"passenger"`
	SegmentID   *models.PNRSegment   `json:"segment_id"`
	CheckinTime *time.Time           `json:"checkin_time"`
}

func ToCheckinResponse(c *models.Checkin) *CheckinResponse {
	if c == nil {
		return nil
	}
	return &CheckinResponse{
		ID:          c.ID,
		PassengerID: c.Passenger,
		SegmentID:   c.Segment,
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

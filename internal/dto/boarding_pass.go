package dto

import (
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)


type IssueBoardingPassRequest struct {
	PassengerID   uuid.UUID `json:"passenger_id"   validate:"required"`
	SegmentID     uuid.UUID `json:"segment_id"     validate:"required"`
	Gate          string    `json:"gate"           validate:"required,max=10"`
	BoardingGroup string    `json:"boarding_group" validate:"required,max=10"`
	BoardingTime  time.Time `json:"boarding_time"  validate:"required"`
}

type BoardingPassResponse struct {
	ID            uuid.UUID  `json:"id"`
	PassengerID   *uuid.UUID `json:"passenger_id"`
	PassengerName string     `json:"passenger_name"`
	SegmentID     *uuid.UUID `json:"segment_id"`
	FlightNumber  string     `json:"flight_number"`
	Origin        string     `json:"origin"`
	Destination   string     `json:"destination"`
	DepartureTime *time.Time `json:"departure_time"`
	SeatNumber    string     `json:"seat_number"`
	BoardingGroup string     `json:"boarding_group"`
	Gate          string     `json:"gate"`
	BoardingTime  *time.Time `json:"boarding_time"`
	QRCode        string     `json:"qr_code"`
}

func ToBoardingPassResponse(bp *models.BoardingPass) *BoardingPassResponse {
	if bp == nil {
		return nil
	}

	r := &BoardingPassResponse{
		ID:            bp.ID,
		PassengerID:   bp.PassengerID,
		SegmentID:     bp.SegmentID,
		BoardingGroup: bp.BoardingGroup,
		Gate:          bp.Gate,
		BoardingTime:  bp.BoardingTime,
		QRCode:        bp.QRCode,
	}

	if bp.Passenger != nil {
		r.PassengerName = bp.Passenger.FirstName + " " + bp.Passenger.LastName
	}

	if bp.Segment != nil && bp.Segment.Flight != nil {
		f := bp.Segment.Flight
		r.DepartureTime = f.DepartureTime
		if f.Schedule != nil {
			r.FlightNumber = f.Schedule.FlightNumber
			r.Origin = f.Schedule.DepartureAirport.Code
			r.Destination = f.Schedule.ArrivalAirport.Code
		}
	}

	return r
}

func ToBoardingPassResponseList(bps []models.BoardingPass) []BoardingPassResponse {
	out := make([]BoardingPassResponse, 0, len(bps))
	for i := range bps {
		if r := ToBoardingPassResponse(&bps[i]); r != nil {
			out = append(out, *r)
		}
	}
	return out
}

package dto

import (
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)


type SearchFlightRequest struct {
	DepartureCode string    `json:"departure_code" validate:"required,len=3"`
	ArrivalCode   string    `json:"arrival_code"   validate:"required,len=3"`
	Date          time.Time `json:"date"           validate:"required"`
}

type FlightResult struct {
	Flight         models.Flight `json:"flight"`
	AvailableSeats int64         `json:"available_seats"`
	LowestPrice    float64       `json:"lowest_price"`
}

type FlightSeatResult struct {
	FlightSeat models.FlightSeat `json:"flight_seat"`
	SeatNumber string            `json:"seat_number"`
	RowNumber  int               `json:"row_number"`
	SeatLetter string            `json:"seat_letter"`
	SeatType   string            `json:"seat_type"`
	IsExitRow  bool              `json:"is_exit_row"`
	ClassCode  string            `json:"class_code"`
	ClassName  string            `json:"class_name"`
}

// ─────────────────────────────────────────────
// FlightSeat
// ─────────────────────────────────────────────

type FlightSeatResponse struct {
	ID           uuid.UUID             `json:"id"`
	FlightID     *uuid.UUID            `json:"flight_id"`
	Price        float64               `json:"price"`
	Status       models.FlightSeatStatus `json:"status"`
	AircraftSeat *AircraftSeatResponse `json:"aircraft_seat"`
}

func ToFlightSeatResponse(fs *models.FlightSeat) *FlightSeatResponse {
	if fs == nil {
		return nil
	}
	return &FlightSeatResponse{
		ID:           fs.ID,
		FlightID:     fs.FlightID,
		Price:        fs.Price,
		Status:       fs.Status,
		AircraftSeat: ToAircraftSeatResponse(fs.AircraftSeat),
	}
}

// ─────────────────────────────────────────────
// Flight
// ─────────────────────────────────────────────

type FlightResponse struct {
	ID             uuid.UUID               `json:"id"`
	FlightNumber   string                  `json:"flight_number"`
	Schedule       *FlightScheduleResponse `json:"schedule"`
	Aircraft       *AircraftResponse       `json:"aircraft"`
	DepartureTime  *time.Time              `json:"departure_time"`
	ArrivalTime    *time.Time              `json:"arrival_time"`
	DurationMin    int                     `json:"duration_minutes"`
	Status         models.FlightStatus     `json:"status"`
	AvailableSeats int64                   `json:"available_seats,omitempty"`
	LowestPrice    float64                 `json:"lowest_price,omitempty"`
}

type FlightListResponse struct {
	Flights []FlightResponse `json:"flights"`
}

func ToFlightResponse(f *models.Flight) *FlightResponse {
	if f == nil {
		return nil
	}
	flightNumber := ""
	if f.Schedule != nil {
		flightNumber = f.Schedule.FlightNumber
	}

	durationMin := 0
	if f.DepartureTime != nil && f.ArrivalTime != nil {
		durationMin = int(f.ArrivalTime.Sub(*f.DepartureTime).Minutes())
	}

	return &FlightResponse{
		ID:            f.ID,
		FlightNumber:  flightNumber,
		Schedule:      ToFlightScheduleResponse(f.Schedule),
		Aircraft:      ToAircraftResponse(f.Aircraft),
		DepartureTime: f.DepartureTime,
		ArrivalTime:   f.ArrivalTime,
		DurationMin:   durationMin,
		Status:        f.Status,
	}
}

// ToFlightSearchResponse enriches with availability data from service layer.
func ToFlightSearchResponse(r FlightResult) *FlightResponse {
	resp := ToFlightResponse(&r.Flight)
	if resp == nil {
		return nil
	}
	resp.AvailableSeats = r.AvailableSeats
	resp.LowestPrice = r.LowestPrice
	return resp
}

func ToFlightSearchResponseList(results []FlightResult) *FlightListResponse {
	out := make([]FlightResponse, 0, len(results))
	for _, r := range results {
		if resp := ToFlightSearchResponse(r); resp != nil {
			out = append(out, *resp)
		}
	}
	return &FlightListResponse{Flights: out}
}

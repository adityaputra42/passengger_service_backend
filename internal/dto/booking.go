package dto

import (
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)



type TripType string

const (
	TripOneWay    TripType = "one_way"
	TripRoundTrip TripType = "round_trip"
	TripMultiCity TripType = "multi_city"
)

type SegmentRequest struct {
	FlightID       uuid.UUID              `json:"flight_id"       validate:"required"`
	SeatSelections []SeatSelectionRequest `json:"seat_selections"`
}

type SeatSelectionRequest struct {
	PassengerIndex int       `json:"passenger_index" validate:"min=0"`
	FlightSeatID   uuid.UUID `json:"flight_seat_id"  validate:"required"`
}

type PassengerRequest struct {
	FirstName      string               `json:"first_name"      validate:"required,max=100"`
	LastName       string               `json:"last_name"       validate:"required,max=100"`
	PassengerType  models.PassengerType `json:"passenger_type"  validate:"required,oneof=ADT CHD INF"`
	BirthDate      *time.Time           `json:"birth_date"      validate:"omitempty"`
	PassportNumber string               `json:"passport_number" validate:"omitempty,max=50"`
}

type ContactRequest struct {
	Name  string `json:"name"  validate:"required,max=255"`
	Email string `json:"email" validate:"required,email,max=255"`
	Phone string `json:"phone" validate:"required,max=50"`
}

// CreateBookingRequest adalah payload utama booking.
//
// ── One Way ──────────────────────────────────
//
//	{
//	  "trip_type": "one_way",
//	  "contact":   { "name": "Budi", "email": "budi@mail.com", "phone": "0811" },
//	  "passengers": [
//	    { "first_name": "Budi", "last_name": "Santoso", "passenger_type": "ADT" }
//	  ],
//	  "segments": [
//	    {
//	      "flight_id": "uuid-GA401",
//	      "seat_selections": [
//	        { "passenger_index": 0, "flight_seat_id": "uuid-12A" }
//	      ]
//	    }
//	  ]
//	}
//
// ── Round Trip ───────────────────────────────
//
//	{
//	  "trip_type": "round_trip",
//	  "contact":   { ... },
//	  "passengers": [
//	    { "first_name": "Budi",  "last_name": "S", "passenger_type": "ADT" },
//	    { "first_name": "Siti",  "last_name": "R", "passenger_type": "ADT" }
//	  ],
//	  "segments": [
//	    {
//	      "flight_id": "uuid-GA401-outbound",
//	      "seat_selections": [
//	        { "passenger_index": 0, "flight_seat_id": "uuid-12A" },
//	        { "passenger_index": 1, "flight_seat_id": "uuid-12B" }
//	      ]
//	    },
//	    {
//	      "flight_id": "uuid-GA402-return",
//	      "seat_selections": [
//	        { "passenger_index": 0, "flight_seat_id": "uuid-14C" },
//	        { "passenger_index": 1, "flight_seat_id": "uuid-14D" }
//	      ]
//	    }
//	  ]
//	}
//
// ── Multi City ───────────────────────────────
//
//	{
//	  "trip_type": "multi_city",
//	  "segments": [
//	    { "flight_id": "uuid-CGK-DPS", "seat_selections": [...] },
//	    { "flight_id": "uuid-DPS-SUB", "seat_selections": [...] },
//	    { "flight_id": "uuid-SUB-CGK", "seat_selections": [...] }
//	  ],
//	  ...
//	}
type CreateBookingRequest struct {
	TripType   TripType           `json:"trip_type"  validate:"required,oneof=one_way round_trip multi_city"`
	Contact    ContactRequest     `json:"contact"    validate:"required"`
	Passengers []PassengerRequest `json:"passengers" validate:"required,min=1,max=9"`
	Segments   []SegmentRequest   `json:"segments"   validate:"required,min=1,max=8"`
}



type UpdateContactRequest struct {
	Name  string `json:"name"  validate:"omitempty,max=255"`
	Email string `json:"email" validate:"omitempty,email,max=255"`
	Phone string `json:"phone" validate:"omitempty,max=50"`
}

type AddSSRRequest struct {
	SSRCode   string    `json:"ssr_code"   validate:"required,max=10"`
	SegmentID uuid.UUID `json:"segment_id" validate:"required"`
}

type AddMealRequest struct {
	MealCode  string    `json:"meal_code"  validate:"required,max=10"`
	SegmentID uuid.UUID `json:"segment_id" validate:"required"`
}




type SeatAssignmentResponse struct {
	ID           uuid.UUID           `json:"id"`
	SegmentID    *uuid.UUID          `json:"segment_id"`
	FlightSeatID *uuid.UUID          `json:"flight_seat_id"`
	SeatNumber   string              `json:"seat_number"`
	SeatClass    *SeatClassResponse  `json:"seat_class"`
	AssignedAt   *time.Time          `json:"assigned_at"`
}

func ToSeatAssignmentResponse(sa *models.SeatAssignment) *SeatAssignmentResponse {
	if sa == nil {
		return nil
	}
	seatNumber := ""
	var seatClass *SeatClassResponse
	if sa.FlightSeat != nil && sa.FlightSeat.AircraftSeat != nil {
		seatNumber = sa.FlightSeat.AircraftSeat.SeatNumber
		seatClass = ToSeatClassResponse(sa.FlightSeat.AircraftSeat.SeatClass)
	}
	return &SeatAssignmentResponse{
		ID:           sa.ID,
		SegmentID:    sa.SegmentID,
		FlightSeatID: sa.FlightSeatID,
		SeatNumber:   seatNumber,
		SeatClass:    seatClass,
		AssignedAt:   sa.AssignedAt,
	}
}

// ─────────────────────────────────────────────
// SSR & Meal (inline, used inside passenger)
// ─────────────────────────────────────────────

type PassengerSSRResponse struct {
	ID        uuid.UUID  `json:"id"`
	SegmentID *uuid.UUID `json:"segment_id"`
	SSRCode   string     `json:"ssr_code"`
	SSRName   string     `json:"ssr_name"`
}

func ToPassengerSSRResponse(s *models.PassengerSSR) *PassengerSSRResponse {
	if s == nil {
		return nil
	}
	code, name := "", ""
	if s.SSRType != nil {
		code = s.SSRType.Code
		name = s.SSRType.Name
	}
	return &PassengerSSRResponse{
		ID:        s.ID,
		SegmentID: s.SegmentID,
		SSRCode:   code,
		SSRName:   name,
	}
}

type PassengerMealResponse struct {
	ID        uuid.UUID  `json:"id"`
	SegmentID *uuid.UUID `json:"segment_id"`
	MealCode  string     `json:"meal_code"`
	MealName  string     `json:"meal_name"`
}

func ToPassengerMealResponse(m *models.PassengerMeal) *PassengerMealResponse {
	if m == nil {
		return nil
	}
	code, name := "", ""
	if m.Meal != nil {
		code = m.Meal.Code
		name = m.Meal.Name
	}
	return &PassengerMealResponse{
		ID:        m.ID,
		SegmentID: m.SegmentID,
		MealCode:  code,
		MealName:  name,
	}
}

// ─────────────────────────────────────────────
// Ticket (inline, used inside passenger)
// ─────────────────────────────────────────────

type TicketResponse struct {
	ID           uuid.UUID  `json:"id"`
	TicketNumber string     `json:"ticket_number"`
	IssuedAt     *time.Time `json:"issued_at"`
}

func ToTicketResponse(t *models.Ticket) *TicketResponse {
	if t == nil {
		return nil
	}
	return &TicketResponse{
		ID:           t.ID,
		TicketNumber: t.TicketNumber,
		IssuedAt:     t.IssuedAt,
	}
}



type PNRPassengerResponse struct {
	ID             uuid.UUID               `json:"id"`
	FirstName      string                  `json:"first_name"`
	LastName       string                  `json:"last_name"`
	FullName       string                  `json:"full_name"`
	PassengerType  models.PassengerType    `json:"passenger_type"`
	BirthDate      *time.Time              `json:"birth_date"`
	PassportNumber string                  `json:"passport_number"`
	Ticket         *TicketResponse         `json:"ticket"`
	SeatAssignment *SeatAssignmentResponse `json:"seat_assignment"`
	SSRs           []PassengerSSRResponse  `json:"ssrs"`
	Meals          []PassengerMealResponse `json:"meals"`
	Baggage        []BaggageResponse       `json:"baggage"`
}

func ToPNRPassengerResponse(p *models.PNRPassenger) *PNRPassengerResponse {
	if p == nil {
		return nil
	}

	ssrs := make([]PassengerSSRResponse, 0, len(p.SSRs))
	for i := range p.SSRs {
		if r := ToPassengerSSRResponse(&p.SSRs[i]); r != nil {
			ssrs = append(ssrs, *r)
		}
	}

	meals := make([]PassengerMealResponse, 0, len(p.Meals))
	for i := range p.Meals {
		if r := ToPassengerMealResponse(&p.Meals[i]); r != nil {
			meals = append(meals, *r)
		}
	}

	return &PNRPassengerResponse{
		ID:             p.ID,
		FirstName:      p.FirstName,
		LastName:       p.LastName,
		FullName:       p.FirstName + " " + p.LastName,
		PassengerType:  p.PassengerType,
		BirthDate:      p.BirthDate,
		PassportNumber: p.PassportNumber,
		Ticket:         ToTicketResponse(p.Ticket),
		SeatAssignment: ToSeatAssignmentResponse(p.SeatAssignment),
		SSRs:           ssrs,
		Meals:          meals,
		Baggage:        ToBaggageResponseList(p.Baggage),
	}
}

// ─────────────────────────────────────────────
// PNRSegment
// ─────────────────────────────────────────────

type PNRSegmentResponse struct {
	ID           uuid.UUID       `json:"id"`
	SegmentOrder int             `json:"segment_order"`
	Flight       *FlightResponse `json:"flight"`
}

func ToPNRSegmentResponse(s *models.PNRSegment) *PNRSegmentResponse {
	if s == nil {
		return nil
	}
	return &PNRSegmentResponse{
		ID:           s.ID,
		SegmentOrder: s.SegmentOrder,
		Flight:       ToFlightResponse(s.Flight),
	}
}



// ─────────────────────────────────────────────
// PNR — top-level booking response
// ─────────────────────────────────────────────

type PNRContactResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Phone string    `json:"phone"`
}

func ToPNRContactResponse(c *models.PNRContact) *PNRContactResponse {
	if c == nil {
		return nil
	}
	return &PNRContactResponse{
		ID:    c.ID,
		Name:  c.Name,
		Email: c.Email,
		Phone: c.Phone,
	}
}

// PNRResponse adalah response utama setelah booking — dipakai di semua endpoint
// yang return detail PNR: CreateBooking, GetPNR, GetPNRByID.
type PNRResponse struct {
	ID            uuid.UUID              `json:"id"`
	RecordLocator string                 `json:"record_locator"`
	Status        models.PNRStatus       `json:"status"`
	TripType      string                 `json:"trip_type"`
	TTL           *time.Time             `json:"ttl"`
	CreatedAt     time.Time              `json:"created_at"`
	Contact       *PNRContactResponse    `json:"contact"`
	Passengers    []PNRPassengerResponse `json:"passengers"`
	Segments      []PNRSegmentResponse   `json:"segments"`
	Payments      []PaymentResponse      `json:"payments"`
	TotalAmount   float64                `json:"total_amount"`
}

type PNRListResponse struct {
	PNRs []PNRResponse `json:"pnrs"`
	PaginationMeta
}

func ToPNRResponse(p *models.PNR) *PNRResponse {
	if p == nil {
		return nil
	}

	passengers := make([]PNRPassengerResponse, 0, len(p.Passengers))
	for i := range p.Passengers {
		if r := ToPNRPassengerResponse(&p.Passengers[i]); r != nil {
			passengers = append(passengers, *r)
		}
	}

	segments := make([]PNRSegmentResponse, 0, len(p.Segments))
	for i := range p.Segments {
		if r := ToPNRSegmentResponse(&p.Segments[i]); r != nil {
			segments = append(segments, *r)
		}
	}

	payments := ToPaymentResponseList(p.Payments)

	// trip_type di-infer dari jumlah segment dan airport pattern
	tripType := inferTripType(p.Segments)

	// total_amount dari semua payment yang success
	var totalAmount float64
	for _, pay := range p.Payments {
		if pay.Status == models.PaymentStatusSuccess {
			totalAmount += pay.Amount
		}
	}

	return &PNRResponse{
		ID:            p.ID,
		RecordLocator: p.RecordLocator,
		Status:        p.Status,
		TripType:      tripType,
		TTL:           p.TTL,
		CreatedAt:     p.CreatedAt,
		Contact:       ToPNRContactResponse(p.Contact),
		Passengers:    passengers,
		Segments:      segments,
		Payments:      payments,
		TotalAmount:   totalAmount,
	}
}

func ToPNRListResponse(pnrs []models.PNR, total int64, page, limit int) *PNRListResponse {
	out := make([]PNRResponse, 0, len(pnrs))
	for i := range pnrs {
		if r := ToPNRResponse(&pnrs[i]); r != nil {
			out = append(out, *r)
		}
	}
	return &PNRListResponse{
		PNRs:           out,
		PaginationMeta: newPagination(total, page, limit),
	}
}

// inferTripType infers trip type from segment count and airport routing.
// Dipanggil di ToPNRResponse karena PNR model tidak menyimpan trip_type
// secara eksplisit — trip_type adalah konsep booking request, bukan storage.
func inferTripType(segments []models.PNRSegment) string {
	n := len(segments)
	switch {
	case n <= 1:
		return string(TripOneWay)
	case n == 2:
		// round trip: seg[1].dep_airport == seg[0].arr_airport
		// AND seg[1].arr_airport == seg[0].dep_airport
		if len(segments) == 2 &&
			segments[0].Flight != nil && segments[0].Flight.Schedule != nil &&
			segments[1].Flight != nil && segments[1].Flight.Schedule != nil {
			s0 := segments[0].Flight.Schedule
			s1 := segments[1].Flight.Schedule
			if s0.ArrivalAirportID == s1.DepartureAirportID &&
				s0.DepartureAirportID == s1.ArrivalAirportID {
				return string(TripRoundTrip)
			}
		}
		return string(TripMultiCity)
	default:
		return string(TripMultiCity)
	}
}

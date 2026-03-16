package services

// import (
// 	"context"
// 	"passenger_service_backend/internal/models"
// 	"time"

// 	"github.com/google/uuid"
// )

// // ─────────────────────────────────────────────
// // Aircraft
// // ─────────────────────────────────────────────

// // ─────────────────────────────────────────────
// // Flight Schedule
// // ─────────────────────────────────────────────

// // ─────────────────────────────────────────────
// // Flight
// // ─────────────────────────────────────────────

// // ─────────────────────────────────────────────
// // Booking (PNR)
// // ─────────────────────────────────────────────

// type BookingService interface {
// 	// Full booking flow: lock seat → create PNR → add passengers → return PNR
// 	CreateBooking(ctx context.Context, req CreateBookingRequest) (*models.PNR, error)
// 	// Retrieve PNR with all relations
// 	GetPNR(ctx context.Context, locator string) (*models.PNR, error)
// 	GetPNRByID(ctx context.Context, id uuid.UUID) (*models.PNR, error)
// 	GetAllPNRs(ctx context.Context, page, limit int) ([]models.PNR, int64, error)
// 	// Update contact info
// 	UpdateContact(ctx context.Context, pnrID uuid.UUID, req UpdateContactRequest) (*models.PNRContact, error)
// 	// Passenger ancillary
// 	AddSSR(ctx context.Context, passengerID uuid.UUID, req AddSSRRequest) (*models.PassengerSSR, error)
// 	RemoveSSR(ctx context.Context, ssrID uuid.UUID) error
// 	AddMeal(ctx context.Context, passengerID uuid.UUID, req AddMealRequest) (*models.PassengerMeal, error)
// 	RemoveMeal(ctx context.Context, mealID uuid.UUID) error
// 	// Cancel PNR — releases seats, voids tickets if any
// 	CancelPNR(ctx context.Context, pnrID uuid.UUID) error
// }

// // ─────────────────────────────────────────────
// // Seat Lock
// // ─────────────────────────────────────────────

// type SeatLockService interface {
// 	// Lock a flight seat for a PNR, returns error if already locked/booked
// 	Lock(ctx context.Context, flightSeatID, pnrID uuid.UUID, ttl time.Duration) (*models.SeatLock, error)
// 	// Release a specific lock
// 	Release(ctx context.Context, lockID uuid.UUID) error
// 	// Background job: release all expired locks and reset seat status to available
// 	ReleaseExpired(ctx context.Context) (int, error)
// }

// // ─────────────────────────────────────────────
// // Payment
// // ─────────────────────────────────────────────

// type PaymentService interface {
// 	// Initiate a payment for a PNR
// 	Initiate(ctx context.Context, req InitiatePaymentRequest) (*models.Payment, error)
// 	// Handle payment gateway callback — confirm or fail
// 	Confirm(ctx context.Context, paymentID uuid.UUID, success bool) (*models.Payment, error)
// 	// Refund a successful payment
// 	Refund(ctx context.Context, paymentID uuid.UUID) (*models.Payment, error)
// 	GetByPNR(ctx context.Context, pnrID uuid.UUID) ([]models.Payment, error)
// 	GetByStatus(ctx context.Context, status models.PaymentStatus, page, limit int) ([]models.Payment, int64, error)
// }

// // ─────────────────────────────────────────────
// // Ticket
// // ─────────────────────────────────────────────

// type TicketService interface {
// 	// Issue tickets after successful payment
// 	IssueForPNR(ctx context.Context, pnrID uuid.UUID) ([]models.Ticket, error)
// 	GetByTicketNumber(ctx context.Context, number string) (*models.Ticket, error)
// 	GetByPassenger(ctx context.Context, passengerID uuid.UUID) (*models.Ticket, error)
// }

// // ─────────────────────────────────────────────
// // Check-in
// // ─────────────────────────────────────────────

// type CheckinService interface {
// 	// Perform check-in for a passenger on a segment
// 	Checkin(ctx context.Context, req CheckinRequest) (*CheckinResult, error)
// 	GetByPassenger(ctx context.Context, passengerID uuid.UUID) ([]models.Checkin, error)
// 	IsCheckedIn(ctx context.Context, passengerID, segmentID uuid.UUID) (bool, error)
// }

// // ─────────────────────────────────────────────
// // Baggage
// // ─────────────────────────────────────────────

// type BaggageService interface {
// 	Add(ctx context.Context, req AddBaggageRequest) (*models.Baggage, error)
// 	UpdateStatus(ctx context.Context, baggageID uuid.UUID, status models.BaggageStatus) (*models.Baggage, error)
// 	GetByPassenger(ctx context.Context, passengerID uuid.UUID) ([]models.Baggage, error)
// 	GetBySegment(ctx context.Context, segmentID uuid.UUID) ([]models.Baggage, error)
// }

// // ─────────────────────────────────────────────
// // Boarding Pass
// // ─────────────────────────────────────────────

// type BoardingPassService interface {
// 	// Issue boarding pass after check-in
// 	Issue(ctx context.Context, req IssueBoardingPassRequest) (*models.BoardingPass, error)
// 	GetByPassengerAndSegment(ctx context.Context, passengerID, segmentID uuid.UUID) (*models.BoardingPass, error)
// 	GetBySegment(ctx context.Context, segmentID uuid.UUID) ([]models.BoardingPass, error)
// }

package services

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/utils"
	"time"

	"github.com/google/uuid"
)

type CheckinService interface {
	Checkin(ctx context.Context, req dto.CheckinRequest) (*dto.CheckinResult, error)
	GetByPassenger(ctx context.Context, passengerID uuid.UUID) ([]models.Checkin, error)
	IsCheckedIn(ctx context.Context, passengerID, segmentID uuid.UUID) (bool, error)
}

type checkinService struct {
	checkinRepo    repository.CheckinRepository
	passengerRepo  repository.PNRPassengerRepository
	segmentRepo    repository.PNRSegmentRepository
	ticketRepo     repository.TicketRepository
	flightRepo     repository.FlightRepository
	seatAssignRepo repository.SeatAssignmentRepository
}

func NewCheckinService(
	checkinRepo repository.CheckinRepository,
	passengerRepo repository.PNRPassengerRepository,
	segmentRepo repository.PNRSegmentRepository,
	ticketRepo repository.TicketRepository,
	flightRepo repository.FlightRepository,
	seatAssignRepo repository.SeatAssignmentRepository,
) CheckinService {
	return &checkinService{
		checkinRepo:    checkinRepo,
		passengerRepo:  passengerRepo,
		segmentRepo:    segmentRepo,
		ticketRepo:     ticketRepo,
		flightRepo:     flightRepo,
		seatAssignRepo: seatAssignRepo,
	}
}

// Checkin validates the passenger's ticket and flight, then records the check-in.
// Check-in window: from 24h before departure to 45min before departure.
func (s *checkinService) Checkin(ctx context.Context, req dto.CheckinRequest) (*dto.CheckinResult, error) {
	// Validate passenger exists
	passenger, err := s.passengerRepo.FindByID(ctx, req.PassengerID)
	if err != nil {
		return nil, utils.ErrPassengerNotFound
	}

	// Validate segment and flight
	segment, err := s.segmentRepo.FindWithFlight(ctx, req.SegmentID)
	if err != nil {
		return nil, fmt.Errorf("segment tidak ditemukan")
	}

	if segment.Flight == nil {
		return nil, utils.ErrFlightNotFound
	}
	flight := segment.Flight

	// Flight status check
	if flight.Status == models.FlightStatusDeparted || flight.Status == models.FlightStatusArrived {
		return nil, utils.ErrFlightAlreadyDeparted
	}
	if flight.Status == models.FlightStatusCancelled {
		return nil, fmt.Errorf("penerbangan dibatalkan")
	}

	// Check-in window
	if flight.DepartureTime != nil {
		now := time.Now()
		openAt := flight.DepartureTime.Add(-24 * time.Hour)
		closeAt := flight.DepartureTime.Add(-45 * time.Minute)

		if now.Before(openAt) {
			return nil, utils.ErrCheckinTooEarly
		}
		if now.After(closeAt) {
			return nil, utils.ErrCheckinClosed
		}
	}

	// Ticket must be issued
	if _, err := s.ticketRepo.FindByPassengerID(ctx, passenger.ID); err != nil {
		return nil, utils.ErrTicketRequiredCheckin
	}

	// Prevent duplicate check-in
	alreadyCheckedIn, err := s.checkinRepo.IsCheckedIn(ctx, req.PassengerID, req.SegmentID)
	if err != nil {
		return nil, fmt.Errorf("cek status check-in: %w", err)
	}
	if alreadyCheckedIn {
		return nil, utils.ErrAlreadyCheckedIn
	}

	// Validate seat assignment exists
	if _, err := s.seatAssignRepo.FindByPassengerID(ctx, req.PassengerID); err != nil {
		return nil, fmt.Errorf("seat assignment tidak ditemukan untuk penumpang ini")
	}

	// Record check-in
	now := time.Now()
	checkin := &models.Checkin{
		PassengerID: &req.PassengerID,
		SegmentID:   &req.SegmentID,
		CheckinTime: &now,
	}
	if err := s.checkinRepo.Create(ctx, checkin); err != nil {
		return nil, fmt.Errorf("create checkin: %w", err)
	}

	return &dto.CheckinResult{Checkin: checkin}, nil
}

func (s *checkinService) GetByPassenger(ctx context.Context, passengerID uuid.UUID) ([]models.Checkin, error) {
	return s.checkinRepo.FindByPassengerID(ctx, passengerID)
}

func (s *checkinService) IsCheckedIn(ctx context.Context, passengerID, segmentID uuid.UUID) (bool, error) {
	return s.checkinRepo.IsCheckedIn(ctx, passengerID, segmentID)
}

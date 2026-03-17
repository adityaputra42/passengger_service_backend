package services

import (
	"context"
	"fmt"
	"math/rand"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/utils"
	"time"

	"github.com/google/uuid"
)

type TicketService interface {
	IssueForPNR(ctx context.Context, pnrID uuid.UUID) ([]models.Ticket, error)
	GetByTicketNumber(ctx context.Context, number string) (*models.Ticket, error)
	GetByPassenger(ctx context.Context, passengerID uuid.UUID) (*models.Ticket, error)
}

type ticketService struct {
	ticketRepo        repository.TicketRepository
	ticketSegmentRepo repository.TicketSegmentRepository
	pnrRepo           repository.PNRRepository
	passengerRepo     repository.PNRPassengerRepository
	segmentRepo       repository.PNRSegmentRepository
	flightSeatRepo    repository.FlightSeatRepository
}

func NewTicketService(
	ticketRepo repository.TicketRepository,
	ticketSegmentRepo repository.TicketSegmentRepository,
	pnrRepo repository.PNRRepository,
	passengerRepo repository.PNRPassengerRepository,
	segmentRepo repository.PNRSegmentRepository,
	flightSeatRepo repository.FlightSeatRepository,
) TicketService {
	return &ticketService{
		ticketRepo:        ticketRepo,
		ticketSegmentRepo: ticketSegmentRepo,
		pnrRepo:           pnrRepo,
		passengerRepo:     passengerRepo,
		segmentRepo:       segmentRepo,
		flightSeatRepo:    flightSeatRepo,
	}
}

// IssueForPNR issues one ticket per passenger, links all segments,
// marks flight seats as booked, and updates PNR → ticketed.
func (s *ticketService) IssueForPNR(ctx context.Context, pnrID uuid.UUID) ([]models.Ticket, error) {
	pnr, err := s.pnrRepo.FindWithFull(ctx, pnrID)
	if err != nil {
		return nil, utils.ErrPNRNotFound
	}

	if pnr.Status == models.PNRStatusCancelled {
		return nil, utils.ErrPNRAlreadyCancelled
	}
	if pnr.Status == models.PNRStatusTicketed {
		return nil, utils.ErrPNRAlreadyTicketed
	}

	segments, err := s.segmentRepo.FindByPNRID(ctx, pnrID)
	if err != nil || len(segments) == 0 {
		return nil, fmt.Errorf("no segments found for PNR")
	}

	var tickets []models.Ticket
	now := time.Now()

	for _, passenger := range pnr.Passengers {
		// Check duplicate
		if _, err := s.ticketRepo.FindByPassengerID(ctx, passenger.ID); err == nil {
			return nil, utils.ErrTicketAlreadyIssued
		}

		ticketNum := generateTicketNumber()
		ticket := &models.Ticket{
			PassengerID:  &passenger.ID,
			TicketNumber: ticketNum,
			IssuedAt:     &now,
		}
		if err := s.ticketRepo.Create(ctx, ticket); err != nil {
			return nil, fmt.Errorf("create ticket for passenger %s: %w", passenger.ID, err)
		}

		// Link each segment
		var tsSegs []models.TicketSegment
		for _, seg := range segments {
			segID := seg.ID
			tsSegs = append(tsSegs, models.TicketSegment{
				TicketID:  &ticket.ID,
				SegmentID: &segID,
			})
		}
		if err := s.ticketSegmentRepo.BulkCreate(ctx, tsSegs); err != nil {
			return nil, fmt.Errorf("create ticket segments: %w", err)
		}

		// Mark seat as booked
		if passenger.SeatAssignment != nil && passenger.SeatAssignment.FlightSeatID != nil {
			_ = s.flightSeatRepo.UpdateStatus(ctx, *passenger.SeatAssignment.FlightSeatID, models.FlightSeatBooked)
		}

		tickets = append(tickets, *ticket)
	}

	// Update PNR → ticketed
	_ = s.pnrRepo.UpdateStatus(ctx, pnrID, models.PNRStatusTicketed)

	return tickets, nil
}

func (s *ticketService) GetByTicketNumber(ctx context.Context, number string) (*models.Ticket, error) {
	t, err := s.ticketRepo.FindByTicketNumber(ctx, number)
	if err != nil {
		return nil, utils.ErrTicketNotFound
	}
	return t, nil
}

func (s *ticketService) GetByPassenger(ctx context.Context, passengerID uuid.UUID) (*models.Ticket, error) {
	t, err := s.ticketRepo.FindByPassengerID(ctx, passengerID)
	if err != nil {
		return nil, utils.ErrTicketNotFound
	}
	return t, nil
}
func generateTicketNumber() string {
	return fmt.Sprintf("126-%010d", rand.Int63n(9_999_999_999))
}

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
// marks flight seats as booked, and updates PNR → ticketed.// internal/services/ticket_service.go

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

	type result struct {
		ticket *models.Ticket
		err    error
	}

	ch := make(chan result, len(pnr.Passengers))
	now := time.Now()

	for _, passenger := range pnr.Passengers {
		go func(p models.PNRPassenger) {
			if _, err := s.ticketRepo.FindByPassengerID(ctx, p.ID); err == nil {
				ch <- result{err: utils.ErrTicketAlreadyIssued}
				return
			}

			ticketNum := generateTicketNumber()
			ticket := &models.Ticket{
				PassengerID:  &p.ID,
				TicketNumber: ticketNum,
				IssuedAt:     &now,
			}
			if err := s.ticketRepo.Create(ctx, ticket); err != nil {
				ch <- result{err: fmt.Errorf("create ticket passenger %s: %w", p.ID, err)}
				return
			}

			var tsSegs []models.TicketSegment
			for _, seg := range segments {
				segID := seg.ID
				tsSegs = append(tsSegs, models.TicketSegment{
					TicketID:  &ticket.ID,
					SegmentID: &segID,
				})
			}
			if err := s.ticketSegmentRepo.BulkCreate(ctx, tsSegs); err != nil {
				ch <- result{err: fmt.Errorf("create ticket segments: %w", err)}
				return
			}

			if p.SeatAssignment != nil && p.SeatAssignment.FlightSeatID != nil {
				_ = s.flightSeatRepo.UpdateStatus(ctx, *p.SeatAssignment.FlightSeatID, models.FlightSeatBooked)
			}

			ch <- result{ticket: ticket}
		}(passenger)
	}

	tickets := make([]models.Ticket, 0, len(pnr.Passengers))
	for range pnr.Passengers {
		r := <-ch
		if r.err != nil {
			// ticket issuance is critical — if any fail, return error
			return nil, r.err
		}
		tickets = append(tickets, *r.ticket)
	}

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

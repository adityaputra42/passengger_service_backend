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

type PaymentService interface {
	Initiate(ctx context.Context, req dto.InitiatePaymentRequest) (*models.Payment, error)
	Confirm(ctx context.Context, paymentID uuid.UUID, success bool) (*models.Payment, error)
	Refund(ctx context.Context, paymentID uuid.UUID) (*models.Payment, error)
	GetByPNR(ctx context.Context, pnrID uuid.UUID) ([]models.Payment, error)
	GetByStatus(ctx context.Context, status models.PaymentStatus, page, limit int) ([]models.Payment, int64, error)
}

type paymentService struct {
	paymentRepo    repository.PaymentRepository
	pnrRepo        repository.PNRRepository
	flightSeatRepo repository.FlightSeatRepository
	ticketSvc      TicketService
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	pnrRepo repository.PNRRepository,
	flightSeatRepo repository.FlightSeatRepository,
	ticketSvc TicketService,
) PaymentService {
	return &paymentService{
		paymentRepo:    paymentRepo,
		pnrRepo:        pnrRepo,
		flightSeatRepo: flightSeatRepo,
		ticketSvc:      ticketSvc,
	}
}

// Initiate creates a pending payment for a PNR.
// Calculates total from all assigned flight seats.
func (s *paymentService) Initiate(ctx context.Context, req dto.InitiatePaymentRequest) (*models.Payment, error) {
	pnr, err := s.pnrRepo.FindWithFull(ctx, req.PNRID)
	if err != nil {
		return nil, utils.ErrPNRNotFound
	}

	switch pnr.Status {
	case models.PNRStatusCancelled:
		return nil, utils.ErrPNRAlreadyCancelled
	case models.PNRStatusTicketed:
		return nil, utils.ErrPNRAlreadyTicketed
	}

	// Check if TTL has passed
	if pnr.TTL != nil && time.Now().After(*pnr.TTL) && pnr.Status == models.PNRStatusHold {
		return nil, utils.ErrPNRHoldExpired
	}

	// Calculate total amount from seat assignments
	var total float64
	for _, passenger := range pnr.Passengers {
		if passenger.SeatAssignment != nil && passenger.SeatAssignment.FlightSeatID != nil {
			seat, err := s.flightSeatRepo.FindByID(ctx, *passenger.SeatAssignment.FlightSeatID)
			if err == nil {
				total += seat.Price
			}
		}
	}
	if total == 0 {
		return nil, fmt.Errorf("tidak ada seat assignment ditemukan untuk menghitung total")
	}

	payment := &models.Payment{
		PNRID:  &req.PNRID,
		Amount: total,
		Method: req.Method,
		Status: models.PaymentStatusPending,
	}
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}
	return payment, nil
}

// Confirm handles payment gateway callback.
// On success: mark seats as booked, update PNR → confirmed, issue tickets.
// On failure: mark payment failed (seats remain locked, retry allowed).
func (s *paymentService) Confirm(ctx context.Context, paymentID uuid.UUID, success bool) (*models.Payment, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, utils.ErrPaymentNotFound
	}
	if payment.Status != models.PaymentStatusPending {
		return nil, utils.ErrPaymentNotPending
	}

	if !success {
		if err := s.paymentRepo.UpdateStatus(ctx, paymentID, models.PaymentStatusFailed); err != nil {
			return nil, err
		}
		payment.Status = models.PaymentStatusFailed
		return payment, nil
	}

	// Payment succeeded
	now := time.Now()
	payment.Status = models.PaymentStatusSuccess
	payment.PaidAt = &now

	if err := s.paymentRepo.UpdateStatus(ctx, paymentID, models.PaymentStatusSuccess); err != nil {
		return nil, err
	}

	if payment.PNRID != nil {
		_ = s.pnrRepo.UpdateStatus(ctx, *payment.PNRID, models.PNRStatusConfirmed)

		_, _ = s.ticketSvc.IssueForPNR(ctx, *payment.PNRID)
	}

	return payment, nil
}

// Refund processes a refund for a successful payment.
func (s *paymentService) Refund(ctx context.Context, paymentID uuid.UUID) (*models.Payment, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, utils.ErrPaymentNotFound
	}
	if payment.Status != models.PaymentStatusSuccess {
		return nil, utils.ErrPaymentNotSuccess
	}

	if err := s.paymentRepo.UpdateStatus(ctx, paymentID, models.PaymentStatusRefunded); err != nil {
		return nil, err
	}

	if payment.PNRID != nil {
		pnr, err := s.pnrRepo.FindWithFull(ctx, *payment.PNRID)
		if err == nil && pnr.Status == models.PNRStatusCancelled {
			for _, passenger := range pnr.Passengers {
				if passenger.SeatAssignment != nil && passenger.SeatAssignment.FlightSeatID != nil {
					_ = s.flightSeatRepo.UpdateStatus(ctx, *passenger.SeatAssignment.FlightSeatID, models.FlightSeatAvailable)
				}
			}
		}
	}

	payment.Status = models.PaymentStatusRefunded
	return payment, nil
}

func (s *paymentService) GetByPNR(ctx context.Context, pnrID uuid.UUID) ([]models.Payment, error) {
	return s.paymentRepo.FindByPNRID(ctx, pnrID)
}

func (s *paymentService) GetByStatus(ctx context.Context, status models.PaymentStatus, page, limit int) ([]models.Payment, int64, error) {
	return s.paymentRepo.FindByStatus(ctx, status, page, limit)
}

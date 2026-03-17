package services

import (
	"context"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type CheckinService interface {
	// Perform check-in for a passenger on a segment
	Checkin(ctx context.Context, req CheckinRequest) (*CheckinResult, error)
	GetByPassenger(ctx context.Context, passengerID uuid.UUID) ([]models.Checkin, error)
	IsCheckedIn(ctx context.Context, passengerID, segmentID uuid.UUID) (bool, error)
}

// ─────────────────────────────────────────────
// Baggage
// ─────────────────────────────────────────────

type BaggageService interface {
	Add(ctx context.Context, req AddBaggageRequest) (*models.Baggage, error)
	UpdateStatus(ctx context.Context, baggageID uuid.UUID, status models.BaggageStatus) (*models.Baggage, error)
	GetByPassenger(ctx context.Context, passengerID uuid.UUID) ([]models.Baggage, error)
	GetBySegment(ctx context.Context, segmentID uuid.UUID) ([]models.Baggage, error)
}

// ─────────────────────────────────────────────
// Boarding Pass
// ─────────────────────────────────────────────

type BoardingPassService interface {
	// Issue boarding pass after check-in
	Issue(ctx context.Context, req IssueBoardingPassRequest) (*models.BoardingPass, error)
	GetByPassengerAndSegment(ctx context.Context, passengerID, segmentID uuid.UUID) (*models.BoardingPass, error)
	GetBySegment(ctx context.Context, segmentID uuid.UUID) ([]models.BoardingPass, error)
}

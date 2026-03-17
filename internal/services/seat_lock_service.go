package services

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/utils"
	"time"

	"github.com/google/uuid"
)

type SeatLockService interface {
	Lock(ctx context.Context, flightSeatID, pnrID uuid.UUID, ttl time.Duration) (*models.SeatLock, error)
	Release(ctx context.Context, lockID uuid.UUID) error
	ReleaseExpired(ctx context.Context) (int, error)
}

type seatLockService struct {
	seatLockRepo   repository.SeatLockRepository
	flightSeatRepo repository.FlightSeatRepository
}

func NewSeatLockService(
	seatLockRepo repository.SeatLockRepository,
	flightSeatRepo repository.FlightSeatRepository,
) SeatLockService {
	return &seatLockService{
		seatLockRepo:   seatLockRepo,
		flightSeatRepo: flightSeatRepo,
	}
}

func (s *seatLockService) Lock(ctx context.Context, flightSeatID, pnrID uuid.UUID, ttl time.Duration) (*models.SeatLock, error) {
	seat, err := s.flightSeatRepo.FindByID(ctx, flightSeatID)
	if err != nil {
		return nil, utils.ErrFlightSeatNotFound
	}

	switch seat.Status {
	case models.FlightSeatBooked:
		return nil, utils.ErrSeatAlreadyBooked
	case models.FlightSeatLocked:
		// Check if existing lock belongs to same PNR (re-lock allowed)
		existingLock, err := s.seatLockRepo.FindByFlightSeatID(ctx, flightSeatID)
		if err == nil && existingLock.PNRID != nil && *existingLock.PNRID == pnrID {
			// Same PNR — extend lock
			expires := time.Now().Add(ttl)
			existingLock.ExpiresAt = &expires
			return existingLock, nil
		}
		return nil, utils.ErrSeatAlreadyLocked
	}

	// Update seat status → locked
	if err := s.flightSeatRepo.UpdateStatus(ctx, flightSeatID, models.FlightSeatLocked); err != nil {
		return nil, fmt.Errorf("lock seat: %w", err)
	}

	now := time.Now()
	expires := now.Add(ttl)
	lock := &models.SeatLock{
		FlightSeatID: &flightSeatID,
		PNRID:        &pnrID,
		LockedAt:     &now,
		ExpiresAt:    &expires,
	}
	if err := s.seatLockRepo.Create(ctx, lock); err != nil {

		_ = s.flightSeatRepo.UpdateStatus(ctx, flightSeatID, models.FlightSeatAvailable)
		return nil, fmt.Errorf("create seat lock: %w", err)
	}
	return lock, nil
}

func (s *seatLockService) Release(ctx context.Context, lockID uuid.UUID) error {
	return s.seatLockRepo.Release(ctx, lockID)
}

func (s *seatLockService) ReleaseExpired(ctx context.Context) (int, error) {
	locks, err := s.seatLockRepo.FindExpired(ctx)
	if err != nil {
		return 0, fmt.Errorf("find expired locks: %w", err)
	}

	released := 0
	for _, lock := range locks {
		if lock.FlightSeatID != nil {
			if err := s.flightSeatRepo.UpdateStatus(ctx, *lock.FlightSeatID, models.FlightSeatAvailable); err != nil {
				continue
			}
		}
		if err := s.seatLockRepo.Release(ctx, lock.ID); err != nil {
			continue
		}
		released++
	}
	return released, nil
}

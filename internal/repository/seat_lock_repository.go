package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatLockRepository interface {
	Create(ctx context.Context, lock *models.SeatLock) error
	FindByFlightSeatID(ctx context.Context, flightSeatID uuid.UUID) (*models.SeatLock, error)
	FindByPNRID(ctx context.Context, pnrID uuid.UUID) ([]models.SeatLock, error)
	FindExpired(ctx context.Context) ([]models.SeatLock, error)
	Release(ctx context.Context, id uuid.UUID) error
	ReleaseByPNR(ctx context.Context, pnrID uuid.UUID) error
}

type seatLockRepository struct {
	db *gorm.DB
}

func NewSeatLockRepository(db *gorm.DB) SeatLockRepository {
	return &seatLockRepository{db: db}
}

func (r *seatLockRepository) Create(ctx context.Context, lock *models.SeatLock) error {
	if err := r.db.WithContext(ctx).Create(lock).Error; err != nil {
		return fmt.Errorf("SeatLockRepo.Create: %w", err)
	}
	return nil
}

func (r *seatLockRepository) FindByFlightSeatID(ctx context.Context, flightSeatID uuid.UUID) (*models.SeatLock, error) {
	var lock models.SeatLock
	if err := r.db.WithContext(ctx).
		Where("flight_seat_id = ? AND expires_at > ?", flightSeatID, time.Now()).
		First(&lock).Error; err != nil {
		return nil, fmt.Errorf("SeatLockRepo.FindByFlightSeatID: %w", err)
	}
	return &lock, nil
}

func (r *seatLockRepository) FindByPNRID(ctx context.Context, pnrID uuid.UUID) ([]models.SeatLock, error) {
	var locks []models.SeatLock
	if err := r.db.WithContext(ctx).
		Where("pnr_id = ?", pnrID).
		Find(&locks).Error; err != nil {
		return nil, fmt.Errorf("SeatLockRepo.FindByPNRID: %w", err)
	}
	return locks, nil
}

// FindExpired returns all seat locks that have passed their TTL.
// Used by a background job to release stale locks.
func (r *seatLockRepository) FindExpired(ctx context.Context) ([]models.SeatLock, error) {
	var locks []models.SeatLock
	if err := r.db.WithContext(ctx).
		Preload("FlightSeat").
		Where("expires_at <= ?", time.Now()).
		Find(&locks).Error; err != nil {
		return nil, fmt.Errorf("SeatLockRepo.FindExpired: %w", err)
	}
	return locks, nil
}

func (r *seatLockRepository) Release(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.SeatLock{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("SeatLockRepo.Release: %w", err)
	}
	return nil
}

func (r *seatLockRepository) ReleaseByPNR(ctx context.Context, pnrID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("pnr_id = ?", pnrID).
		Delete(&models.SeatLock{}).Error; err != nil {
		return fmt.Errorf("SeatLockRepo.ReleaseByPNR: %w", err)
	}
	return nil
}

package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type CheckinRepository interface {
	Create(ctx context.Context, checkin *models.Checkin) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Checkin, error)
	FindByPassengerID(ctx context.Context, passengerID uuid.UUID) ([]models.Checkin, error)
	FindBySegmentID(ctx context.Context, segmentID uuid.UUID) ([]models.Checkin, error)
	IsCheckedIn(ctx context.Context, passengerID, segmentID uuid.UUID) (bool, error)
}

type checkinRepository struct {

}

func NewCheckinRepository() CheckinRepository {
	return &checkinRepository{}
}

func (r *checkinRepository) Create(ctx context.Context, c *models.Checkin) error {
	if err := db.DB.WithContext(ctx).Create(c).Error; err != nil {
		return fmt.Errorf("CheckinRepo.Create: %w", err)
	}
	return nil
}

func (r *checkinRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Checkin, error) {
	var c models.Checkin
	if err := db.DB.WithContext(ctx).
		Preload("Passenger").
		Preload("Segment").
		First(&c, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("CheckinRepo.FindByID: %w", err)
	}
	return &c, nil
}

func (r *checkinRepository) FindByPassengerID(ctx context.Context, passengerID uuid.UUID) ([]models.Checkin, error) {
	var checkins []models.Checkin
	if err := db.DB.WithContext(ctx).
		Preload("Segment").
		Where("passenger_id = ?", passengerID).
		Order("checkin_time DESC").
		Find(&checkins).Error; err != nil {
		return nil, fmt.Errorf("CheckinRepo.FindByPassengerID: %w", err)
	}
	return checkins, nil
}

func (r *checkinRepository) FindBySegmentID(ctx context.Context, segmentID uuid.UUID) ([]models.Checkin, error) {
	var checkins []models.Checkin
	if err := db.DB.WithContext(ctx).
		Preload("Passenger").
		Where("segment_id = ?", segmentID).
		Find(&checkins).Error; err != nil {
		return nil, fmt.Errorf("CheckinRepo.FindBySegmentID: %w", err)
	}
	return checkins, nil
}

func (r *checkinRepository) IsCheckedIn(ctx context.Context, passengerID, segmentID uuid.UUID) (bool, error) {
	var count int64
	if err := db.DB.WithContext(ctx).
		Model(&models.Checkin{}).
		Where("passenger_id = ? AND segment_id = ?", passengerID, segmentID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("CheckinRepo.IsCheckedIn: %w", err)
	}
	return count > 0, nil
}

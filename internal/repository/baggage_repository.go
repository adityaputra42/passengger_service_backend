package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaggageRepository interface {
	Create(ctx context.Context, baggage *models.Baggage) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Baggage, error)
	FindByPassengerID(ctx context.Context, passengerID uuid.UUID) ([]models.Baggage, error)
	FindBySegmentID(ctx context.Context, segmentID uuid.UUID) ([]models.Baggage, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.BaggageStatus) error
}

type baggageRepository struct{ db *gorm.DB }

func NewBaggageRepository(db *gorm.DB) BaggageRepository {
	return &baggageRepository{db: db}
}
func (r *baggageRepository) Create(ctx context.Context, b *models.Baggage) error {
	if err := r.db.WithContext(ctx).Create(b).Error; err != nil {
		return fmt.Errorf("BaggageRepo.Create: %w", err)
	}
	return nil
}

func (r *baggageRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Baggage, error) {
	var b models.Baggage
	if err := r.db.WithContext(ctx).First(&b, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("BaggageRepo.FindByID: %w", err)
	}
	return &b, nil
}

func (r *baggageRepository) FindByPassengerID(ctx context.Context, id uuid.UUID) ([]models.Baggage, error) {
	var bags []models.Baggage
	if err := r.db.WithContext(ctx).Where("passenger_id = ?", id).Find(&bags).Error; err != nil {
		return nil, fmt.Errorf("BaggageRepo.FindByPassengerID: %w", err)
	}
	return bags, nil
}

func (r *baggageRepository) FindBySegmentID(ctx context.Context, id uuid.UUID) ([]models.Baggage, error) {
	var bags []models.Baggage
	if err := r.db.WithContext(ctx).Where("segment_id = ?", id).Find(&bags).Error; err != nil {
		return nil, fmt.Errorf("BaggageRepo.FindBySegmentID: %w", err)
	}
	return bags, nil
}

func (r *baggageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.BaggageStatus) error {
	if err := r.db.WithContext(ctx).
		Model(&models.Baggage{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("BaggageRepo.UpdateStatus: %w", err)
	}
	return nil
}

package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PassengerSSRRepository interface {
	Create(ctx context.Context, ssr *models.PassengerSSR) error
	FindByPassengerID(ctx context.Context, passengerID uuid.UUID) ([]models.PassengerSSR, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type PassengerSSRRepositoryImpl struct {
	db *gorm.DB
}

func NewPassengerSSRRepository(db *gorm.DB) PassengerSSRRepository {
	return &PassengerSSRRepositoryImpl{db: db}
}

func (r *PassengerSSRRepositoryImpl) Create(ctx context.Context, ssr *models.PassengerSSR) error {
	if err := r.db.WithContext(ctx).Create(ssr).Error; err != nil {
		return fmt.Errorf("PassengerSSRRepo.Create: %w", err)
	}
	return nil
}

func (r *PassengerSSRRepositoryImpl) FindByPassengerID(ctx context.Context, id uuid.UUID) ([]models.PassengerSSR, error) {
	var ssrs []models.PassengerSSR
	if err := r.db.WithContext(ctx).Preload("SSRType").Where("passenger_id = ?", id).Find(&ssrs).Error; err != nil {
		return nil, fmt.Errorf("PassengerSSRRepo.FindByPassengerID: %w", err)
	}
	return ssrs, nil
}

func (r *PassengerSSRRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.PassengerSSR{}, "id = ?", id).Error
}

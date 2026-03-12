package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type PassengerSSRRepository interface {
	Create(ctx context.Context, ssr *models.PassengerSSR) error
	FindByPassengerID(ctx context.Context, passengerID uuid.UUID) ([]models.PassengerSSR, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type passengerSSRRepository struct{  }

func NewPassengerSSRRepository() PassengerSSRRepository { return &passengerSSRRepository{} }



func (r *passengerSSRRepository) Create(ctx context.Context, ssr *models.PassengerSSR) error {
	if err := db.DB.WithContext(ctx).Create(ssr).Error; err != nil {
		return fmt.Errorf("PassengerSSRRepo.Create: %w", err)
	}
	return nil
}

func (r *passengerSSRRepository) FindByPassengerID(ctx context.Context, id uuid.UUID) ([]models.PassengerSSR, error) {
	var ssrs []models.PassengerSSR
	if err := db.DB.WithContext(ctx).Preload("SSRType").Where("passenger_id = ?", id).Find(&ssrs).Error; err != nil {
		return nil, fmt.Errorf("PassengerSSRRepo.FindByPassengerID: %w", err)
	}
	return ssrs, nil
}

func (r *passengerSSRRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return db.DB.WithContext(ctx).Delete(&models.PassengerSSR{}, "id = ?", id).Error
}

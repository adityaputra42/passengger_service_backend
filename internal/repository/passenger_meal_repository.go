package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PassengerMealRepository interface {
	Create(ctx context.Context, meal *models.PassengerMeal) error
	FindByPassengerID(ctx context.Context, passengerID uuid.UUID) ([]models.PassengerMeal, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type passengerMealRepository struct {
	db *gorm.DB
}

func NewPassengerMealRepository(db *gorm.DB) PassengerMealRepository {
	return &passengerMealRepository{db: db}
}

func (r *passengerMealRepository) Create(ctx context.Context, m *models.PassengerMeal) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("PassengerMealRepo.Create: %w", err)
	}
	return nil
}

func (r *passengerMealRepository) FindByPassengerID(ctx context.Context, id uuid.UUID) ([]models.PassengerMeal, error) {
	var meals []models.PassengerMeal
	if err := r.db.WithContext(ctx).Preload("Meal").Where("passenger_id = ?", id).Find(&meals).Error; err != nil {
		return nil, fmt.Errorf("PassengerMealRepo.FindByPassengerID: %w", err)
	}
	return meals, nil
}

func (r *passengerMealRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.PassengerMeal{}, "id = ?", id).Error
}

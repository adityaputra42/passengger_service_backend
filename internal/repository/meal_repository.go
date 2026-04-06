package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"gorm.io/gorm"
)

type MealRepository interface {
	FindAll(ctx context.Context) ([]models.Meal, error)
	FindByCode(ctx context.Context, code string) (*models.Meal, error)
}

type mealRepository struct {
	db *gorm.DB
}

func NewMealRepository(db *gorm.DB) MealRepository {
	return &mealRepository{db: db}
}

func (r *mealRepository) FindAll(ctx context.Context) ([]models.Meal, error) {
	var meals []models.Meal
	if err := r.db.WithContext(ctx).Order("code").Find(&meals).Error; err != nil {
		return nil, fmt.Errorf("MealRepo.FindAll: %w", err)
	}
	return meals, nil
}

func (r *mealRepository) FindByCode(ctx context.Context, code string) (*models.Meal, error) {
	var m models.Meal
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, fmt.Errorf("MealRepo.FindByCode: %w", err)
	}
	return &m, nil
}

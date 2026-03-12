package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"
)

type MealRepository interface {
	FindAll(ctx context.Context) ([]models.Meal, error)
	FindByCode(ctx context.Context, code string) (*models.Meal, error)
}

type mealRepository struct{  }

func NewMealRepository() MealRepository               { return &mealRepository{} }


func (r *mealRepository) FindAll(ctx context.Context) ([]models.Meal, error) {
	var meals []models.Meal
	if err := db.DB.WithContext(ctx).Order("code").Find(&meals).Error; err != nil {
		return nil, fmt.Errorf("MealRepo.FindAll: %w", err)
	}
	return meals, nil
}

func (r *mealRepository) FindByCode(ctx context.Context, code string) (*models.Meal, error) {
	var m models.Meal
	if err := db.DB.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, fmt.Errorf("MealRepo.FindByCode: %w", err)
	}
	return &m, nil
}

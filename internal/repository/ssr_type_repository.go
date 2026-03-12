package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"
)

type SSRTypeRepository interface {
	FindAll(ctx context.Context) ([]models.SSRType, error)
	FindByCode(ctx context.Context, code string) (*models.SSRType, error)
}

type ssrTypeRepositoryImpl struct{  }

func NewSSRTypeRepository() SSRTypeRepository{
	 return &ssrTypeRepositoryImpl{} }

func (r *ssrTypeRepositoryImpl) FindAll(ctx context.Context) ([]models.SSRType, error) {
	var types []models.SSRType
	if err := db.DB.WithContext(ctx).Order("code").Find(&types).Error; err != nil {
		return nil, fmt.Errorf("SSRTypeRepo.FindAll: %w", err)
	}
	return types, nil
}
func (r *ssrTypeRepositoryImpl) FindByCode(ctx context.Context, code string) (*models.SSRType, error) {
	var t models.SSRType
	if err := db.DB.WithContext(ctx).Where("code = ?", code).First(&t).Error; err != nil {
		return nil, fmt.Errorf("SSRTypeRepo.FindByCode: %w", err)
	}
	return &t, nil
}

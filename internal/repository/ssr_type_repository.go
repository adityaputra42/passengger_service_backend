package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"gorm.io/gorm"
)

type SSRTypeRepository interface {
	FindAll(ctx context.Context) ([]models.SSRType, error)
	FindByCode(ctx context.Context, code string) (*models.SSRType, error)
}

type ssrTypeRepositoryImpl struct {
	db *gorm.DB
}

func NewSSRTypeRepository(db *gorm.DB) SSRTypeRepository {
	return &ssrTypeRepositoryImpl{db: db}
}

func (r *ssrTypeRepositoryImpl) FindAll(ctx context.Context) ([]models.SSRType, error) {
	var types []models.SSRType
	if err := r.db.WithContext(ctx).Order("code").Find(&types).Error; err != nil {
		return nil, fmt.Errorf("SSRTypeRepo.FindAll: %w", err)
	}
	return types, nil
}

func (r *ssrTypeRepositoryImpl) FindByCode(ctx context.Context, code string) (*models.SSRType, error) {
	var t models.SSRType
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&t).Error; err != nil {
		return nil, fmt.Errorf("SSRTypeRepo.FindByCode: %w", err)
	}
	return &t, nil
}

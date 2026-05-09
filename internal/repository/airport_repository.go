package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AirportRepository interface {
	Create(ctx context.Context, airport *models.Airport) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Airport, error)
	FindByCode(ctx context.Context, code string) (*models.Airport, error)
	FindAll(ctx context.Context) ([]models.Airport, error)
	Search(ctx context.Context, query string) ([]models.Airport, error)
	Update(ctx context.Context, airport *models.Airport) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AirportRepositoryImpl struct {
	db *gorm.DB
}

func NewAirportRepository(db *gorm.DB) *AirportRepositoryImpl {
	return &AirportRepositoryImpl{db: db}
}

func (r *AirportRepositoryImpl) Create(ctx context.Context, airport *models.Airport) error {
	if err := r.db.WithContext(ctx).Create(airport).Error; err != nil {
		return fmt.Errorf("AirportRepo.Create: %w", err)
	}
	return nil
}

func (r *AirportRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Airport, error) {
	var airport models.Airport
	if err := r.db.WithContext(ctx).First(&airport, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("AirportRepo.FindByID: %w", err)
	}
	return &airport, nil
}

func (r *AirportRepositoryImpl) FindByCode(ctx context.Context, code string) (*models.Airport, error) {
	var airport models.Airport
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&airport).Error; err != nil {
		return nil, fmt.Errorf("AirportRepo.FindByCode: %w", err)
	}
	return &airport, nil
}

func (r *AirportRepositoryImpl) FindAll(ctx context.Context) ([]models.Airport, error) {
	var airports []models.Airport
	if err := r.db.WithContext(ctx).Order("country, city").Find(&airports).Error; err != nil {
		return nil, fmt.Errorf("AirportRepo.FindAll: %w", err)
	}
	return airports, nil
}

func (r *AirportRepositoryImpl) Search(ctx context.Context, query string) ([]models.Airport, error) {
	var airports []models.Airport
	like := "%" + query + "%"
	if err := r.db.WithContext(ctx).
		Where("code ILIKE ? OR name ILIKE ? OR city ILIKE ? OR country ILIKE ?", like, like, like, like).
		Limit(20).
		Find(&airports).Error; err != nil {
		return nil, fmt.Errorf("AirportRepo.Search: %w", err)
	}
	return airports, nil
}

func (r *AirportRepositoryImpl) Update(ctx context.Context, airport *models.Airport) error {
	if err := r.db.WithContext(ctx).Save(airport).Error; err != nil {
		return fmt.Errorf("AirportRepo.Update: %w", err)
	}
	return nil
}

func (r *AirportRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Airport{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("AirportRepo.Delete: %w", err)
	}
	return nil
}

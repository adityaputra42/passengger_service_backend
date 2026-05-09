package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatClassRepository interface {
	FindAll(ctx context.Context) ([]models.SeatClass, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.SeatClass, error)
	FindByCode(ctx context.Context, tx *gorm.DB, code string) (*models.SeatClass, error)
}

// ─────────────────────────────────────────────
// Implementation
// ─────────────────────────────────────────────

type SeatClassRepositoryImpl struct {
	db *gorm.DB
}

func NewSeatClassRepository(db *gorm.DB) SeatClassRepository {
	return &SeatClassRepositoryImpl{db: db}
}

func (r *SeatClassRepositoryImpl) FindAll(ctx context.Context) ([]models.SeatClass, error) {
	var classes []models.SeatClass
	if err := r.db.WithContext(ctx).Order("code").Find(&classes).Error; err != nil {
		return nil, fmt.Errorf("SeatClassRepo.FindAll: %w", err)
	}
	return classes, nil
}

func (r *SeatClassRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.SeatClass, error) {
	var sc models.SeatClass
	if err := r.db.WithContext(ctx).First(&sc, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("SeatClassRepo.FindByID: %w", err)
	}
	return &sc, nil
}

func (r *SeatClassRepositoryImpl) FindByCode(ctx context.Context, tx *gorm.DB, code string) (*models.SeatClass, error) {
	var sc models.SeatClass
	database := db.DB

	if tx != nil {
		database = tx
	}
	if err := database.WithContext(ctx).Where("code = ?", code).First(&sc).Error; err != nil {
		return nil, fmt.Errorf("SeatClassRepo.FindByCode %q: %w", code, err)
	}
	return &sc, nil
}

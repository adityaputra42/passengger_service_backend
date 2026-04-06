package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PNRContactRepository interface {
	Create(ctx context.Context, contact *models.PNRContact) error
	FindByPNRID(ctx context.Context, pnrID uuid.UUID) (*models.PNRContact, error)
	Update(ctx context.Context, contact *models.PNRContact) error
}

type pnrContactRepository struct {
	db *gorm.DB
}

func NewPNRContactRepository(db *gorm.DB) PNRContactRepository {
	return &pnrContactRepository{db: db}
}

func (r *pnrContactRepository) Create(ctx context.Context, contact *models.PNRContact) error {
	if err := r.db.WithContext(ctx).Create(contact).Error; err != nil {
		return fmt.Errorf("PNRContactRepo.Create: %w", err)
	}
	return nil
}

func (r *pnrContactRepository) FindByPNRID(ctx context.Context, pnrID uuid.UUID) (*models.PNRContact, error) {
	var contact models.PNRContact
	if err := r.db.WithContext(ctx).
		Where("pnr_id = ?", pnrID).
		First(&contact).Error; err != nil {
		return nil, fmt.Errorf("PNRContactRepo.FindByPNRID: %w", err)
	}
	return &contact, nil
}

func (r *pnrContactRepository) Update(ctx context.Context, contact *models.PNRContact) error {
	if err := r.db.WithContext(ctx).Save(contact).Error; err != nil {
		return fmt.Errorf("PNRContactRepo.Update: %w", err)
	}
	return nil
}

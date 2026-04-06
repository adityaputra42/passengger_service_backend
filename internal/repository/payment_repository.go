package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	FindByPNRID(ctx context.Context, pnrID uuid.UUID) ([]models.Payment, error)
	FindByStatus(ctx context.Context, status models.PaymentStatus, page, limit int) ([]models.Payment, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.PaymentStatus) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, p *models.Payment) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return fmt.Errorf("PaymentRepo.Create: %w", err)
	}
	return nil
}

func (r *paymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	var p models.Payment
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("PaymentRepo.FindByID: %w", err)
	}
	return &p, nil
}

func (r *paymentRepository) FindByPNRID(ctx context.Context, pnrID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.WithContext(ctx).
		Where("pnr_id = ?", pnrID).
		Order("created_at DESC").
		Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("PaymentRepo.FindByPNRID: %w", err)
	}
	return payments, nil
}

func (r *paymentRepository) FindByStatus(ctx context.Context, status models.PaymentStatus, page, limit int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	offset := (page - 1) * limit
	q := r.db.WithContext(ctx).Model(&models.Payment{}).Where("status = ?", status)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("PaymentRepo.FindByStatus count: %w", err)
	}
	if err := q.Preload("PNR").
		Order("paid_at DESC").
		Offset(offset).Limit(limit).
		Find(&payments).Error; err != nil {
		return nil, 0, fmt.Errorf("PaymentRepo.FindByStatus: %w", err)
	}
	return payments, total, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PaymentStatus) error {
	if err := r.db.WithContext(ctx).
		Model(&models.Payment{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("PaymentRepo.UpdateStatus: %w", err)
	}
	return nil
}

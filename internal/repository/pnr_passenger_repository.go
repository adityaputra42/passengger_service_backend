package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PNRPassengerRepository interface {
	Create(ctx context.Context, passenger *models.PNRPassenger) error
	BulkCreate(ctx context.Context, passengers []models.PNRPassenger) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.PNRPassenger, error)
	FindByPNRID(ctx context.Context, pnrID uuid.UUID) ([]models.PNRPassenger, error)
	Update(ctx context.Context, passenger *models.PNRPassenger) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type pnrPassengerRepository struct {
	db *gorm.DB
}

func NewPNRPassengerRepository(db *gorm.DB) PNRPassengerRepository {
	return &pnrPassengerRepository{db: db}
}

func (r *pnrPassengerRepository) Create(ctx context.Context, p *models.PNRPassenger) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return fmt.Errorf("PNRPassengerRepo.Create: %w", err)
	}
	return nil
}

func (r *pnrPassengerRepository) BulkCreate(ctx context.Context, passengers []models.PNRPassenger) error {
	if len(passengers) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).CreateInBatches(passengers, 50).Error; err != nil {
		return fmt.Errorf("PNRPassengerRepo.BulkCreate: %w", err)
	}
	return nil
}

func (r *pnrPassengerRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PNRPassenger, error) {
	var p models.PNRPassenger
	if err := r.db.WithContext(ctx).
		Preload("SeatAssignment").
		Preload("Ticket").
		First(&p, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("PNRPassengerRepo.FindByID: %w", err)
	}
	return &p, nil
}

func (r *pnrPassengerRepository) FindByPNRID(ctx context.Context, pnrID uuid.UUID) ([]models.PNRPassenger, error) {
	var passengers []models.PNRPassenger
	if err := r.db.WithContext(ctx).
		Preload("SeatAssignment").
		Preload("SeatAssignment.FlightSeat").
		Preload("SeatAssignment.FlightSeat.AircraftSeat").
		Preload("Ticket").
		Where("pnr_id = ?", pnrID).
		Find(&passengers).Error; err != nil {
		return nil, fmt.Errorf("PNRPassengerRepo.FindByPNRID: %w", err)
	}
	return passengers, nil
}

func (r *pnrPassengerRepository) Update(ctx context.Context, p *models.PNRPassenger) error {
	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return fmt.Errorf("PNRPassengerRepo.Update: %w", err)
	}
	return nil
}

func (r *pnrPassengerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.PNRPassenger{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("PNRPassengerRepo.Delete: %w", err)
	}
	return nil
}

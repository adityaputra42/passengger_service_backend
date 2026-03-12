package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type PNRSegmentRepository interface {
	Create(ctx context.Context, segment *models.PNRSegment) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.PNRSegment, error)
	FindByPNRID(ctx context.Context, pnrID uuid.UUID) ([]models.PNRSegment, error)
	FindWithFlight(ctx context.Context, id uuid.UUID) (*models.PNRSegment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type pnrSegmentRepository struct {

}

func NewPNRSegmentRepository() PNRSegmentRepository {
	return &pnrSegmentRepository{}
}

func (r *pnrSegmentRepository) Create(ctx context.Context, s *models.PNRSegment) error {
	if err := db.DB.WithContext(ctx).Create(s).Error; err != nil {
		return fmt.Errorf("PNRSegmentRepo.Create: %w", err)
	}
	return nil
}

func (r *pnrSegmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PNRSegment, error) {
	var s models.PNRSegment
	if err := db.DB.WithContext(ctx).First(&s, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("PNRSegmentRepo.FindByID: %w", err)
	}
	return &s, nil
}

func (r *pnrSegmentRepository) FindByPNRID(ctx context.Context, pnrID uuid.UUID) ([]models.PNRSegment, error) {
	var segments []models.PNRSegment
	if err := db.DB.WithContext(ctx).
		Preload("Flight").
		Preload("Flight.Schedule").
		Preload("Flight.Schedule.DepartureAirport").
		Preload("Flight.Schedule.ArrivalAirport").
		Where("pnr_id = ?", pnrID).
		Order("segment_order ASC").
		Find(&segments).Error; err != nil {
		return nil, fmt.Errorf("PNRSegmentRepo.FindByPNRID: %w", err)
	}
	return segments, nil
}

func (r *pnrSegmentRepository) FindWithFlight(ctx context.Context, id uuid.UUID) (*models.PNRSegment, error) {
	var s models.PNRSegment
	if err := db.DB.WithContext(ctx).
		Preload("Flight").
		Preload("Flight.Schedule").
		Preload("Flight.Schedule.DepartureAirport").
		Preload("Flight.Schedule.ArrivalAirport").
		Preload("Flight.Aircraft").
		First(&s, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("PNRSegmentRepo.FindWithFlight: %w", err)
	}
	return &s, nil
}

func (r *pnrSegmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := db.DB.WithContext(ctx).Delete(&models.PNRSegment{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("PNRSegmentRepo.Delete: %w", err)
	}
	return nil
}

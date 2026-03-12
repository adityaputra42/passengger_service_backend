package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type BoardingPassRepository interface {
	Create(ctx context.Context, bp *models.BoardingPass) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.BoardingPass, error)
	FindByPassengerID(ctx context.Context, passengerID uuid.UUID) ([]models.BoardingPass, error)
	FindBySegmentID(ctx context.Context, segmentID uuid.UUID) ([]models.BoardingPass, error)
	FindByPassengerAndSegment(ctx context.Context, passengerID, segmentID uuid.UUID) (*models.BoardingPass, error)
}

type boardingPassRepository struct {

}

func NewBoardingPassRepository() BoardingPassRepository {
	return &boardingPassRepository{}
}

func (r *boardingPassRepository) Create(ctx context.Context, bp *models.BoardingPass) error {
	if err := db.DB.WithContext(ctx).Create(bp).Error; err != nil {
		return fmt.Errorf("BoardingPassRepo.Create: %w", err)
	}
	return nil
}

func (r *boardingPassRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.BoardingPass, error) {
	var bp models.BoardingPass
	if err := db.DB.WithContext(ctx).
		Preload("Passenger").
		Preload("Segment").
		Preload("Segment.Flight").
		Preload("Segment.Flight.Schedule").
		Preload("Segment.Flight.Schedule.DepartureAirport").
		Preload("Segment.Flight.Schedule.ArrivalAirport").
		First(&bp, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("BoardingPassRepo.FindByID: %w", err)
	}
	return &bp, nil
}

func (r *boardingPassRepository) FindByPassengerID(ctx context.Context, passengerID uuid.UUID) ([]models.BoardingPass, error) {
	var bps []models.BoardingPass
	if err := db.DB.WithContext(ctx).
		Preload("Segment").
		Where("passenger_id = ?", passengerID).
		Find(&bps).Error; err != nil {
		return nil, fmt.Errorf("BoardingPassRepo.FindByPassengerID: %w", err)
	}
	return bps, nil
}

func (r *boardingPassRepository) FindBySegmentID(ctx context.Context, segmentID uuid.UUID) ([]models.BoardingPass, error) {
	var bps []models.BoardingPass
	if err := db.DB.WithContext(ctx).
		Preload("Passenger").
		Where("segment_id = ?", segmentID).
		Find(&bps).Error; err != nil {
		return nil, fmt.Errorf("BoardingPassRepo.FindBySegmentID: %w", err)
	}
	return bps, nil
}

func (r *boardingPassRepository) FindByPassengerAndSegment(ctx context.Context, passengerID, segmentID uuid.UUID) (*models.BoardingPass, error) {
	var bp models.BoardingPass
	if err := db.DB.WithContext(ctx).
		Preload("Segment").
		Preload("Segment.Flight").
		Preload("Segment.Flight.Schedule").
		Preload("Segment.Flight.Schedule.DepartureAirport").
		Preload("Segment.Flight.Schedule.ArrivalAirport").
		Where("passenger_id = ? AND segment_id = ?", passengerID, segmentID).
		First(&bp).Error; err != nil {
		return nil, fmt.Errorf("BoardingPassRepo.FindByPassengerAndSegment: %w", err)
	}
	return &bp, nil
}

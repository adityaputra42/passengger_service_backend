package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type PNRRepository interface {
	Create(ctx context.Context, pnr *models.PNR) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.PNR, error)
	FindByLocator(ctx context.Context, locator string) (*models.PNR, error)
	FindWithFull(ctx context.Context, id uuid.UUID) (*models.PNR, error) // preload all relations
	FindAll(ctx context.Context, page, limit int) ([]models.PNR, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.PNRStatus) error
	Update(ctx context.Context, pnr *models.PNR) error
}



type pnrRepository struct {

}

func NewPNRRepository() PNRRepository {
	return &pnrRepository{}
}

func (r *pnrRepository) Create(ctx context.Context, pnr *models.PNR) error {
	if err := db.DB.WithContext(ctx).Create(pnr).Error; err != nil {
		return fmt.Errorf("PNRRepo.Create: %w", err)
	}
	return nil
}

func (r *pnrRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PNR, error) {
	var pnr models.PNR
	if err := db.DB.WithContext(ctx).First(&pnr, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("PNRRepo.FindByID: %w", err)
	}
	return &pnr, nil
}

func (r *pnrRepository) FindByLocator(ctx context.Context, locator string) (*models.PNR, error) {
	var pnr models.PNR
	if err := db.DB.WithContext(ctx).
		Where("record_locator = ?", locator).
		First(&pnr).Error; err != nil {
		return nil, fmt.Errorf("PNRRepo.FindByLocator: %w", err)
	}
	return &pnr, nil
}

// FindWithFull loads PNR with all relations needed for booking summary.
func (r *pnrRepository) FindWithFull(ctx context.Context, id uuid.UUID) (*models.PNR, error) {
	var pnr models.PNR
	if err := db.DB.WithContext(ctx).
		Preload("Contact").
		Preload("Passengers").
		Preload("Passengers.SeatAssignment").
		Preload("Passengers.SeatAssignment.FlightSeat").
		Preload("Passengers.SeatAssignment.FlightSeat.AircraftSeat").
		Preload("Passengers.SeatAssignment.FlightSeat.AircraftSeat.SeatClass").
		Preload("Passengers.Ticket").
		Preload("Passengers.SSRs").
		Preload("Passengers.SSRs.SSRType").
		Preload("Passengers.Meals").
		Preload("Passengers.Meals.Meal").
		Preload("Passengers.Baggage").
		Preload("Segments").
		Preload("Segments.Flight").
		Preload("Segments.Flight.Schedule").
		Preload("Segments.Flight.Schedule.DepartureAirport").
		Preload("Segments.Flight.Schedule.ArrivalAirport").
		Preload("Payments").
		First(&pnr, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("PNRRepo.FindWithFull: %w", err)
	}
	return &pnr, nil
}

func (r *pnrRepository) FindAll(ctx context.Context, page, limit int) ([]models.PNR, int64, error) {
	var pnrs []models.PNR
	var total int64

	offset := (page - 1) * limit
	q := db.DB.WithContext(ctx).Model(&models.PNR{})

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("PNRRepo.FindAll count: %w", err)
	}
	if err := q.
		Preload("Contact").
		Preload("Passengers").
		Preload("Segments").
		Preload("Segments.Flight").
		Preload("Segments.Flight.Schedule").
		Preload("Segments.Flight.Schedule.DepartureAirport").
		Preload("Segments.Flight.Schedule.ArrivalAirport").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&pnrs).Error; err != nil {
		return nil, 0, fmt.Errorf("PNRRepo.FindAll: %w", err)
	}
	return pnrs, total, nil
}

func (r *pnrRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PNRStatus) error {
	if err := db.DB.WithContext(ctx).
		Model(&models.PNR{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("PNRRepo.UpdateStatus: %w", err)
	}
	return nil
}

func (r *pnrRepository) Update(ctx context.Context, pnr *models.PNR) error {
	if err := db.DB.WithContext(ctx).Save(pnr).Error; err != nil {
		return fmt.Errorf("PNRRepo.Update: %w", err)
	}
	return nil
}

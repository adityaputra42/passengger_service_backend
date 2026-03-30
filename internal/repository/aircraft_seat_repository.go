package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AircraftSeatRepository interface {
	BulkCreate(ctx context.Context, tx *gorm.DB, seats []models.AircraftSeat) error
	FindByAircraftID(ctx context.Context, aircraftID uuid.UUID) ([]models.AircraftSeat, error)
	FindByAircraftAndClass(ctx context.Context, aircraftID uuid.UUID, classCode string) ([]models.AircraftSeat, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.AircraftSeat, error)
	Update(ctx context.Context, seat *models.AircraftSeat) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type aircraftSeatRepository struct {
}

func NewAircraftSeatRepository() AircraftSeatRepository {
	return &aircraftSeatRepository{}
}

func (r *aircraftSeatRepository) BulkCreate(ctx context.Context, tx *gorm.DB, seats []models.AircraftSeat) error {
	if len(seats) == 0 {
		return nil
	}
	database := db.DB

	if tx != nil {
		database = tx
	}

	if err := database.WithContext(ctx).CreateInBatches(seats, 100).Error; err != nil {
		return fmt.Errorf("AircraftSeatRepo.BulkCreate: %w", err)
	}
	return nil
}

func (r *aircraftSeatRepository) FindByAircraftID(ctx context.Context, aircraftID uuid.UUID) ([]models.AircraftSeat, error) {
	var seats []models.AircraftSeat
	if err := db.DB.WithContext(ctx).
		Preload("SeatClass").
		Where("aircraft_id = ?", aircraftID).
		Order("row_number, seat_letter").
		Find(&seats).Error; err != nil {
		return nil, fmt.Errorf("AircraftSeatRepo.FindByAircraftID: %w", err)
	}
	return seats, nil
}

func (r *aircraftSeatRepository) FindByAircraftAndClass(ctx context.Context, aircraftID uuid.UUID, classCode string) ([]models.AircraftSeat, error) {
	var seats []models.AircraftSeat
	if err := db.DB.WithContext(ctx).
		Joins("JOIN seat_classes ON seat_classes.id = aircraft_seats.seat_class_id").
		Where("aircraft_seats.aircraft_id = ? AND seat_classes.code = ?", aircraftID, classCode).
		Order("aircraft_seats.row_number, aircraft_seats.seat_letter").
		Find(&seats).Error; err != nil {
		return nil, fmt.Errorf("AircraftSeatRepo.FindByAircraftAndClass: %w", err)
	}
	return seats, nil
}

func (r *aircraftSeatRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.AircraftSeat, error) {
	var seat models.AircraftSeat
	if err := db.DB.WithContext(ctx).
		Preload("SeatClass").
		First(&seat, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("AircraftSeatRepo.FindByID: %w", err)
	}
	return &seat, nil
}

func (r *aircraftSeatRepository) Update(ctx context.Context, seat *models.AircraftSeat) error {
	if err := db.DB.WithContext(ctx).Save(seat).Error; err != nil {
		return fmt.Errorf("AircraftSeatRepo.Update: %w", err)
	}
	return nil
}

func (r *aircraftSeatRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := db.DB.WithContext(ctx).Delete(&models.AircraftSeat{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("AircraftSeatRepo.Delete: %w", err)
	}
	return nil
}

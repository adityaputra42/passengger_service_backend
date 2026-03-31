package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type FlightSeatRepository interface {
	BulkCreate(ctx context.Context, seats []models.FlightSeat) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.FlightSeat, error)
	FindByFlight(ctx context.Context, flightID uuid.UUID) ([]models.FlightSeat, error)
	FindAvailableByFlight(ctx context.Context, flightID uuid.UUID) ([]models.FlightSeat, error)
	FindAvailableByFlightAndClass(ctx context.Context, flightID uuid.UUID, classCode string) ([]models.FlightSeat, error)
	FindWithSeatDetail(ctx context.Context, id uuid.UUID) (*models.FlightSeat, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.FlightSeatStatus) error
	CountAvailable(ctx context.Context, flightID uuid.UUID) (int64, error)
}

type flightSeatRepository struct {
}

func NewFlightSeatRepository() FlightSeatRepository {
	return &flightSeatRepository{}
}

func (r *flightSeatRepository) BulkCreate(ctx context.Context, seats []models.FlightSeat) error {
	if len(seats) == 0 {
		return nil
	}
	if err := db.DB.WithContext(ctx).CreateInBatches(seats, 100).Error; err != nil {
		return fmt.Errorf("FlightSeatRepo.BulkCreate: %w", err)
	}
	return nil
}

func (r *flightSeatRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.FlightSeat, error) {
	var fs models.FlightSeat
	if err := db.DB.WithContext(ctx).First(&fs, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("FlightSeatRepo.FindByID: %w", err)
	}
	return &fs, nil
}

func (r *flightSeatRepository) FindByFlight(ctx context.Context, flightID uuid.UUID) ([]models.FlightSeat, error) {
	var seats []models.FlightSeat
	if err := db.DB.WithContext(ctx).
		Preload("AircraftSeat").
		Preload("AircraftSeat.SeatClass").
		Where("flight_id = ?", flightID).
		Order("aircraft_seats.row_number, aircraft_seats.seat_letter").
		Find(&seats).Error; err != nil {
		return nil, fmt.Errorf("FlightSeatRepo.FindByFlight: %w", err)
	}
	return seats, nil
}

func (r *flightSeatRepository) FindAvailableByFlight(ctx context.Context, flightID uuid.UUID) ([]models.FlightSeat, error) {
	var seats []models.FlightSeat
	if err := db.DB.WithContext(ctx).
		Preload("AircraftSeat").
		Preload("AircraftSeat.SeatClass").
		Where("flight_id = ? AND status = ?", flightID, models.FlightSeatAvailable).
		Find(&seats).Error; err != nil {
		return nil, fmt.Errorf("FlightSeatRepo.FindAvailableByFlight: %w", err)
	}
	return seats, nil
}

func (r *flightSeatRepository) FindAvailableByFlightAndClass(ctx context.Context, flightID uuid.UUID, classCode string) ([]models.FlightSeat, error) {
	var seats []models.FlightSeat
	if err := db.DB.WithContext(ctx).
		Preload("AircraftSeat").
		Preload("AircraftSeat.SeatClass").
		Joins("JOIN aircraft_seats ON aircraft_seats.id = flight_seats.aircraft_seat_id").
		Joins("JOIN seat_classes ON seat_classes.id = aircraft_seats.seat_class_id").
		Where("flight_seats.flight_id = ? AND flight_seats.status = ? AND seat_classes.code = ?",
			flightID, models.FlightSeatAvailable, classCode).
		Order("aircraft_seats.row_number, aircraft_seats.seat_letter").
		Find(&seats).Error; err != nil {
		return nil, fmt.Errorf("FlightSeatRepo.FindAvailableByFlightAndClass: %w", err)
	}
	return seats, nil
}

func (r *flightSeatRepository) FindWithSeatDetail(ctx context.Context, id uuid.UUID) (*models.FlightSeat, error) {
	var fs models.FlightSeat
	if err := db.DB.WithContext(ctx).
		Preload("AircraftSeat").
		Preload("AircraftSeat.SeatClass").
		Preload("Flight").
		First(&fs, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("FlightSeatRepo.FindWithSeatDetail: %w", err)
	}
	return &fs, nil
}

func (r *flightSeatRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.FlightSeatStatus) error {
	if err := db.DB.WithContext(ctx).
		Model(&models.FlightSeat{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("FlightSeatRepo.UpdateStatus: %w", err)
	}
	return nil
}

func (r *flightSeatRepository) CountAvailable(ctx context.Context, flightID uuid.UUID) (int64, error) {
	var count int64
	if err := db.DB.WithContext(ctx).
		Model(&models.FlightSeat{}).
		Where("flight_id = ? AND status = ?", flightID, models.FlightSeatAvailable).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("FlightSeatRepo.CountAvailable: %w", err)
	}
	return count, nil
}

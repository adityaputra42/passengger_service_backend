package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlightRepository interface {
	Create(ctx context.Context, flight *models.Flight) error
	BulkCreate(ctx context.Context, flights []models.Flight) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Flight, error)
	FindAvailable(ctx context.Context, depAirportID, arrAirportID uuid.UUID, date time.Time) ([]models.Flight, error)
	FindBySchedule(ctx context.Context, scheduleID uuid.UUID) ([]models.Flight, error)
	FindByStatus(ctx context.Context, status models.FlightStatus) ([]models.Flight, error)
	FindWithDetails(ctx context.Context, id uuid.UUID) (*models.Flight, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.FlightStatus) error
	Update(ctx context.Context, flight *models.Flight) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type flightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) FlightRepository {
	return &flightRepository{db: db}
}

func (r *flightRepository) Create(ctx context.Context, f *models.Flight) error {
	if err := r.db.WithContext(ctx).Create(f).Error; err != nil {
		return fmt.Errorf("FlightRepo.Create: %w", err)
	}
	return nil
}

func (r *flightRepository) BulkCreate(ctx context.Context, flights []models.Flight) error {
	if len(flights) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).CreateInBatches(flights, 50).Error; err != nil {
		return fmt.Errorf("FlightRepo.BulkCreate: %w", err)
	}
	return nil
}

func (r *flightRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	var f models.Flight
	if err := r.db.WithContext(ctx).First(&f, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("FlightRepo.FindByID: %w", err)
	}
	return &f, nil
}

// FindAvailable searches flights by route and date with available seats.
// Joins through schedule to get departure/arrival airport.
func (r *flightRepository) FindAvailable(ctx context.Context, depAirportID, arrAirportID uuid.UUID, date time.Time) ([]models.Flight, error) {
	var flights []models.Flight

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	if err := r.db.WithContext(ctx).
		Preload("Schedule").
		Preload("Schedule.DepartureAirport").
		Preload("Schedule.ArrivalAirport").
		Preload("Aircraft").
		Joins("JOIN flight_schedules ON flight_schedules.id = flights.schedule_id").
		Where(`
			flight_schedules.departure_airport_id = ?
			AND flight_schedules.arrival_airport_id = ?
			AND flights.departure_time >= ?
			AND flights.departure_time < ?
			AND flights.status NOT IN (?)
		`, depAirportID, arrAirportID, startOfDay, endOfDay,
			[]string{string(models.FlightStatusCancelled), string(models.FlightStatusArrived)},
		).
		Order("flights.departure_time ASC").
		Find(&flights).Error; err != nil {
		return nil, fmt.Errorf("FlightRepo.FindAvailable: %w", err)
	}
	return flights, nil
}

func (r *flightRepository) FindBySchedule(ctx context.Context, scheduleID uuid.UUID) ([]models.Flight, error) {
	var flights []models.Flight
	if err := r.db.WithContext(ctx).
		Where("schedule_id = ?", scheduleID).
		Order("departure_time ASC").
		Find(&flights).Error; err != nil {
		return nil, fmt.Errorf("FlightRepo.FindBySchedule: %w", err)
	}
	return flights, nil
}

func (r *flightRepository) FindByStatus(ctx context.Context, status models.FlightStatus) ([]models.Flight, error) {
	var flights []models.Flight
	if err := r.db.WithContext(ctx).
		Preload("Schedule").
		Preload("Aircraft").
		Where("status = ?", status).
		Order("departure_time ASC").
		Find(&flights).Error; err != nil {
		return nil, fmt.Errorf("FlightRepo.FindByStatus: %w", err)
	}
	return flights, nil
}

// FindWithDetails loads full flight info: schedule, airports, aircraft, seat summary.
func (r *flightRepository) FindWithDetails(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	var f models.Flight
	if err := r.db.WithContext(ctx).
		Preload("Schedule").
		Preload("Schedule.DepartureAirport").
		Preload("Schedule.ArrivalAirport").
		Preload("Aircraft").
		First(&f, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("FlightRepo.FindWithDetails: %w", err)
	}
	return &f, nil
}

func (r *flightRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.FlightStatus) error {
	if err := r.db.WithContext(ctx).
		Model(&models.Flight{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("FlightRepo.UpdateStatus: %w", err)
	}
	return nil
}

func (r *flightRepository) Update(ctx context.Context, f *models.Flight) error {
	if err := r.db.WithContext(ctx).Save(f).Error; err != nil {
		return fmt.Errorf("FlightRepo.Update: %w", err)
	}
	return nil
}

func (r *flightRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Flight{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("FlightRepo.Delete: %w", err)
	}
	return nil
}

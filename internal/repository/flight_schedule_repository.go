package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)


type FlightScheduleRepository interface {
	Create(ctx context.Context, schedule *models.FlightSchedule) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.FlightSchedule, error)
	FindAll(ctx context.Context) ([]models.FlightSchedule, error)
	FindByRoute(ctx context.Context, depAirportID, arrAirportID uuid.UUID) ([]models.FlightSchedule, error)
	FindByFlightNumber(ctx context.Context, flightNumber string) (*models.FlightSchedule, error)
	Update(ctx context.Context, schedule *models.FlightSchedule) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type flightScheduleRepository struct {

}

func NewFlightScheduleRepository() FlightScheduleRepository {
	return &flightScheduleRepository{}
}

func (r *flightScheduleRepository) Create(ctx context.Context, s *models.FlightSchedule) error {
	if err := db.DB.WithContext(ctx).Create(s).Error; err != nil {
		return fmt.Errorf("FlightScheduleRepo.Create: %w", err)
	}
	return nil
}

func (r *flightScheduleRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.FlightSchedule, error) {
	var s models.FlightSchedule
	if err := db.DB.WithContext(ctx).
		Preload("DepartureAirport").
		Preload("ArrivalAirport").
		First(&s, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("FlightScheduleRepo.FindByID: %w", err)
	}
	return &s, nil
}

func (r *flightScheduleRepository) FindAll(ctx context.Context) ([]models.FlightSchedule, error) {
	var schedules []models.FlightSchedule
	if err := db.DB.WithContext(ctx).
		Preload("DepartureAirport").
		Preload("ArrivalAirport").
		Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("FlightScheduleRepo.FindAll: %w", err)
	}
	return schedules, nil
}

func (r *flightScheduleRepository) FindByRoute(ctx context.Context, depID, arrID uuid.UUID) ([]models.FlightSchedule, error) {
	var schedules []models.FlightSchedule
	if err := db.DB.WithContext(ctx).
		Preload("DepartureAirport").
		Preload("ArrivalAirport").
		Where("departure_airport_id = ? AND arrival_airport_id = ?", depID, arrID).
		Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("FlightScheduleRepo.FindByRoute: %w", err)
	}
	return schedules, nil
}

func (r *flightScheduleRepository) FindByFlightNumber(ctx context.Context, number string) (*models.FlightSchedule, error) {
	var s models.FlightSchedule
	if err := db.DB.WithContext(ctx).
		Preload("DepartureAirport").
		Preload("ArrivalAirport").
		Where("flight_number = ?", number).
		First(&s).Error; err != nil {
		return nil, fmt.Errorf("FlightScheduleRepo.FindByFlightNumber: %w", err)
	}
	return &s, nil
}

func (r *flightScheduleRepository) Update(ctx context.Context, s *models.FlightSchedule) error {
	if err := db.DB.WithContext(ctx).Save(s).Error; err != nil {
		return fmt.Errorf("FlightScheduleRepo.Update: %w", err)
	}
	return nil
}

func (r *flightScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := db.DB.WithContext(ctx).Delete(&models.FlightSchedule{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("FlightScheduleRepo.Delete: %w", err)
	}
	return nil
}

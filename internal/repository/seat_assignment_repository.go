package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
)

type SeatAssignmentRepository interface {
	Create(ctx context.Context, assignment *models.SeatAssignment) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.SeatAssignment, error)
	FindByPassengerID(ctx context.Context, passengerID uuid.UUID) (*models.SeatAssignment, error)
	FindBySegmentID(ctx context.Context, segmentID uuid.UUID) ([]models.SeatAssignment, error)
	FindByFlightSeatID(ctx context.Context, flightSeatID uuid.UUID) (*models.SeatAssignment, error)
	Update(ctx context.Context, assignment *models.SeatAssignment) error
	Delete(ctx context.Context, id uuid.UUID) error
}


type seatAssignmentRepository struct {

}

func NewSeatAssignmentRepository() SeatAssignmentRepository {
	return &seatAssignmentRepository{}
}

func (r *seatAssignmentRepository) Create(ctx context.Context, a *models.SeatAssignment) error {
	if err := db.DB.WithContext(ctx).Create(a).Error; err != nil {
		return fmt.Errorf("SeatAssignmentRepo.Create: %w", err)
	}
	return nil
}

func (r *seatAssignmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.SeatAssignment, error) {
	var a models.SeatAssignment
	if err := db.DB.WithContext(ctx).
		Preload("Passenger").
		Preload("FlightSeat").
		Preload("FlightSeat.AircraftSeat").
		Preload("FlightSeat.AircraftSeat.SeatClass").
		First(&a, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("SeatAssignmentRepo.FindByID: %w", err)
	}
	return &a, nil
}

func (r *seatAssignmentRepository) FindByPassengerID(ctx context.Context, passengerID uuid.UUID) (*models.SeatAssignment, error) {
	var a models.SeatAssignment
	if err := db.DB.WithContext(ctx).
		Preload("FlightSeat").
		Preload("FlightSeat.AircraftSeat").
		Preload("FlightSeat.AircraftSeat.SeatClass").
		Where("passenger_id = ?", passengerID).
		First(&a).Error; err != nil {
		return nil, fmt.Errorf("SeatAssignmentRepo.FindByPassengerID: %w", err)
	}
	return &a, nil
}

func (r *seatAssignmentRepository) FindBySegmentID(ctx context.Context, segmentID uuid.UUID) ([]models.SeatAssignment, error) {
	var assignments []models.SeatAssignment
	if err := db.DB.WithContext(ctx).
		Preload("Passenger").
		Preload("FlightSeat").
		Preload("FlightSeat.AircraftSeat").
		Where("segment_id = ?", segmentID).
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("SeatAssignmentRepo.FindBySegmentID: %w", err)
	}
	return assignments, nil
}

func (r *seatAssignmentRepository) FindByFlightSeatID(ctx context.Context, flightSeatID uuid.UUID) (*models.SeatAssignment, error) {
	var a models.SeatAssignment
	if err := db.DB.WithContext(ctx).
		Where("flight_seat_id = ?", flightSeatID).
		First(&a).Error; err != nil {
		return nil, fmt.Errorf("SeatAssignmentRepo.FindByFlightSeatID: %w", err)
	}
	return &a, nil
}

func (r *seatAssignmentRepository) Update(ctx context.Context, a *models.SeatAssignment) error {
	if err := db.DB.WithContext(ctx).Save(a).Error; err != nil {
		return fmt.Errorf("SeatAssignmentRepo.Update: %w", err)
	}
	return nil
}

func (r *seatAssignmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := db.DB.WithContext(ctx).Delete(&models.SeatAssignment{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("SeatAssignmentRepo.Delete: %w", err)
	}
	return nil
}

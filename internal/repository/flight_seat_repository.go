package repository

import (
	"context"
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

package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository interface {
	Create(ctx context.Context, ticket *models.Ticket) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error)
	FindByTicketNumber(ctx context.Context, number string) (*models.Ticket, error)
	FindByPassengerID(ctx context.Context, passengerID uuid.UUID) (*models.Ticket, error)
	FindWithSegments(ctx context.Context, id uuid.UUID) (*models.Ticket, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, t *models.Ticket) error {
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		return fmt.Errorf("TicketRepo.Create: %w", err)
	}
	return nil
}

func (r *ticketRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error) {
	var t models.Ticket
	if err := r.db.WithContext(ctx).First(&t, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("TicketRepo.FindByID: %w", err)
	}
	return &t, nil
}

func (r *ticketRepository) FindByTicketNumber(ctx context.Context, number string) (*models.Ticket, error) {
	var t models.Ticket
	if err := r.db.WithContext(ctx).
		Preload("Passenger").
		Preload("Segments").
		Preload("Segments.Segment").
		Preload("Segments.Segment.Flight").
		Preload("Segments.Segment.Flight.Schedule").
		Preload("Segments.Segment.Flight.Schedule.DepartureAirport").
		Preload("Segments.Segment.Flight.Schedule.ArrivalAirport").
		Where("ticket_number = ?", number).
		First(&t).Error; err != nil {
		return nil, fmt.Errorf("TicketRepo.FindByTicketNumber: %w", err)
	}
	return &t, nil
}

func (r *ticketRepository) FindByPassengerID(ctx context.Context, passengerID uuid.UUID) (*models.Ticket, error) {
	var t models.Ticket
	if err := r.db.WithContext(ctx).
		Preload("Segments").
		Where("passenger_id = ?", passengerID).
		First(&t).Error; err != nil {
		return nil, fmt.Errorf("TicketRepo.FindByPassengerID: %w", err)
	}
	return &t, nil
}

func (r *ticketRepository) FindWithSegments(ctx context.Context, id uuid.UUID) (*models.Ticket, error) {
	var t models.Ticket
	if err := r.db.WithContext(ctx).
		Preload("Passenger").
		Preload("Segments").
		Preload("Segments.Segment").
		Preload("Segments.Segment.Flight").
		Preload("Segments.Segment.Flight.Schedule").
		Preload("Segments.Segment.Flight.Schedule.DepartureAirport").
		Preload("Segments.Segment.Flight.Schedule.ArrivalAirport").
		First(&t, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("TicketRepo.FindWithSegments: %w", err)
	}
	return &t, nil
}

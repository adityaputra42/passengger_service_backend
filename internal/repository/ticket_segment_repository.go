package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketSegmentRepository interface {
	Create(ctx context.Context, ts *models.TicketSegment) error
	BulkCreate(ctx context.Context, segments []models.TicketSegment) error
	FindByTicketID(ctx context.Context, ticketID uuid.UUID) ([]models.TicketSegment, error)
}

type ticketSegmentRepository struct {
	db *gorm.DB
}

func NewTicketSegmentRepository(db *gorm.DB) TicketSegmentRepository {
	return &ticketSegmentRepository{db: db}
}

func (r *ticketSegmentRepository) Create(ctx context.Context, ts *models.TicketSegment) error {
	if err := r.db.WithContext(ctx).Create(ts).Error; err != nil {
		return fmt.Errorf("TicketSegmentRepo.Create: %w", err)
	}
	return nil
}

func (r *ticketSegmentRepository) BulkCreate(ctx context.Context, segments []models.TicketSegment) error {
	if len(segments) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).CreateInBatches(segments, 50).Error; err != nil {
		return fmt.Errorf("TicketSegmentRepo.BulkCreate: %w", err)
	}
	return nil
}

func (r *ticketSegmentRepository) FindByTicketID(ctx context.Context, ticketID uuid.UUID) ([]models.TicketSegment, error) {
	var segments []models.TicketSegment
	if err := r.db.WithContext(ctx).
		Preload("Segment").
		Preload("Segment.Flight").
		Preload("Segment.Flight.Schedule").
		Where("ticket_id = ?", ticketID).
		Find(&segments).Error; err != nil {
		return nil, fmt.Errorf("TicketSegmentRepo.FindByTicketID: %w", err)
	}
	return segments, nil
}

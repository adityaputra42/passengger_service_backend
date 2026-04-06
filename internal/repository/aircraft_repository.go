package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AircraftRepository interface {
	Create(ctx context.Context, aircraft *models.Aircraft) error
	FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*models.Aircraft, error)
	FindAll(ctx context.Context) ([]models.Aircraft, error)
	FindWithSeats(ctx context.Context, id uuid.UUID) (*models.Aircraft, error)
	Update(ctx context.Context, tx *gorm.DB, aircraft *models.Aircraft) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type aircraftRepository struct {
	db *gorm.DB
}

func NewAircraftRepository(db *gorm.DB) AircraftRepository {
	return &aircraftRepository{db: db}
}

func (r *aircraftRepository) Create(ctx context.Context, aircraft *models.Aircraft) error {
	if err := r.db.WithContext(ctx).Create(aircraft).Error; err != nil {
		return fmt.Errorf("AircraftRepo.Create: %w", err)
	}
	return nil
}

func (r *aircraftRepository) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*models.Aircraft, error) {
	var aircraft models.Aircraft
	database := r.db

	if tx != nil {
		database = tx
	}
	if err := database.WithContext(ctx).First(&aircraft, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("AircraftRepo.FindByID: %w", err)
	}
	return &aircraft, nil
}

func (r *aircraftRepository) FindAll(ctx context.Context) ([]models.Aircraft, error) {
	var aircrafts []models.Aircraft
	if err := r.db.WithContext(ctx).Order("manufacturer, model").Find(&aircrafts).Error; err != nil {
		return nil, fmt.Errorf("AircraftRepo.FindAll: %w", err)
	}
	return aircrafts, nil
}

func (r *aircraftRepository) FindWithSeats(ctx context.Context, id uuid.UUID) (*models.Aircraft, error) {
	var aircraft models.Aircraft
	if err := r.db.WithContext(ctx).
		Preload("Seats", func(db *gorm.DB) *gorm.DB {
			return db.Order("row_number, seat_letter")
		}).
		Preload("Seats.SeatClass").
		First(&aircraft, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("AircraftRepo.FindWithSeats: %w", err)
	}
	return &aircraft, nil
}

func (r *aircraftRepository) Update(ctx context.Context, tx *gorm.DB, aircraft *models.Aircraft) error {
	database := r.db
	if tx != nil {
		database = tx
	}
	if err := database.WithContext(ctx).Save(aircraft).Error; err != nil {
		return fmt.Errorf("AircraftRepo.Update: %w", err)
	}
	return nil
}

func (r *aircraftRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Aircraft{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("AircraftRepo.Delete: %w", err)
	}
	return nil
}

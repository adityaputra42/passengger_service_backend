package repository

import (
	"context"
	"errors"
	"passenger_service_backend/internal/cache"
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────
// CachedFlightSeatRepository
// ─────────────────────────────────────────────────────────────

type CachedFlightSeatRepository struct {
	inner *FlightSeatRepositoryImpl // ← concrete
	cache *cache.Client
}

func NewCachedFlightSeatRepository(
	inner *FlightSeatRepositoryImpl,
	cache *cache.Client,
) *CachedFlightSeatRepository {
	return &CachedFlightSeatRepository{
		inner: inner,
		cache: cache,
	}
}

func (r *CachedFlightSeatRepository) FindByFlight(ctx context.Context, flightID uuid.UUID) ([]models.FlightSeat, error) {
	key := cache.KeyFlightSeatMap(flightID)
	var seats []models.FlightSeat
	if err := r.cache.Get(ctx, key, &seats); err == nil {
		return seats, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindByFlight(ctx, flightID)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLFlightSeatMap*time.Second)
	return result, nil
}

func (r *CachedFlightSeatRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.FlightSeatStatus) error {
	// Fetch to get flightID before updating, so we can bust the right seat map key
	seat, lookupErr := r.inner.FindByID(ctx, id)

	err := r.inner.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	if lookupErr == nil && seat.FlightID != nil {
		_ = r.cache.Del(ctx, cache.KeyFlightSeatMap(*seat.FlightID))
	}
	return nil
}

func (r *CachedFlightSeatRepository) BulkCreate(ctx context.Context, seats []models.FlightSeat) error {
	err := r.inner.BulkCreate(ctx, seats)
	if err != nil {
		return err
	}
	if len(seats) > 0 && seats[0].FlightID != nil {
		_ = r.cache.Del(ctx, cache.KeyFlightSeatMap(*seats[0].FlightID))
	}
	return nil
}

// Pass-through methods
func (r *CachedFlightSeatRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.FlightSeat, error) {
	return r.inner.FindByID(ctx, id)
}
func (r *CachedFlightSeatRepository) FindAvailableByFlight(ctx context.Context, flightID uuid.UUID) ([]models.FlightSeat, error) {
	return r.inner.FindAvailableByFlight(ctx, flightID)
}
func (r *CachedFlightSeatRepository) FindAvailableByFlightAndClass(ctx context.Context, flightID uuid.UUID, classCode string) ([]models.FlightSeat, error) {
	return r.inner.FindAvailableByFlightAndClass(ctx, flightID, classCode)
}
func (r *CachedFlightSeatRepository) FindWithSeatDetail(ctx context.Context, id uuid.UUID) (*models.FlightSeat, error) {
	return r.inner.FindWithSeatDetail(ctx, id)
}
func (r *CachedFlightSeatRepository) CountAvailable(ctx context.Context, flightID uuid.UUID) (int64, error) {
	return r.inner.CountAvailable(ctx, flightID)
}

var _ FlightSeatRepository = (*CachedFlightSeatRepository)(nil)

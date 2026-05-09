package repository

import (
	"context"
	"errors"
	"fmt"
	"passenger_service_backend/internal/cache"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────
// CachedFlightRepository
//
// Cache strategy:
//   FindByID / FindWithDetails → cache per flight ID, 10 min
//   FindAvailable (search)     → cache per route+date, 5 min
//   UpdateStatus               → bust all keys for that flight
//   Create / BulkCreate        → bust search cache for affected routes
// ─────────────────────────────────────────────────────────────

type CachedFlightRepository struct {
	inner *FlightRepositoryImpl // ← concrete
	cache *cache.Client
}

func NewCachedFlightRepository(
	inner *FlightRepositoryImpl,
	cache *cache.Client,
) *CachedFlightRepository {
	return &CachedFlightRepository{
		inner: inner,
		cache: cache,
	}
}

func (r *CachedFlightRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	key := cache.KeyFlightByID(id)
	var flight models.Flight
	if err := r.cache.Get(ctx, key, &flight); err == nil {
		return &flight, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLFlightDetail*time.Second)
	return result, nil
}

func (r *CachedFlightRepository) FindWithDetails(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	key := fmt.Sprintf("flight:detail:%s", id)
	var flight models.Flight
	if err := r.cache.Get(ctx, key, &flight); err == nil {
		return &flight, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLFlightDetail*time.Second)
	return result, nil
}

func (r *CachedFlightRepository) FindAvailable(ctx context.Context, depID, arrID uuid.UUID, date time.Time) ([]models.Flight, error) {
	key := fmt.Sprintf("flight:search:%s:%s:%s", depID, arrID, date.Format("2006-01-02"))
	var flights []models.Flight
	if err := r.cache.Get(ctx, key, &flights); err == nil {
		return flights, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindAvailable(ctx, depID, arrID, date)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLFlightSearch*time.Second)
	return result, nil
}

func (r *CachedFlightRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.FlightStatus) error {
	err := r.inner.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx,
		cache.KeyFlightByID(id),
		fmt.Sprintf("flight:detail:%s", id),
		cache.KeyFlightSeatMap(id),
	)
	return nil
}

func (r *CachedFlightRepository) Update(ctx context.Context, flight *models.Flight) error {
	err := r.inner.Update(ctx, flight)
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx,
		cache.KeyFlightByID(flight.ID),
		fmt.Sprintf("flight:detail:%s", flight.ID),
	)
	return nil
}

func (r *CachedFlightRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.inner.Delete(ctx, id)
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx,
		cache.KeyFlightByID(id),
		fmt.Sprintf("flight:detail:%s", id),
		cache.KeyFlightSeatMap(id),
	)
	return nil
}

func (r *CachedFlightRepository) Create(ctx context.Context, f *models.Flight) error {
	return r.inner.Create(ctx, f)
}
func (r *CachedFlightRepository) BulkCreate(ctx context.Context, flights []models.Flight) error {
	return r.inner.BulkCreate(ctx, flights)
}
func (r *CachedFlightRepository) FindBySchedule(ctx context.Context, schedID uuid.UUID) ([]models.Flight, error) {
	return r.inner.FindBySchedule(ctx, schedID)
}
func (r *CachedFlightRepository) FindByStatus(ctx context.Context, status models.FlightStatus) ([]models.Flight, error) {
	return r.inner.FindByStatus(ctx, status)
}

var _ FlightRepository = (*CachedFlightRepository)(nil)

var _ = dto.FlightResult{}
var _ = errors.New

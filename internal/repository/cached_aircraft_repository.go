package repository

import (
	"context"
	"errors"
	"passenger_service_backend/internal/cache"
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CachedAircraftRepository struct {
	inner *AircraftRepositoryImpl // ← concrete
	cache *cache.Client
}

func NewCachedAircraftRepository(
	inner *AircraftRepositoryImpl,
	cache *cache.Client,
) *CachedAircraftRepository {
	return &CachedAircraftRepository{
		inner: inner,
		cache: cache,
	}
}

func (r *CachedAircraftRepository) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*models.Aircraft, error) {

	if tx != nil {
		return r.inner.FindByID(ctx, tx, id)
	}

	key := cache.KeyAircraftByID(id)
	var aircraft models.Aircraft
	if err := r.cache.Get(ctx, key, &aircraft); err == nil {
		return &aircraft, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindByID(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLAircraft*time.Second)
	return result, nil
}

func (r *CachedAircraftRepository) FindWithSeats(ctx context.Context, id uuid.UUID) (*models.Aircraft, error) {
	key := cache.KeyAircraftWithSeats(id)
	var aircraft models.Aircraft
	if err := r.cache.Get(ctx, key, &aircraft); err == nil {
		return &aircraft, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindWithSeats(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLAircraft*time.Second)
	return result, nil
}

func (r *CachedAircraftRepository) FindAll(ctx context.Context) ([]models.Aircraft, error) {
	key := cache.KeyAircraftList()
	var aircrafts []models.Aircraft
	if err := r.cache.Get(ctx, key, &aircrafts); err == nil {
		return aircrafts, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLAircraft*time.Second)
	return result, nil
}

func (r *CachedAircraftRepository) Create(ctx context.Context, aircraft *models.Aircraft) error {
	err := r.inner.Create(ctx, aircraft)
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx, cache.KeyAircraftList())
	return nil
}

func (r *CachedAircraftRepository) Update(ctx context.Context, tx *gorm.DB, aircraft *models.Aircraft) error {
	err := r.inner.Update(ctx, tx, aircraft)
	if err != nil {
		return err
	}
	if tx == nil {
		_ = r.cache.Del(ctx,
			cache.KeyAircraftByID(aircraft.ID),
			cache.KeyAircraftWithSeats(aircraft.ID),
			cache.KeyAircraftList(),
		)
	}
	return nil
}

func (r *CachedAircraftRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.inner.Delete(ctx, id)
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx,
		cache.KeyAircraftByID(id),
		cache.KeyAircraftWithSeats(id),
		cache.KeyAircraftList(),
	)
	return nil
}

var _ AircraftRepository = (*CachedAircraftRepository)(nil)

package repository

import (
	"context"
	"errors"
	"passenger_service_backend/internal/cache"
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type CachedAirportRepository struct {
	inner *AirportRepositoryImpl // ← concrete, not AirportRepository interface
	cache *cache.Client
}

func NewCachedAirportRepository(
	inner *AirportRepositoryImpl,
	cache *cache.Client,
) *CachedAirportRepository {
	return &CachedAirportRepository{
		inner: inner,
		cache: cache,
	}
}
func (r *CachedAirportRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Airport, error) {
	key := cache.KeyAirportByID(id)

	var airport models.Airport
	if err := r.cache.Get(ctx, key, &airport); err == nil {
		return &airport, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err // degraded
	}

	result, err := r.inner.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLAirport*time.Second)
	return result, nil
}

func (r *CachedAirportRepository) FindByCode(ctx context.Context, code string) (*models.Airport, error) {
	key := cache.KeyAirportByCode(code)

	var airport models.Airport
	if err := r.cache.Get(ctx, key, &airport); err == nil {
		return &airport, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLAirport*time.Second)
	return result, nil
}

func (r *CachedAirportRepository) FindAll(ctx context.Context) ([]models.Airport, error) {
	key := cache.KeyAirportList()

	var airports []models.Airport
	if err := r.cache.Get(ctx, key, &airports); err == nil {
		return airports, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLAirport*time.Second)
	return result, nil
}

func (r *CachedAirportRepository) Search(ctx context.Context, query string) ([]models.Airport, error) {
	// Not cached — queries are too varied, dataset is small
	return r.inner.Search(ctx, query)
}

func (r *CachedAirportRepository) Create(ctx context.Context, airport *models.Airport) error {
	err := r.inner.Create(ctx, airport)
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx, cache.KeyAirportList())
	return nil
}

func (r *CachedAirportRepository) Update(ctx context.Context, airport *models.Airport) error {
	err := r.inner.Update(ctx, airport)
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx,
		cache.KeyAirportByID(airport.ID),
		cache.KeyAirportByCode(airport.Code),
		cache.KeyAirportList(),
	)
	return nil
}

func (r *CachedAirportRepository) Delete(ctx context.Context, id uuid.UUID) error {
	airport, _ := r.inner.FindByID(ctx, id)
	err := r.inner.Delete(ctx, id)
	if err != nil {
		return err
	}
	keys := []string{cache.KeyAirportByID(id), cache.KeyAirportList()}
	if airport != nil {
		keys = append(keys, cache.KeyAirportByCode(airport.Code))
	}
	_ = r.cache.Del(ctx, keys...)
	return nil
}

var _ AirportRepository = (*CachedAirportRepository)(nil)

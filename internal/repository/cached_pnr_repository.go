package repository

import (
	"context"
	"errors"
	"passenger_service_backend/internal/cache"
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type CachedPNRRepository struct {
	inner *PnrRepositoryImpl
	cache *cache.Client
}

func NewCachedPNRRepository(
	inner *PnrRepositoryImpl,
	cache *cache.Client,
) *CachedPNRRepository {
	return &CachedPNRRepository{
		inner: inner,
		cache: cache,
	}
}

func (r *CachedPNRRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PNR, error) {
	key := cache.KeyPNRByID(id)
	var pnr models.PNR
	if err := r.cache.Get(ctx, key, &pnr); err == nil {
		return &pnr, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLBookingPNR*time.Second)
	return result, nil
}

func (r *CachedPNRRepository) FindByLocator(ctx context.Context, locator string) (*models.PNR, error) {
	key := cache.KeyPNRByLocator(locator)
	var pnr models.PNR
	if err := r.cache.Get(ctx, key, &pnr); err == nil {
		return &pnr, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindByLocator(ctx, locator)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLBookingPNR*time.Second)
	_ = r.cache.Set(ctx, cache.KeyPNRByID(result.ID), result, cache.TTLBookingPNR*time.Second)
	return result, nil
}

func (r *CachedPNRRepository) FindWithFull(ctx context.Context, id uuid.UUID) (*models.PNR, error) {
	key := "pnr:full:" + id.String()
	var pnr models.PNR
	if err := r.cache.Get(ctx, key, &pnr); err == nil {
		return &pnr, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindWithFull(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLBookingPNR*time.Second)
	return result, nil
}

func (r *CachedPNRRepository) FindAll(ctx context.Context, page, limit int) ([]models.PNR, int64, error) {
	// Not cached — paginated admin list, too varied
	return r.inner.FindAll(ctx, page, limit)
}

func (r *CachedPNRRepository) Create(ctx context.Context, pnr *models.PNR) error {
	return r.inner.Create(ctx, pnr)
}

func (r *CachedPNRRepository) Update(ctx context.Context, pnr *models.PNR) error {
	err := r.inner.Update(ctx, pnr)
	if err != nil {
		return err
	}
	r.bust(ctx, pnr.ID, pnr.RecordLocator)
	return nil
}

func (r *CachedPNRRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PNRStatus) error {
	err := r.inner.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx,
		cache.KeyPNRByID(id),
		"pnr:full:"+id.String(),
	)
	return nil
}

func (r *CachedPNRRepository) bust(ctx context.Context, id uuid.UUID, locator string) {
	_ = r.cache.Del(ctx,
		cache.KeyPNRByID(id),
		cache.KeyPNRByLocator(locator),
		"pnr:full:"+id.String(),
	)
}

var _ PNRRepository = (*CachedPNRRepository)(nil)

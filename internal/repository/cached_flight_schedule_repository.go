package repository

import (
	"context"
	"errors"
	"fmt"
	"passenger_service_backend/internal/cache"
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────
// CachedFlightScheduleRepository
// ─────────────────────────────────────────────────────────────

type CachedFlightScheduleRepository struct {
	inner *FlightScheduleRepositoryImpl // ← concrete
	cache *cache.Client
}

func NewCachedFlightScheduleRepository(
	inner *FlightScheduleRepositoryImpl,
	cache *cache.Client,
) *CachedFlightScheduleRepository {
	return &CachedFlightScheduleRepository{
		inner: inner,
		cache: cache,
	}
}

func (r *CachedFlightScheduleRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.FlightSchedule, error) {
	key := cache.KeyFlightScheduleByID(id)
	var sched models.FlightSchedule
	if err := r.cache.Get(ctx, key, &sched); err == nil {
		return &sched, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLFlightSchedule*time.Second)
	return result, nil
}

func (r *CachedFlightScheduleRepository) FindAll(ctx context.Context) ([]models.FlightSchedule, error) {
	key := cache.KeyFlightScheduleList()
	var scheds []models.FlightSchedule
	if err := r.cache.Get(ctx, key, &scheds); err == nil {
		return scheds, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLFlightSchedule*time.Second)
	return result, nil
}

func (r *CachedFlightScheduleRepository) FindByRoute(ctx context.Context, depID, arrID uuid.UUID) ([]models.FlightSchedule, error) {
	key := fmt.Sprintf("schedule:route:%s:%s", depID, arrID)
	var scheds []models.FlightSchedule
	if err := r.cache.Get(ctx, key, &scheds); err == nil {
		return scheds, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindByRoute(ctx, depID, arrID)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLFlightSchedule*time.Second)
	return result, nil
}

func (r *CachedFlightScheduleRepository) FindByFlightNumber(ctx context.Context, number string) (*models.FlightSchedule, error) {
	key := fmt.Sprintf("schedule:number:%s", number)
	var sched models.FlightSchedule
	if err := r.cache.Get(ctx, key, &sched); err == nil {
		return &sched, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	result, err := r.inner.FindByFlightNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, result, cache.TTLFlightSchedule*time.Second)
	return result, nil
}

func (r *CachedFlightScheduleRepository) Create(ctx context.Context, s *models.FlightSchedule) error {
	err := r.inner.Create(ctx, s)
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx, cache.KeyFlightScheduleList())
	return nil
}

func (r *CachedFlightScheduleRepository) Update(ctx context.Context, s *models.FlightSchedule) error {
	err := r.inner.Update(ctx, s)
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx,
		cache.KeyFlightScheduleByID(s.ID),
		cache.KeyFlightScheduleList(),
		fmt.Sprintf("schedule:route:%s:%s", s.DepartureAirportID, s.ArrivalAirportID),
		fmt.Sprintf("schedule:number:%s", s.FlightNumber),
	)
	return nil
}

func (r *CachedFlightScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sched, _ := r.inner.FindByID(ctx, id)
	err := r.inner.Delete(ctx, id)
	if err != nil {
		return err
	}
	keys := []string{cache.KeyFlightScheduleByID(id), cache.KeyFlightScheduleList()}
	if sched != nil {
		keys = append(keys,
			fmt.Sprintf("schedule:route:%s:%s", sched.DepartureAirportID, sched.ArrivalAirportID),
			fmt.Sprintf("schedule:number:%s", sched.FlightNumber),
		)
	}
	_ = r.cache.Del(ctx, keys...)
	return nil
}

var _ FlightScheduleRepository = (*CachedFlightScheduleRepository)(nil)

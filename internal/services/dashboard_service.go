package services

import (
	"context"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/repository"
)

type DashboardService interface {
	GetSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error)
	GetRevenueTrend(ctx context.Context, days int) ([]dto.RevenueTrendResponse, error)
	GetBookingStatus(ctx context.Context) ([]dto.BookingStatusResponse, error)
	GetTodayFlights(ctx context.Context) ([]dto.TodayFlightResponse, error)
	GetRecentBookings(ctx context.Context, limit int) ([]dto.RecentBookingResponse, error)
}

type DashboardServiceImpl struct {
	repo repository.DashboardRepository
}

func NewDashboardService(
	repo repository.DashboardRepository,
) DashboardService {
	return &DashboardServiceImpl{
		repo: repo,
	}
}

func (s *DashboardServiceImpl) GetSummary(
	ctx context.Context,
) (*dto.DashboardSummaryResponse, error) {
	return s.repo.GetSummary(ctx)
}

func (s *DashboardServiceImpl) GetRevenueTrend(
	ctx context.Context,
	days int,
) ([]dto.RevenueTrendResponse, error) {
	return s.repo.GetRevenueTrend(ctx, days)
}

func (s *DashboardServiceImpl) GetBookingStatus(
	ctx context.Context,
) ([]dto.BookingStatusResponse, error) {
	return s.repo.GetBookingStatus(ctx)
}

func (s *DashboardServiceImpl) GetTodayFlights(
	ctx context.Context,
) ([]dto.TodayFlightResponse, error) {
	return s.repo.GetTodayFlights(ctx)
}

func (s *DashboardServiceImpl) GetRecentBookings(
	ctx context.Context,
	limit int,
) ([]dto.RecentBookingResponse, error) {
	return s.repo.GetRecentBookings(ctx, limit)
}

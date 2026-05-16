package repository

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/dto"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error)
	GetRevenueTrend(ctx context.Context, days int) ([]dto.RevenueTrendResponse, error)
	GetBookingStatus(ctx context.Context) ([]dto.BookingStatusResponse, error)
	GetTodayFlights(ctx context.Context) ([]dto.TodayFlightResponse, error)
	GetRecentBookings(ctx context.Context, limit int) ([]dto.RecentBookingResponse, error)
}

type DashboardRepositoryImpl struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &DashboardRepositoryImpl{
		db: db,
	}
}

func (r *DashboardRepositoryImpl) GetSummary(
	ctx context.Context,
) (*dto.DashboardSummaryResponse, error) {

	var result dto.DashboardSummaryResponse

	if err := r.db.WithContext(ctx).
		Table("pnrs").
		Count(&result.TotalBookings).Error; err != nil {
		return nil, fmt.Errorf("GetSummary booking: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Table("pnr_passengers").
		Count(&result.TotalPassengers).Error; err != nil {
		return nil, fmt.Errorf("GetSummary passenger: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Table("flights").
		Where("DATE(departure_time) = CURRENT_DATE").
		Count(&result.TodayFlights).Error; err != nil {
		return nil, fmt.Errorf("GetSummary flights: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Table("payments").
		Select("COALESCE(SUM(amount),0)").
		Where("status = ?", "paid").
		Scan(&result.TotalRevenue).Error; err != nil {
		return nil, fmt.Errorf("GetSummary revenue: %w", err)
	}

	return &result, nil
}

func (r *DashboardRepositoryImpl) GetRevenueTrend(
	ctx context.Context,
	days int,
) ([]dto.RevenueTrendResponse, error) {

	var result []dto.RevenueTrendResponse

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			DATE(paid_at) as date,
			COALESCE(SUM(amount),0) as revenue
		FROM payments
		WHERE status = 'paid'
		AND paid_at >= NOW() - (? * INTERVAL '1 day')
		GROUP BY DATE(paid_at)
		ORDER BY DATE(paid_at)
	`, days).Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("GetRevenueTrend: %w", err)
	}

	return result, nil
}

func (r *DashboardRepositoryImpl) GetBookingStatus(
	ctx context.Context,
) ([]dto.BookingStatusResponse, error) {

	var result []dto.BookingStatusResponse

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			status,
			COUNT(*) as value
		FROM pnrs
		GROUP BY status
	`).Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("GetBookingStatus: %w", err)
	}

	return result, nil
}

func (r *DashboardRepositoryImpl) GetTodayFlights(
	ctx context.Context,
) ([]dto.TodayFlightResponse, error) {

	var result []dto.TodayFlightResponse

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			f.id,
			fs.flight_number,
			dep.code as origin,
			arr.code as destination,
			TO_CHAR(f.departure_time, 'HH24:MI') as departure_time,
			ac.model as aircraft,
			COUNT(DISTINCT pp.id) as passenger_count,
			f.status
		FROM flights f
		LEFT JOIN flight_schedules fs
			ON fs.id = f.schedule_id

		LEFT JOIN airports dep
			ON dep.id = fs.departure_airport_id

		LEFT JOIN airports arr
			ON arr.id = fs.arrival_airport_id

		LEFT JOIN aircrafts ac
			ON ac.id = f.aircraft_id

		LEFT JOIN pnr_segments ps
			ON ps.flight_id = f.id

		LEFT JOIN pnr_passengers pp
			ON pp.pnr_id = ps.pnr_id

		WHERE DATE(f.departure_time) = CURRENT_DATE

		GROUP BY
			f.id,
			fs.flight_number,
			dep.code,
			arr.code,
			ac.model

		ORDER BY f.departure_time ASC
	`).Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("GetTodayFlights: %w", err)
	}

	return result, nil
}

func (r *DashboardRepositoryImpl) GetRecentBookings(
	ctx context.Context,
	limit int,
) ([]dto.RecentBookingResponse, error) {

	var result []dto.RecentBookingResponse

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			p.id,
			p.record_locator as booking_code,

			COALESCE(
				MIN(pp.first_name || ' ' || pp.last_name),
				'Unknown Passenger'
			) as passenger_name,

			MIN(dep.code || ' → ' || arr.code) as route,

			COALESCE(pay.status, 'unpaid') as payment_status

		FROM pnrs p

		LEFT JOIN pnr_passengers pp
			ON pp.pnr_id = p.id

		LEFT JOIN pnr_segments ps
			ON ps.pnr_id = p.id

		LEFT JOIN flights f
			ON f.id = ps.flight_id

		LEFT JOIN flight_schedules fs
			ON fs.id = f.schedule_id

		LEFT JOIN airports dep
			ON dep.id = fs.departure_airport_id

		LEFT JOIN airports arr
			ON arr.id = fs.arrival_airport_id

		LEFT JOIN payments pay
			ON pay.pnr_id = p.id

		GROUP BY
			p.id,
			p.record_locator,
			pay.status,
			p.created_at

		ORDER BY p.created_at DESC
		LIMIT ?
	`, limit).Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("GetRecentBookings: %w", err)
	}

	return result, nil
}

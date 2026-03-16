package services

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/utils"
	"time"

	"github.com/google/uuid"
)

type FlightService interface {
	Search(ctx context.Context, req dto.SearchFlightRequest) ([]dto.FlightResult, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Flight, error)
	GetSeatMap(ctx context.Context, flightID uuid.UUID) ([]dto.FlightSeatResult, error)
	// Generate flights from a schedule for a date range
	GenerateFromSchedule(ctx context.Context, scheduleID uuid.UUID, from, to time.Time) (int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.FlightStatus) error
	GetByStatus(ctx context.Context, status models.FlightStatus) ([]models.Flight, error)
}

type flightService struct {
	flightRepo     repository.FlightRepository
	flightSeatRepo repository.FlightSeatRepository
	scheduleRepo   repository.FlightScheduleRepository
	airportRepo    repository.AirportRepository
	aircraftRepo   repository.AircraftRepository
	acSeatRepo     repository.AircraftSeatRepository
}

func NewFlightService(
	flightRepo repository.FlightRepository,
	flightSeatRepo repository.FlightSeatRepository,
	scheduleRepo repository.FlightScheduleRepository,
	airportRepo repository.AirportRepository,
	aircraftRepo repository.AircraftRepository,
	acSeatRepo repository.AircraftSeatRepository,
) FlightService {
	return &flightService{
		flightRepo:     flightRepo,
		flightSeatRepo: flightSeatRepo,
		scheduleRepo:   scheduleRepo,
		airportRepo:    airportRepo,
		aircraftRepo:   aircraftRepo,
		acSeatRepo:     acSeatRepo,
	}
}

// Search finds available flights on a given date for a route.
func (s *flightService) Search(ctx context.Context, req dto.SearchFlightRequest) ([]dto.FlightResult, error) {
	dep, err := s.airportRepo.FindByCode(ctx, req.DepartureCode)
	if err != nil {
		return nil, utils.ErrAirportNotFound
	}
	arr, err := s.airportRepo.FindByCode(ctx, req.ArrivalCode)
	if err != nil {
		return nil, utils.ErrAirportNotFound
	}

	flights, err := s.flightRepo.FindAvailable(ctx, dep.ID, arr.ID, req.Date)
	if err != nil {
		return nil, err
	}

	results := make([]dto.FlightResult, 0, len(flights))
	for _, f := range flights {
		available, _ := s.flightSeatRepo.CountAvailable(ctx, f.ID)

		// Get lowest available seat price
		seats, _ := s.flightSeatRepo.FindAvailableByFlight(ctx, f.ID)
		var lowestPrice float64
		for _, seat := range seats {
			if lowestPrice == 0 || seat.Price < lowestPrice {
				lowestPrice = seat.Price
			}
		}

		results = append(results, dto.FlightResult{
			Flight:         f,
			AvailableSeats: available,
			LowestPrice:    lowestPrice,
		})
	}
	return results, nil
}

func (s *flightService) GetByID(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	f, err := s.flightRepo.FindWithDetails(ctx, id)
	if err != nil {
		return nil, utils.ErrFlightNotFound
	}
	return f, nil
}

// GetSeatMap returns the full seat map of a flight enriched with class info.
func (s *flightService) GetSeatMap(ctx context.Context, flightID uuid.UUID) ([]dto.FlightSeatResult, error) {
	if _, err := s.flightRepo.FindByID(ctx, flightID); err != nil {
		return nil, utils.ErrFlightNotFound
	}

	seats, err := s.flightSeatRepo.FindByFlight(ctx, flightID)
	if err != nil {
		return nil, err
	}

	results := make([]dto.FlightSeatResult, 0, len(seats))
	for _, fs := range seats {
		r := dto.FlightSeatResult{FlightSeat: fs}
		if fs.AircraftSeat != nil {
			r.SeatNumber = fs.AircraftSeat.SeatNumber
			r.RowNumber = fs.AircraftSeat.RowNumber
			r.SeatLetter = fs.AircraftSeat.SeatLetter
			r.SeatType = fs.AircraftSeat.SeatType
			r.IsExitRow = fs.AircraftSeat.IsExitRow
			if fs.AircraftSeat.SeatClass != nil {
				r.ClassCode = fs.AircraftSeat.SeatClass.Code
				r.ClassName = fs.AircraftSeat.SeatClass.Name
			}
		}
		results = append(results, r)
	}
	return results, nil
}

// GenerateFromSchedule creates Flight instances for every day in [from, to]
// that matches the schedule's operating days, then bulk-creates FlightSeats.
func (s *flightService) GenerateFromSchedule(ctx context.Context, scheduleID uuid.UUID, from, to time.Time) (int, error) {
	sched, err := s.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return 0, utils.ErrScheduleNotFound
	}

	// Get aircraft attached to the schedule's airline — pick first available
	// In production you'd attach aircraft to schedule directly.
	aircrafts, err := s.aircraftRepo.FindAll(ctx)
	if err != nil || len(aircrafts) == 0 {
		return 0, utils.ErrAircraftNotFound
	}
	aircraft := aircrafts[0]

	// Parse operating days "1,2,3,4,5" → set of weekday numbers (1=Mon, 7=Sun)
	opDays := parseOperatingDays(sched.OperatingDays)

	depH, depM := parseHHMM(sched.DepartureTime)
	arrH, arrM := parseHHMM(sched.ArrivalTime)

	var flights []models.Flight
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		// Go's time.Weekday: Sunday=0..Saturday=6; convert to 1=Mon..7=Sun
		wd := int(d.Weekday())
		if wd == 0 {
			wd = 7
		}
		if !opDays[wd] {
			continue
		}

		depTime := time.Date(d.Year(), d.Month(), d.Day(), depH, depM, 0, 0, d.Location())
		arrTime := time.Date(d.Year(), d.Month(), d.Day(), arrH, arrM, 0, 0, d.Location())
		if arrTime.Before(depTime) {
			arrTime = arrTime.AddDate(0, 0, 1)
		}

		schedID := sched.ID
		acID := aircraft.ID
		flights = append(flights, models.Flight{
			ScheduleID:    &schedID,
			AircraftID:    &acID,
			DepartureTime: &depTime,
			ArrivalTime:   &arrTime,
			Status:        models.FlightStatusScheduled,
		})
	}

	if len(flights) == 0 {
		return 0, utils.ErrNoFlightsGenerated
	}

	if err := s.flightRepo.BulkCreate(ctx, flights); err != nil {
		return 0, fmt.Errorf("bulk create flights: %w", err)
	}

	// For each created flight, generate FlightSeats from aircraft_seats
	acSeats, err := s.acSeatRepo.FindByAircraftID(ctx, aircraft.ID)
	if err != nil {
		return len(flights), fmt.Errorf("load aircraft seats: %w", err)
	}

	for _, f := range flights {
		var flightSeats []models.FlightSeat
		for _, acSeat := range acSeats {
			price := priceByClass(acSeat.SeatClass)
			acSeatID := acSeat.ID
			flightID := f.ID
			flightSeats = append(flightSeats, models.FlightSeat{
				FlightID:       &flightID,
				AircraftSeatID: &acSeatID,
				Price:          price,
				Status:         models.FlightSeatAvailable,
			})
		}
		if err := s.flightSeatRepo.BulkCreate(ctx, flightSeats); err != nil {
			return len(flights), fmt.Errorf("create flight seats for flight %s: %w", f.ID, err)
		}
	}

	return len(flights), nil
}

func (s *flightService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.FlightStatus) error {
	if _, err := s.flightRepo.FindByID(ctx, id); err != nil {
		return utils.ErrFlightNotFound
	}
	return s.flightRepo.UpdateStatus(ctx, id, status)
}

func (s *flightService) GetByStatus(ctx context.Context, status models.FlightStatus) ([]models.Flight, error) {
	return s.flightRepo.FindByStatus(ctx, status)
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func parseOperatingDays(s string) map[int]bool {
	result := map[int]bool{}
	for _, ch := range s {
		if ch >= '1' && ch <= '7' {
			result[int(ch-'0')] = true
		}
	}
	return result
}

func parseHHMM(s string) (int, int) {
	var h, m, sec int
	fmt.Sscanf(s, "%d:%d:%d", &h, &m, &sec)
	return h, m
}

func priceByClass(sc *models.SeatClass) float64 {
	if sc == nil {
		return 500_000
	}
	switch sc.Code {
	case "F":
		return 5_000_000
	case "C":
		return 2_500_000
	default:
		return 850_000
	}
}

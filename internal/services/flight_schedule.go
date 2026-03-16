package services

import (
	"context"
	"fmt"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/utils"

	"github.com/google/uuid"
)

type FlightScheduleService interface {
	Create(ctx context.Context, req dto.CreateFlightScheduleRequest) (*models.FlightSchedule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.FlightSchedule, error)
	GetAll(ctx context.Context) ([]models.FlightSchedule, error)
	GetByRoute(ctx context.Context, depCode, arrCode string) ([]models.FlightSchedule, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateFlightScheduleRequest) (*models.FlightSchedule, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type flightScheduleService struct {
	scheduleRepo repository.FlightScheduleRepository
	airportRepo  repository.AirportRepository
}

func NewFlightScheduleService(
	scheduleRepo repository.FlightScheduleRepository,
	airportRepo repository.AirportRepository,
) FlightScheduleService {
	return &flightScheduleService{scheduleRepo: scheduleRepo, airportRepo: airportRepo}
}

func (s *flightScheduleService) Create(ctx context.Context, req dto.CreateFlightScheduleRequest) (*models.FlightSchedule, error) {
	// Validate flight number uniqueness
	if _, err := s.scheduleRepo.FindByFlightNumber(ctx, req.FlightNumber); err == nil {
		return nil, utils.ErrFlightNumberDuplicate
	}

	dep, err := s.airportRepo.FindByCode(ctx, req.DepartureAirportCode)
	if err != nil {
		return nil, fmt.Errorf("departure %w", utils.ErrAirportNotFound)
	}
	arr, err := s.airportRepo.FindByCode(ctx, req.ArrivalAirportCode)
	if err != nil {
		return nil, fmt.Errorf("arrival %w", utils.ErrAirportNotFound)
	}

	sched := &models.FlightSchedule{
		FlightNumber:       req.FlightNumber,
		DepartureAirportID: dep.ID,
		ArrivalAirportID:   arr.ID,
		DepartureTime:      req.DepartureTime,
		ArrivalTime:        req.ArrivalTime,
		OperatingDays:      req.OperatingDays,
	}
	if err := s.scheduleRepo.Create(ctx, sched); err != nil {
		return nil, err
	}
	return s.scheduleRepo.FindByID(ctx, sched.ID)
}

func (s *flightScheduleService) GetByID(ctx context.Context, id uuid.UUID) (*models.FlightSchedule, error) {
	sched, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.ErrScheduleNotFound
	}
	return sched, nil
}

func (s *flightScheduleService) GetAll(ctx context.Context) ([]models.FlightSchedule, error) {
	return s.scheduleRepo.FindAll(ctx)
}

func (s *flightScheduleService) GetByRoute(ctx context.Context, depCode, arrCode string) ([]models.FlightSchedule, error) {
	dep, err := s.airportRepo.FindByCode(ctx, depCode)
	if err != nil {
		return nil, utils.ErrAirportNotFound
	}
	arr, err := s.airportRepo.FindByCode(ctx, arrCode)
	if err != nil {
		return nil, utils.ErrAirportNotFound
	}
	return s.scheduleRepo.FindByRoute(ctx, dep.ID, arr.ID)
}

func (s *flightScheduleService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateFlightScheduleRequest) (*models.FlightSchedule, error) {
	sched, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.ErrScheduleNotFound
	}
	if req.DepartureTime != "" {
		sched.DepartureTime = req.DepartureTime
	}
	if req.ArrivalTime != "" {
		sched.ArrivalTime = req.ArrivalTime
	}
	if req.OperatingDays != "" {
		sched.OperatingDays = req.OperatingDays
	}
	if err := s.scheduleRepo.Update(ctx, sched); err != nil {
		return nil, err
	}
	return s.scheduleRepo.FindByID(ctx, id)
}

func (s *flightScheduleService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.scheduleRepo.FindByID(ctx, id); err != nil {
		return utils.ErrScheduleNotFound
	}
	return s.scheduleRepo.Delete(ctx, id)
}

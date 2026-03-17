package services

import (
	"context"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/utils"
	"time"

	"github.com/google/uuid"
)

type BaggageService interface {
	Add(ctx context.Context, req dto.AddBaggageRequest) (*models.Baggage, error)
	UpdateStatus(ctx context.Context, baggageID uuid.UUID, status models.BaggageStatus) (*models.Baggage, error)
	GetByPassenger(ctx context.Context, passengerID uuid.UUID) ([]models.Baggage, error)
	GetBySegment(ctx context.Context, segmentID uuid.UUID) ([]models.Baggage, error)
}

type baggageService struct {
	baggageRepo   repository.BaggageRepository
	checkinRepo   repository.CheckinRepository
	passengerRepo repository.PNRPassengerRepository
}

func NewBaggageService(
	baggageRepo repository.BaggageRepository,
	checkinRepo repository.CheckinRepository,
	passengerRepo repository.PNRPassengerRepository,
) BaggageService {
	return &baggageService{
		baggageRepo:   baggageRepo,
		checkinRepo:   checkinRepo,
		passengerRepo: passengerRepo,
	}
}

// Add registers checked baggage for a passenger.
// Passenger must have checked in first.
func (s *baggageService) Add(ctx context.Context, req dto.AddBaggageRequest) (*models.Baggage, error) {
	if _, err := s.passengerRepo.FindByID(ctx, req.PassengerID); err != nil {
		return nil, utils.ErrPassengerNotFound
	}

	checkedIn, err := s.checkinRepo.IsCheckedIn(ctx, req.PassengerID, req.SegmentID)
	if err != nil {
		return nil, err
	}
	if !checkedIn {
		return nil, utils.ErrCheckinRequiredBaggage
	}

	baggage := &models.Baggage{
		PassengerID: &req.PassengerID,
		SegmentID:   &req.SegmentID,
		Weight:      req.Weight,
		TagNumber:   generateBaggageTag(),
		Status:      models.BaggageCheckedIn,
	}
	if err := s.baggageRepo.Create(ctx, baggage); err != nil {
		return nil, err
	}
	return baggage, nil
}

func (s *baggageService) UpdateStatus(ctx context.Context, baggageID uuid.UUID, status models.BaggageStatus) (*models.Baggage, error) {
	baggage, err := s.baggageRepo.FindByID(ctx, baggageID)
	if err != nil {
		return nil, utils.ErrBaggageNotFound
	}
	if err := s.baggageRepo.UpdateStatus(ctx, baggageID, status); err != nil {
		return nil, err
	}
	baggage.Status = status
	return baggage, nil
}

func (s *baggageService) GetByPassenger(ctx context.Context, passengerID uuid.UUID) ([]models.Baggage, error) {
	return s.baggageRepo.FindByPassengerID(ctx, passengerID)
}

func (s *baggageService) GetBySegment(ctx context.Context, segmentID uuid.UUID) ([]models.Baggage, error) {
	return s.baggageRepo.FindBySegmentID(ctx, segmentID)
}

func generateBaggageTag() string {
	const digits = "0123456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = digits[time.Now().UnixNano()%10]
		time.Sleep(1 * time.Nanosecond)
	}
	return "BG" + string(b)
}

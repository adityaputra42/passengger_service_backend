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

type BoardingPassService interface {
	Issue(ctx context.Context, req dto.IssueBoardingPassRequest) (*models.BoardingPass, error)
	GetByPassengerAndSegment(ctx context.Context, passengerID, segmentID uuid.UUID) (*models.BoardingPass, error)
	GetBySegment(ctx context.Context, segmentID uuid.UUID) ([]models.BoardingPass, error)
}

type boardingPassService struct {
	bpRepo         repository.BoardingPassRepository
	checkinRepo    repository.CheckinRepository
	passengerRepo  repository.PNRPassengerRepository
	segmentRepo    repository.PNRSegmentRepository
	seatAssignRepo repository.SeatAssignmentRepository
}

func NewBoardingPassService(
	bpRepo repository.BoardingPassRepository,
	checkinRepo repository.CheckinRepository,
	passengerRepo repository.PNRPassengerRepository,
	segmentRepo repository.PNRSegmentRepository,
	seatAssignRepo repository.SeatAssignmentRepository,
) BoardingPassService {
	return &boardingPassService{
		bpRepo:         bpRepo,
		checkinRepo:    checkinRepo,
		passengerRepo:  passengerRepo,
		segmentRepo:    segmentRepo,
		seatAssignRepo: seatAssignRepo,
	}
}

func (s *boardingPassService) Issue(ctx context.Context, req dto.IssueBoardingPassRequest) (*models.BoardingPass, error) {
	passenger, err := s.passengerRepo.FindByID(ctx, req.PassengerID)
	if err != nil {
		return nil, utils.ErrPassengerNotFound
	}

	if _, err := s.segmentRepo.FindByID(ctx, req.SegmentID); err != nil {
		return nil, fmt.Errorf("segment tidak ditemukan")
	}

	checkedIn, err := s.checkinRepo.IsCheckedIn(ctx, req.PassengerID, req.SegmentID)
	if err != nil {
		return nil, err
	}
	if !checkedIn {
		return nil, utils.ErrCheckinRequiredBoarding
	}

	if _, err := s.bpRepo.FindByPassengerAndSegment(ctx, req.PassengerID, req.SegmentID); err == nil {
		return nil, utils.ErrBoardingPassExists
	}

	seatNumber := "UNASSIGNED"
	if assign, err := s.seatAssignRepo.FindByPassengerID(ctx, req.PassengerID); err == nil {
		if assign.FlightSeat != nil && assign.FlightSeat.AircraftSeat != nil {
			seatNumber = assign.FlightSeat.AircraftSeat.SeatNumber
		}
	}

	qrContent := fmt.Sprintf("%s|%s|%s|%s|%s",
		passenger.ID, req.SegmentID, seatNumber, req.BoardingGroup, req.Gate,
	)
	qrCode := generateQRCode(qrContent)

	boardingTime := req.BoardingTime
	bp := &models.BoardingPass{
		PassengerID:   &req.PassengerID,
		SegmentID:     &req.SegmentID,
		BoardingGroup: req.BoardingGroup,
		Gate:          req.Gate,
		BoardingTime:  &boardingTime,
		QRCode:        qrCode,
	}
	if err := s.bpRepo.Create(ctx, bp); err != nil {
		return nil, fmt.Errorf("create boarding pass: %w", err)
	}

	return s.bpRepo.FindByPassengerAndSegment(ctx, req.PassengerID, req.SegmentID)
}

func (s *boardingPassService) GetByPassengerAndSegment(ctx context.Context, passengerID, segmentID uuid.UUID) (*models.BoardingPass, error) {
	bp, err := s.bpRepo.FindByPassengerAndSegment(ctx, passengerID, segmentID)
	if err != nil {
		return nil, utils.ErrBoardingPassNotFound
	}
	return bp, nil
}

func (s *boardingPassService) GetBySegment(ctx context.Context, segmentID uuid.UUID) ([]models.BoardingPass, error) {
	return s.bpRepo.FindBySegmentID(ctx, segmentID)
}

func generateQRCode(content string) string {
	return fmt.Sprintf("QR::%x", []byte(content))
}


package services

import (
	"context"
	"errors"
	"fmt"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AircraftService interface {
	Create(ctx context.Context, req dto.CreateAircraftRequest) (*models.Aircraft, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Aircraft, error)
	GetAll(ctx context.Context) ([]models.Aircraft, error)
	GetWithSeats(ctx context.Context, id uuid.UUID) (*models.Aircraft, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateAircraftRequest) (*models.Aircraft, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GenerateSeats(ctx context.Context, aircraftID uuid.UUID, req dto.GenerateSeatsRequest) ([]models.AircraftSeat, error)
}

type aircraftService struct {
	aircraftRepo     repository.AircraftRepository
	aircraftSeatRepo repository.AircraftSeatRepository
	seatClassRepo    repository.SeatClassRepository
}

func NewAircraftService(
	aircraftRepo repository.AircraftRepository,
	aircraftSeatRepo repository.AircraftSeatRepository,
	seatClassRepo repository.SeatClassRepository,
) AircraftService {
	return &aircraftService{
		aircraftRepo:     aircraftRepo,
		aircraftSeatRepo: aircraftSeatRepo,
		seatClassRepo:    seatClassRepo,
	}
}

func (s *aircraftService) Create(ctx context.Context, req dto.CreateAircraftRequest) (*models.Aircraft, error) {
	aircraft := &models.Aircraft{
		Model:        req.Model,
		Manufacturer: req.Manufacturer,
	}
	if err := s.aircraftRepo.Create(ctx, aircraft); err != nil {
		return nil, err
	}
	return aircraft, nil
}

func (s *aircraftService) GetByID(ctx context.Context, id uuid.UUID) (*models.Aircraft, error) {
	a, err := s.aircraftRepo.FindByID(ctx, nil, id)
	if err != nil {
		return nil, utils.ErrAircraftNotFound
	}
	return a, nil
}

func (s *aircraftService) GetAll(ctx context.Context) ([]models.Aircraft, error) {
	return s.aircraftRepo.FindAll(ctx)
}

func (s *aircraftService) GetWithSeats(ctx context.Context, id uuid.UUID) (*models.Aircraft, error) {
	a, err := s.aircraftRepo.FindWithSeats(ctx, id)
	if err != nil {
		return nil, utils.ErrAircraftNotFound
	}
	return a, nil
}

func (s *aircraftService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateAircraftRequest) (*models.Aircraft, error) {
	aircraft, err := s.aircraftRepo.FindByID(ctx, nil, id)
	if err != nil {
		return nil, utils.ErrAircraftNotFound
	}
	if req.Model != "" {
		aircraft.Model = req.Model
	}
	if req.Manufacturer != "" {
		aircraft.Manufacturer = req.Manufacturer
	}
	if err := s.aircraftRepo.Update(ctx, nil, aircraft); err != nil {
		return nil, err
	}
	return aircraft, nil
}

func (s *aircraftService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.aircraftRepo.FindByID(ctx, nil, id); err != nil {
		return utils.ErrAircraftNotFound
	}
	return s.aircraftRepo.Delete(ctx, id)
}

// GenerateSeats bulk-creates AircraftSeats from a layout config.
func (s *aircraftService) GenerateSeats(ctx context.Context, aircraftID uuid.UUID, req dto.GenerateSeatsRequest) ([]models.AircraftSeat, error) {

	var seats []models.AircraftSeat
	exitSet := map[int]bool{}

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		aircraft, err := s.aircraftRepo.FindByID(ctx, tx, aircraftID)
		if err != nil {
			return utils.ErrAircraftNotFound
		}

		classMap := map[string]*models.SeatClass{}
		for _, cfg := range req.Classes {
			sc, err := s.seatClassRepo.FindByCode(ctx, tx, cfg.ClassCode)
			if err != nil {
				return fmt.Errorf("seat class %q not found: %w", cfg.ClassCode, err)
			}
			classMap[cfg.ClassCode] = sc
		}

		row := 1
		for _, cfg := range req.Classes {
			sc := classMap[cfg.ClassCode]
			for _, er := range cfg.ExitRowNums {
				exitSet[er] = true
			}
			for r := 0; r < cfg.Rows; r++ {
				for i, letter := range cfg.Letters {
					seatType := determineSeatType(letter, cfg.Letters)
					classID := sc.ID
					seat := models.AircraftSeat{
						AircraftID:  aircraftID,
						SeatNumber:  fmt.Sprintf("%d%s", row, letter),
						RowNumber:   row,
						SeatLetter:  letter,
						XPosition:   row,
						YPosition:   i,
						SeatClassID: &classID,
						SeatType:    seatType,
						IsExitRow:   exitSet[row],
					}
					seats = append(seats, seat)
				}
				row++
			}
		}

		if err := s.aircraftSeatRepo.BulkCreate(ctx, tx, seats); err != nil {
			return err
		}

		aircraft.TotalSeats = len(seats)
		if err = s.aircraftRepo.Update(ctx, tx, aircraft); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return seats, nil
}

func determineSeatType(letter string, letters []string) string {
	n := len(letters)
	if n == 0 {
		return "seat"
	}
	first, last := letters[0], letters[n-1]
	if letter == first || letter == last {
		return "window"
	}
	mid := n / 2
	if letter == letters[mid-1] || letter == letters[mid] {
		return "aisle"
	}
	return "middle"
}

var _ = errors.Is(gorm.ErrRecordNotFound, gorm.ErrRecordNotFound)

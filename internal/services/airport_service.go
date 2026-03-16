package services

import (
	"context"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/utils"
	"strings"

	"github.com/google/uuid"
)

type AirportService interface {
	Create(ctx context.Context, req dto.CreateAirportRequest) (*models.Airport, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Airport, error)
	GetByCode(ctx context.Context, code string) (*models.Airport, error)
	GetAll(ctx context.Context) ([]models.Airport, error)
	Search(ctx context.Context, query string) ([]models.Airport, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateAirportRequest) (*models.Airport, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type airportService struct {
	repo repository.AirportRepository
}

func NewAirportService(repo repository.AirportRepository) AirportService {
	return &airportService{repo: repo}
}

func (s *airportService) Create(ctx context.Context, req dto.CreateAirportRequest) (*models.Airport, error) {
	code := strings.ToUpper(req.Code)
	if _, err := s.repo.FindByCode(ctx, code); err == nil {
		return nil, utils.ErrAirportCodeDuplicate
	}
	airport := &models.Airport{
		Code:     code,
		Name:     req.Name,
		City:     req.City,
		Country:  req.Country,
		Timezone: req.Timezone,
	}
	if err := s.repo.Create(ctx, airport); err != nil {
		return nil, err
	}
	return airport, nil
}

func (s *airportService) GetByID(ctx context.Context, id uuid.UUID) (*models.Airport, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.ErrAirportNotFound
	}
	return a, nil
}

func (s *airportService) GetByCode(ctx context.Context, code string) (*models.Airport, error) {
	a, err := s.repo.FindByCode(ctx, strings.ToUpper(code))
	if err != nil {
		return nil, utils.ErrAirportNotFound
	}
	return a, nil
}

func (s *airportService) GetAll(ctx context.Context) ([]models.Airport, error) {
	return s.repo.FindAll(ctx)
}

func (s *airportService) Search(ctx context.Context, query string) ([]models.Airport, error) {
	return s.repo.Search(ctx, query)
}

func (s *airportService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateAirportRequest) (*models.Airport, error) {
	airport, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.ErrAirportNotFound
	}
	if req.Name != "" {
		airport.Name = req.Name
	}
	if req.City != "" {
		airport.City = req.City
	}
	if req.Country != "" {
		airport.Country = req.Country
	}
	if req.Timezone != "" {
		airport.Timezone = req.Timezone
	}
	if err := s.repo.Update(ctx, airport); err != nil {
		return nil, err
	}
	return airport, nil
}

func (s *airportService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return utils.ErrAirportNotFound
	}
	return s.repo.Delete(ctx, id)
}

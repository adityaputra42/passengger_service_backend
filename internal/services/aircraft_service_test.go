package services_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/mocks"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
)

func TestAircraftService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAircraftRepo := mocks.NewMockAircraftRepository(ctrl)
	mockAircraftSeatRepo := mocks.NewMockAircraftSeatRepository(ctrl)
	mockSeatClassRepo := mocks.NewMockSeatClassRepository(ctrl)

	svc := services.NewAircraftService(mockAircraftRepo, mockAircraftSeatRepo, mockSeatClassRepo)

	ctx := context.Background()
	req := dto.CreateAircraftRequest{
		Model:        "Boeing 737",
		Manufacturer: "Boeing",
	}

	t.Run("Success", func(t *testing.T) {
		mockAircraftRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		a, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, a)
		assert.Equal(t, req.Model, a.Model)
	})
}

func TestAircraftService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAircraftRepo := mocks.NewMockAircraftRepository(ctrl)
	mockAircraftSeatRepo := mocks.NewMockAircraftSeatRepository(ctrl)
	mockSeatClassRepo := mocks.NewMockSeatClassRepository(ctrl)

	svc := services.NewAircraftService(mockAircraftRepo, mockAircraftSeatRepo, mockSeatClassRepo)

	ctx := context.Background()
	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockAircraftRepo.EXPECT().FindByID(ctx, nil, id).Return(&models.Aircraft{ID: id, Model: "Boeing 737"}, nil)

		a, err := svc.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, a)
		assert.Equal(t, id, a.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockAircraftRepo.EXPECT().FindByID(ctx, nil, id).Return(nil, utils.ErrAircraftNotFound)

		a, err := svc.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, utils.ErrAircraftNotFound, err)
		assert.Nil(t, a)
	})
}

func TestAircraftService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAircraftRepo := mocks.NewMockAircraftRepository(ctrl)
	svc := services.NewAircraftService(mockAircraftRepo, nil, nil)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockAircraftRepo.EXPECT().FindAll(ctx).Return([]models.Aircraft{{Model: "B737"}, {Model: "A320"}}, nil)

		aircrafts, err := svc.GetAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, aircrafts, 2)
	})
}

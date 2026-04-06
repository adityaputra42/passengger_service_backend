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

func TestAirportService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAirportRepository(ctrl)
	svc := services.NewAirportService(mockRepo)

	ctx := context.Background()
	req := dto.CreateAirportRequest{
		Code:     "CGK",
		Name:     "Soekarno-Hatta",
		City:     "Jakarta",
		Country:  "Indonesia",
		Timezone: "Asia/Jakarta",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().FindByCode(ctx, "CGK").Return(nil, utils.ErrAirportNotFound)
		mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		res, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "CGK", res.Code)
		assert.Equal(t, "Soekarno-Hatta", res.Name)
	})

	t.Run("Duplicate Code", func(t *testing.T) {
		existingAirport := &models.Airport{Code: "CGK"}
		mockRepo.EXPECT().FindByCode(ctx, "CGK").Return(existingAirport, nil)

		res, err := svc.Create(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, utils.ErrAirportCodeDuplicate, err)
		assert.Nil(t, res)
	})
}

func TestAirportService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAirportRepository(ctrl)
	svc := services.NewAirportService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	expectedAirport := &models.Airport{ID: id, Code: "CGK"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, id).Return(expectedAirport, nil)

		res, err := svc.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, id, res.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, id).Return(nil, utils.ErrAirportNotFound)

		res, err := svc.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, utils.ErrAirportNotFound, err)
		assert.Nil(t, res)
	})
}

func TestAirportService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAirportRepository(ctrl)
	svc := services.NewAirportService(mockRepo)

	ctx := context.Background()
	airports := []models.Airport{
		{Code: "CGK"},
		{Code: "DPS"},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().FindAll(ctx).Return(airports, nil)

		res, err := svc.GetAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})
}

func TestAirportService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAirportRepository(ctrl)
	svc := services.NewAirportService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	existingAirport := &models.Airport{
		ID:       id,
		Code:     "CGK",
		Name:     "Old Name",
		City:     "Jakarta",
		Country:  "Indonesia",
		Timezone: "Asia/Jakarta",
	}

	req := dto.UpdateAirportRequest{
		Name: "New Name",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, id).Return(existingAirport, nil)
		mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		res, err := svc.Update(ctx, id, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "New Name", res.Name) // Only name should be updated
		assert.Equal(t, "Jakarta", res.City)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, id).Return(nil, utils.ErrAirportNotFound)

		res, err := svc.Update(ctx, id, req)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestAirportService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAirportRepository(ctrl)
	svc := services.NewAirportService(mockRepo)

	ctx := context.Background()
	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, id).Return(&models.Airport{ID: id}, nil)
		mockRepo.EXPECT().Delete(ctx, id).Return(nil)

		err := svc.Delete(ctx, id)

		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, id).Return(nil, utils.ErrAirportNotFound)

		err := svc.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, utils.ErrAirportNotFound, err)
	})
}

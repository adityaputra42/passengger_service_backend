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

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	svc := services.NewUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	req := dto.CreateUserRequest{
		Email:    "test@example.com",
		FullName: "Test User",
		Password: "password",
		RoleID:   1,
	}

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.EXPECT().FindByEmail(ctx, req.Email).Return(models.User{}, utils.ErrUserNotFound)
		mockRoleRepo.EXPECT().FindById(ctx, req.RoleID).Return(models.Role{ID: 1}, nil)
		mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(models.User{Email: req.Email}, nil)

		user, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Email, user.Email)
	})

	t.Run("Duplicate Email", func(t *testing.T) {
		mockUserRepo.EXPECT().FindByEmail(ctx, req.Email).Return(models.User{Email: req.Email}, nil)

		user, err := svc.Create(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, utils.ErrEmailAlreadyExists, err)
		assert.Nil(t, user)
	})
}

func TestUserService_GetByUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	svc := services.NewUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	uid := uuid.New()

	t.Run("Success", func(t *testing.T) {
		expectedUser := models.User{UID: uid, Email: "test@example.com"}
		mockUserRepo.EXPECT().FindByUid(ctx, uid).Return(expectedUser, nil)

		user, err := svc.GetByUID(ctx, uid)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uid, user.UID)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUserRepo.EXPECT().FindByUid(ctx, uid).Return(models.User{}, utils.ErrUserNotFound)

		user, err := svc.GetByUID(ctx, uid)

		assert.Error(t, err)
		assert.Equal(t, utils.ErrUserNotFound, err)
		assert.Nil(t, user)
	})
}

// Add more tests for Update, Delete as needed following this pattern

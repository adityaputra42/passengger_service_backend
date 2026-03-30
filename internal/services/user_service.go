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

type UserService interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (*models.User, error)
	GetByUID(ctx context.Context, uid uuid.UUID) (*models.User, error)
	GetAll(ctx context.Context, page, limit int) (*dto.UserListResponse, error)
	Update(ctx context.Context, uid uuid.UUID, req dto.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, uid uuid.UUID) error
	UpdateProfile(ctx context.Context, uid uuid.UUID, req dto.UpdateProfileRequest) (*models.User, error)
}

// ─────────────────────────────────────────────
// User Service
// ─────────────────────────────────────────────

type userService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) UserService {
	return &userService{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (*models.User, error) {
	// Check email uniqueness
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, utils.ErrEmailAlreadyExists
	}

	// Validate role exists
	if _, err := s.roleRepo.FindById(ctx, req.RoleID); err != nil {
		return nil, utils.ErrRoleNotFound
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Email:        req.Email,
		FullName:     req.FullName,
		PasswordHash: string(hash),
		RoleID:       req.RoleID,
	}
	result, err := s.userRepo.Create(ctx, *user)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *userService) GetByUID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}
	return &user, nil
}

func (s *userService) GetAll(ctx context.Context, page, limit int) (*dto.UserListResponse, error) {
	param := dto.UserListRequest{}
	users, err := s.userRepo.FindAll(ctx, param)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}
	return users, nil

}

func (s *userService) Update(ctx context.Context, uid uuid.UUID, req dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.RoleID != 0 {
		if _, err := s.roleRepo.FindById(ctx, req.RoleID); err != nil {
			return nil, utils.ErrRoleNotFound
		}
		user.RoleID = req.RoleID
	}

	userUpdate, err := s.userRepo.Update(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &userUpdate, nil
}

func (s *userService) Delete(ctx context.Context, uid uuid.UUID) error {
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return utils.ErrUserNotFound
	}
	return s.userRepo.Delete(ctx, user)
}

func (s *userService) UpdateProfile(ctx context.Context, uid uuid.UUID, req dto.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	userUpdate, err := s.userRepo.Update(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &userUpdate, nil
}

package services

import (
	"context"
	"errors"
	"fmt"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	LoginAdmin(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	Logout(ctx context.Context, userUID uuid.UUID) error
	Me(ctx context.Context, userUID uuid.UUID) (*models.User, error)
	ChangePassword(ctx context.Context, userUID uuid.UUID, req dto.ChangePasswordRequest) error
	generateTokenResponse(user *models.User) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	jwtService *utils.JWTService
}

// LoginAdmin implements [AuthService].
func (a *authService) LoginAdmin(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	userResult, err := a.userRepo.FindByUsernameOrEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}
	if err := utils.CheckPassword(req.Password, userResult.PasswordHash); err != nil {
		return nil, errors.New("invalid credentials")
	}
	role, err := a.roleRepo.FindById(userResult.RoleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	if role.Level < 3 {
		return nil, errors.New("unauthorized: admin only")
	}
	return a.generateTokenResponse(&userResult)
}

// generateTokenResponse implements [AuthService].
func (a *authService) generateTokenResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, expiresAt, err := a.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := a.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         user,
	}, nil
}

// ChangePassword implements [AuthService].
func (a *authService) ChangePassword(ctx context.Context, userUID uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := a.userRepo.FindByUid(userUID)
	if err != nil {
		return utils.ErrUserNotFound
	}
	newPass, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	user.PasswordHash = string(newPass)
	_, err = a.userRepo.Update(&user)
	return err
}

// Login implements [AuthService].
func (a *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	userResult, err := a.userRepo.FindByUsernameOrEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}
	if err := utils.CheckPassword(req.Password, userResult.PasswordHash); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return a.generateTokenResponse(&userResult)

}

// Logout implements [AuthService].
func (a *authService) Logout(ctx context.Context, userUID uuid.UUID) error {
	panic("unimplemented")
}

// Me implements [AuthService].
func (a *authService) Me(ctx context.Context, userUID uuid.UUID) (*models.User, error) {
	user, err := a.userRepo.FindByUid(userUID)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}
	return &user, nil
}

// RefreshToken implements [AuthService].
func (a *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	claims, err := a.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := a.userRepo.FindByUid(claims.UID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return a.generateTokenResponse(&user)
}

func NewAuthService(userRepo repository.UserRepository,
	jwtService *utils.JWTService) AuthService {
	return &authService{userRepo: userRepo, jwtService: jwtService}
}

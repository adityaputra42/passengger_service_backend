package dto

import (
	"passenger_service_backend/internal/models"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	User         *models.User `json:"user"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type AuthResponseDTO struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresAt    time.Time     `json:"expires_at"`
	User         *UserResponse `json:"user"`
}

func ToAuthResponseDTO(a *AuthResponse) *AuthResponseDTO {
	if a == nil {
		return nil
	}
	return &AuthResponseDTO{
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		ExpiresAt:    a.ExpiresAt,
		User:         ToUserResponse(a.User),
	}
}

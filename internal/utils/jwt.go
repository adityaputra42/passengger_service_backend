package utils

import (
	"errors"
	"log"
	"passenger_service_backend/internal/config"
	"passenger_service_backend/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UID    uuid.UUID `json:"uid"`
	Email  string    `json:"email"`
	RoleID uint      `json:"role_id"`
	Type   string    `json:"type"` // "access" or "refresh"
}

type JWTService struct {
	config *config.Config
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{config: cfg}
}

func (s *JWTService) GenerateAccessToken(user *models.User) (string, time.Time, error) {
	// 🔍 DEBUG: Validasi user sebelum generate token
	if user == nil {
		log.Println("❌ JWT ERROR: user is nil")
		return "", time.Time{}, errors.New("user cannot be nil")
	}

	expiresAt := time.Now().Add(s.config.JWT.AccessTokenExpiry)

	claims := jwt.MapClaims{
		"uid":     user.UID.String(),
		"email":   user.Email,
		"role_id": user.RoleID,
		"type":    "access",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		log.Printf("❌ JWT ERROR: Failed to sign token: %v", err)
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *JWTService) GenerateRefreshToken(user *models.User) (string, time.Time, error) {
	// 🔍 DEBUG: Validasi user sebelum generate token
	if user == nil {
		log.Println("❌ JWT ERROR: user is nil")
		return "", time.Time{}, errors.New("user cannot be nil")
	}

	expiresAt := time.Now().Add(s.config.JWT.RefreshTokenExpiry)

	claims := jwt.MapClaims{
		"uid":     user.UID.String(),
		"email":   user.Email,
		"role_id": user.RoleID,
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	// 🔍 DEBUG: Log claims
	log.Printf("✅ Generating Refresh Token - UID: %d, Email: %s", user.UID, user.Email)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.RefreshSecret))
	if err != nil {
		log.Printf("❌ JWT ERROR: Failed to sign refresh token: %v", err)
		return "", time.Time{}, err
	}

	log.Printf("✅ Refresh Token Generated Successfully for UserID: %d", user.UID)
	return tokenString, expiresAt, nil
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		log.Printf("❌ JWT VALIDATION ERROR: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			log.Println("❌ JWT ERROR: Invalid token type")
			return nil, errors.New("invalid token type")
		}
		uidStr, ok := claims["uid"].(string)
		if !ok {
			log.Printf("❌ JWT ERROR: Invalid uid claim, got: %v (type: %T)", claims["uid"], claims["uid"])
			return nil, errors.New("invalid uid claim")
		}

		uid, err := uuid.Parse(uidStr)
		if err != nil {
			log.Printf("❌ JWT ERROR: Failed to parse UUID: %v", err)
			return nil, err
		}

		email, ok := claims["email"].(string)
		if !ok {
			log.Println("❌ JWT ERROR: Invalid email claim")
			return nil, errors.New("invalid email claim")
		}

		roleID, ok := claims["role_id"].(float64)
		if !ok {
			log.Println("❌ JWT ERROR: Invalid role_id claim")
			return nil, errors.New("invalid role_id claim")
		}

		jwtClaims := &JWTClaims{
			UID:    uid,
			Email:  email,
			RoleID: uint(roleID),
			Type:   tokenType,
		}

		// 🔍 DEBUG: Log parsed claims
		log.Printf("✅ Token Validated - uid: %d, Email: %s, RoleID: %d",
			jwtClaims.UID, jwtClaims.Email, jwtClaims.RoleID)

		return jwtClaims, nil
	}

	log.Println("❌ JWT ERROR: Invalid token")
	return nil, errors.New("invalid token")
}

func (s *JWTService) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWT.RefreshSecret), nil
	})

	if err != nil {
		log.Printf("❌ REFRESH TOKEN VALIDATION ERROR: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "refresh" {
			log.Println("❌ JWT ERROR: Invalid refresh token type")
			return nil, errors.New("invalid token type")
		}

		uidStr, ok := claims["uid"].(string)
		if !ok {
			log.Printf("❌ JWT ERROR: Invalid uid claim, got: %v (type: %T)", claims["uid"], claims["uid"])
			return nil, errors.New("invalid uid claim")
		}

		uid, err := uuid.Parse(uidStr)
		if err != nil {
			log.Printf("❌ JWT ERROR: Failed to parse UUID: %v", err)
			return nil, err
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, errors.New("invalid email claim")
		}

		roleID, ok := claims["role_id"].(float64)
		if !ok {
			return nil, errors.New("invalid role_id claim")
		}

		jwtClaims := &JWTClaims{
			UID:    uid,
			Email:  email,
			RoleID: uint(roleID),
			Type:   tokenType,
		}

		// 🔍 DEBUG: Log parsed refresh token claims
		log.Printf("✅ Refresh Token Validated - UID: %d, Email: %s",
			jwtClaims.UID, jwtClaims.Email)

		return jwtClaims, nil
	}

	return nil, errors.New("invalid token")
}

// Helper function untuk debug token tanpa validasi signature
func (s *JWTService) DebugToken(tokenString string) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Printf("❌ DEBUG TOKEN ERROR: %v", err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		log.Println("🔍 TOKEN DEBUG INFO:")
		log.Printf("   uid: %v (type: %T)", claims["uid"], claims["uid"])
		log.Printf("   email: %v", claims["email"])
		log.Printf("   role_id: %v", claims["role_id"])
		log.Printf("   type: %v", claims["type"])
		log.Printf("   exp: %v", claims["exp"])
		log.Printf("   iat: %v", claims["iat"])
	}
}

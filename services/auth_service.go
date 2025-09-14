package services

import (
	"context"
	"errors"
	"time"

	"go-boilerplate/config"
	"go-boilerplate/logger"
	"go-boilerplate/services/interfaces"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	jwtSecret string
}

func NewAuthService(cfg *config.Config) interfaces.AuthService {
	return &authService{
		jwtSecret: cfg.JWT.Secret,
	}
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *authService) GenerateToken(userID uint) (string, error) {
	ctx := context.Background()
	logger.Debug(ctx, "AuthService.GenerateToken start", map[string]any{"user_id": userID})
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		logger.Error(ctx, "AuthService.GenerateToken failed", map[string]any{"user_id": userID, "error": err.Error()})
		return "", err
	}
	logger.Info(ctx, "AuthService.GenerateToken success", map[string]any{"user_id": userID})
	return signed, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	ctx := context.Background()
	logger.Debug(ctx, "AuthService.ValidateToken start", nil)
	tok, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		logger.Warn(ctx, "AuthService.ValidateToken failed", map[string]any{"error": err.Error()})
		return nil, err
	}
	logger.Debug(ctx, "AuthService.ValidateToken success", nil)
	return tok, nil
}

func (s *authService) GetUserIDFromToken(token *jwt.Token) (uint, error) {
	ctx := context.Background()
	logger.Debug(ctx, "AuthService.GetUserIDFromToken start", nil)
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		logger.Warn(ctx, "AuthService.GetUserIDFromToken invalid token", nil)
		return 0, errors.New("invalid token")
	}
	logger.Debug(ctx, "AuthService.GetUserIDFromToken success", map[string]any{"user_id": claims.UserID})
	return claims.UserID, nil
}

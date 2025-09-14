package interfaces

import "github.com/golang-jwt/jwt/v5"

type AuthService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserIDFromToken(token *jwt.Token) (uint, error)
}

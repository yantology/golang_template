package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	secretKey []byte
}

func NewAuthService(secretKey string) *AuthService {
	return &AuthService{
		secretKey: []byte(secretKey),
	}
}

func (a *AuthService) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})
	return token.SignedString(a.secretKey)
}
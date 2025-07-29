package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token has expired")
	ErrTokenNotFound  = errors.New("token not found")
	ErrInvalidClaims  = errors.New("invalid token claims")
	ErrTokenBlacklist = errors.New("token is blacklisted")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	TokenType TokenType `json:"token_type"`
	SessionID uuid.UUID `json:"session_id"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret             []byte
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
	issuer             string
	audience           string
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration, issuer, audience string) *JWTManager {
	return &JWTManager{
		secret:             []byte(secret),
		accessTokenTTL:     accessTTL,
		refreshTokenTTL:    refreshTTL,
		issuer:             issuer,
		audience:           audience,
	}
}

func (j *JWTManager) GenerateTokenPair(userID uuid.UUID, email string, sessionID uuid.UUID) (*TokenPair, error) {
	accessToken, err := j.generateToken(userID, email, sessionID, AccessToken, j.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.generateToken(userID, email, sessionID, RefreshToken, j.refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(j.accessTokenTTL).Unix(),
	}, nil
}

func (j *JWTManager) generateToken(userID uuid.UUID, email string, sessionID uuid.UUID, tokenType TokenType, ttl time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)

	claims := Claims{
		UserID:    userID,
		Email:     email,
		TokenType: tokenType,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Audience:  jwt.ClaimStrings{j.audience},
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	// Validate audience and issuer
	if len(claims.Audience) == 0 || claims.Audience[0] != j.audience {
		return nil, ErrInvalidClaims
	}

	if claims.Issuer != j.issuer {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

func (j *JWTManager) RefreshAccessToken(refreshTokenString string, sessionID uuid.UUID) (*TokenPair, error) {
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != RefreshToken {
		return nil, ErrInvalidToken
	}

	if claims.SessionID != sessionID {
		return nil, ErrInvalidToken
	}

	return j.GenerateTokenPair(claims.UserID, claims.Email, sessionID)
}

func (j *JWTManager) ExtractUserID(tokenString string) (uuid.UUID, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

func (j *JWTManager) ExtractSessionID(tokenString string) (uuid.UUID, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.SessionID, nil
}
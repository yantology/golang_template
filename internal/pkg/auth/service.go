package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionNotFound    = errors.New("session not found")
	ErrInvalidSession     = errors.New("invalid session")
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"-"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `json:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type Service struct {
	userRepo      UserRepository
	sessionRepo   SessionRepository
	jwtManager    *JWTManager
	passwordHasher *PasswordHasher
}

type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	UserAgent string `json:"-"`
	IPAddress string `json:"-"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	UserAgent string `json:"-"`
	IPAddress string `json:"-"`
}

type AuthResponse struct {
	User   *User      `json:"user"`
	Tokens *TokenPair `json:"tokens"`
}

func NewService(userRepo UserRepository, sessionRepo SessionRepository, jwtManager *JWTManager) *Service {
	return &Service{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		jwtManager:     jwtManager,
		passwordHasher: NewPasswordHasher(),
	}
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := s.passwordHasher.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create session and tokens
	return s.createSessionAndTokens(ctx, user, req.UserAgent, req.IPAddress)
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	valid, err := s.passwordHasher.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil || !valid {
		return nil, ErrInvalidCredentials
	}

	// Create session and tokens
	return s.createSessionAndTokens(ctx, user, req.UserAgent, req.IPAddress)
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Get session by refresh token
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Check if session is expired
	if session.ExpiresAt.Before(time.Now()) {
		// Delete expired session
		_ = s.sessionRepo.Delete(ctx, session.ID)
		return nil, ErrInvalidSession
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new token pair
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, session.ID)
	if err != nil {
		return nil, err
	}

	// Update session with new refresh token
	session.RefreshToken = tokens.RefreshToken
	session.UpdatedAt = time.Now()
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *Service) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

func (s *Service) LogoutAllSessions(ctx context.Context, userID uuid.UUID) error {
	return s.sessionRepo.DeleteByUserID(ctx, userID)
}

func (s *Service) ValidateToken(ctx context.Context, tokenString string) (*User, *Session, error) {
	// Validate JWT token
	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, nil, err
	}

	// Ensure it's an access token
	if claims.TokenType != AccessToken {
		return nil, nil, ErrInvalidToken
	}

	// Get session to verify it's still active
	session, err := s.sessionRepo.GetByID(ctx, claims.SessionID)
	if err != nil {
		return nil, nil, ErrSessionNotFound
	}

	// Check if session is expired
	if session.ExpiresAt.Before(time.Now()) {
		// Delete expired session
		_ = s.sessionRepo.Delete(ctx, session.ID)
		return nil, nil, ErrInvalidSession
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, nil, ErrUserNotFound
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, nil, ErrInvalidCredentials
	}

	return user, session, nil
}

func (s *Service) CleanupExpiredSessions(ctx context.Context) error {
	return s.sessionRepo.DeleteExpired(ctx)
}

func (s *Service) createSessionAndTokens(ctx context.Context, user *User, userAgent, ipAddress string) (*AuthResponse, error) {
	// Create session
	session := &Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		ExpiresAt: time.Now().Add(24 * time.Hour * 30), // 30 days
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Generate tokens
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, session.ID)
	if err != nil {
		return nil, err
	}

	// Set refresh token in session
	session.RefreshToken = tokens.RefreshToken

	// Save session
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:   user,
		Tokens: tokens,
	}, nil
}
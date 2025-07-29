package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserContextKey    = "user"
	SessionContextKey = "session"
	UserIDContextKey  = "user_id"
)

type Middleware struct {
	authService *Service
}

func NewMiddleware(authService *Service) *Middleware {
	return &Middleware{
		authService: authService,
	}
}

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractTokenFromHeader(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
			})
			c.Abort()
			return
		}

		user, session, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			status := http.StatusUnauthorized
			message := "Invalid token"

			switch err {
			case ErrExpiredToken:
				message = "Token has expired"
			case ErrTokenNotFound:
				message = "Token not found"
			case ErrInvalidSession:
				message = "Session is invalid or expired"
			case ErrUserNotFound:
				message = "User not found"
			}

			c.JSON(status, gin.H{
				"error": message,
			})
			c.Abort()
			return
		}

		// Set user and session in context
		c.Set(UserContextKey, user)
		c.Set(SessionContextKey, session)
		c.Set(UserIDContextKey, user.ID)

		c.Next()
	}
}

func (m *Middleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractTokenFromHeader(c)
		if token == "" {
			c.Next()
			return
		}

		user, session, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err == nil {
			// Set user and session in context only if validation succeeds
			c.Set(UserContextKey, user)
			c.Set(SessionContextKey, session)
			c.Set(UserIDContextKey, user.ID)
		}

		c.Next()
	}
}

func (m *Middleware) extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// Helper functions to extract user information from context

func GetUserFromContext(c *gin.Context) (*User, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}
	
	authUser, ok := user.(*User)
	return authUser, ok
}

func GetSessionFromContext(c *gin.Context) (*Session, bool) {
	session, exists := c.Get(SessionContextKey)
	if !exists {
		return nil, false
	}
	
	authSession, ok := session.(*Session)
	return authSession, ok
}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get(UserIDContextKey)
	if !exists {
		return uuid.Nil, false
	}
	
	id, ok := userID.(uuid.UUID)
	return id, ok
}

// Context helpers for non-Gin contexts

func GetUserFromStdContext(ctx context.Context) (*User, bool) {
	user := ctx.Value(UserContextKey)
	if user == nil {
		return nil, false
	}
	
	authUser, ok := user.(*User)
	return authUser, ok
}

func GetSessionFromStdContext(ctx context.Context) (*Session, bool) {
	session := ctx.Value(SessionContextKey)
	if session == nil {
		return nil, false
	}
	
	authSession, ok := session.(*Session)
	return authSession, ok
}

func GetUserIDFromStdContext(ctx context.Context) (uuid.UUID, bool) {
	userID := ctx.Value(UserIDContextKey)
	if userID == nil {
		return uuid.Nil, false
	}
	
	id, ok := userID.(uuid.UUID)
	return id, ok
}

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, SessionContextKey, session)
}

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}
package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/yantology/golang_template/internal/api/handlers"
	"github.com/yantology/golang_template/internal/api/routes"
	"github.com/yantology/golang_template/internal/config"
	"github.com/yantology/golang_template/pkg/response"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	db     *sql.DB
	router *gin.Engine
	server *http.Server
}

// New creates a new server instance
func New(cfg *config.Config, db *sql.DB) *Server {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add global middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", healthCheckHandler(db))

	server := &Server{
		config: cfg,
		db:     db,
		router: router,
	}

	// Setup API routes
	server.setupRoutes()

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	fmt.Printf("Server starting on %s\n", addr)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Server shutting down...")
	
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	
	return nil
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Initialize handlers
	handler := handlers.NewHandler()

	// API version 1
	v1 := s.router.Group("/api/v1")
	
	// Setup routes
	routes.SetupRoutes(v1, handler)
}


// healthCheckHandler performs a basic health check
func healthCheckHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db != nil {
			if err := db.Ping(); err != nil {
				response.Error(c, http.StatusServiceUnavailable, "Health check failed", err.Error())
				return
			}
		}

		response.Success(c, http.StatusOK, "Health check passed", map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
		})
	}
}
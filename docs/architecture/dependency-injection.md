# Dependency Injection

Dependency Injection (DI) is a crucial pattern in the Go Backend Template that enables testability, maintainability, and loose coupling between components.

## 🎯 Why Dependency Injection?

### Benefits

1. **Testability**: Easy to mock dependencies for unit testing
2. **Flexibility**: Swap implementations without changing dependent code
3. **Maintainability**: Clear separation of concerns
4. **Configuration**: Centralized dependency management

### Without DI (Problematic)

```go
// ❌ Hard to test, tightly coupled
type UserService struct{}

func (s *UserService) CreateUser(user *User) error {
    // Direct dependency creation - hard to test
    db := sql.Open("postgres", "connection-string")
    repo := &PostgreSQLUserRepository{db: db}
    
    return repo.Create(user)
}
```

### With DI (Clean)

```go
// ✅ Easy to test, loosely coupled
type UserService struct {
    userRepo UserRepository // Interface dependency
}

func NewUserService(userRepo UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user *User) error {
    return s.userRepo.Create(user)
}
```

## 🏗️ DI Architecture in the Template

### Dependency Flow

```
main.go
  ↓
config.Load()
  ↓
database.Connect()
  ↓
repositories.New*()
  ↓
services.New*()
  ↓
handlers.New*()
  ↓
routes.Setup()
  ↓
server.Start()
```

## 🔧 Implementation Patterns

### 1. Constructor Injection

The primary pattern used throughout the template:

```go
// Repository layer
type UserRepository interface {
    Create(ctx context.Context, user *models.User) (*models.User, error)
    GetByID(ctx context.Context, id int64) (*models.User, error)
}

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{db: db}
}

// Service layer
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error)
}

type userService struct {
    userRepo     UserRepository
    authService  AuthService
    emailService EmailService
    logger       Logger
}

func NewUserService(
    userRepo UserRepository,
    authService AuthService,
    emailService EmailService,
    logger Logger,
) UserService {
    return &userService{
        userRepo:     userRepo,
        authService:  authService,
        emailService: emailService,
        logger:       logger,
    }
}

// Handler layer
type UserHandler struct {
    userService UserService
    logger      Logger
}

func NewUserHandler(userService UserService, logger Logger) *UserHandler {
    return &UserHandler{
        userService: userService,
        logger:      logger,
    }
}
```

### 2. Interface Segregation

Each dependency is defined as a focused interface:

```go
// Focused interfaces for specific needs
type Logger interface {
    Info(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Debug(msg string, args ...interface{})
}

type EmailService interface {
    SendWelcomeEmail(userEmail, userName string) error
    SendPasswordResetEmail(userEmail, resetToken string) error
}

type AuthService interface {
    HashPassword(password string) (string, error)
    CheckPassword(password, hash string) bool
    GenerateJWT(userID int64, email, role string) (string, error)
    ValidateJWT(token string) (*JWTClaims, error)
}
```

### 3. Dependency Container

A simple container to manage dependency creation:

```go
// internal/container/container.go
package container

import (
    "database/sql"
    
    "github.com/yantology/golang_template/internal/api/handlers"
    "github.com/yantology/golang_template/internal/business/services"
    "github.com/yantology/golang_template/internal/config"
    "github.com/yantology/golang_template/internal/data/repositories"
    "github.com/yantology/golang_template/internal/pkg/auth"
    "github.com/yantology/golang_template/internal/pkg/logger"
)

type Container struct {
    Config   *config.Config
    DB       *sql.DB
    Logger   logger.Logger
    
    // Repositories
    UserRepo     repositories.UserRepository
    ArticleRepo  repositories.ArticleRepository
    CategoryRepo repositories.CategoryRepository
    
    // Services
    AuthService    services.AuthService
    UserService    services.UserService
    ArticleService services.ArticleService
    EmailService   services.EmailService
    
    // Handlers
    UserHandler    *handlers.UserHandler
    ArticleHandler *handlers.ArticleHandler
    AuthHandler    *handlers.AuthHandler
}

func NewContainer(cfg *config.Config) (*Container, error) {
    container := &Container{Config: cfg}
    
    if err := container.initializeInfrastructure(); err != nil {
        return nil, err
    }
    
    container.initializeRepositories()
    container.initializeServices()
    container.initializeHandlers()
    
    return container, nil
}

func (c *Container) initializeInfrastructure() error {
    // Initialize logger
    var err error
    c.Logger, err = logger.New(c.Config.Logger)
    if err != nil {
        return fmt.Errorf("failed to initialize logger: %w", err)
    }
    
    // Initialize database
    c.DB, err = database.Connect(c.Config.Database)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }
    
    return nil
}

func (c *Container) initializeRepositories() {
    c.UserRepo = repositories.NewUserRepository(c.DB)
    c.ArticleRepo = repositories.NewArticleRepository(c.DB)
    c.CategoryRepo = repositories.NewCategoryRepository(c.DB)
}

func (c *Container) initializeServices() {
    // Auth service
    c.AuthService = auth.NewService(c.Config.JWT)
    
    // Email service
    c.EmailService = email.NewService(c.Config.Email, c.Logger)
    
    // Business services
    c.UserService = services.NewUserService(
        c.UserRepo,
        c.AuthService,
        c.EmailService,
        c.Logger,
    )
    
    c.ArticleService = services.NewArticleService(
        c.ArticleRepo,
        c.UserRepo,
        c.CategoryRepo,
        c.Logger,
    )
}

func (c *Container) initializeHandlers() {
    c.UserHandler = handlers.NewUserHandler(c.UserService, c.Logger)
    c.ArticleHandler = handlers.NewArticleHandler(c.ArticleService, c.Logger)
    c.AuthHandler = handlers.NewAuthHandler(c.AuthService, c.UserService, c.Logger)
}

func (c *Container) Close() error {
    if c.DB != nil {
        return c.DB.Close()
    }
    return nil
}
```

## 🚀 Application Bootstrap

### Main Application Wiring

```go
// cmd/api/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/yantology/golang_template/internal/api/routes"
    "github.com/yantology/golang_template/internal/config"
    "github.com/yantology/golang_template/internal/container"
    "github.com/yantology/golang_template/internal/server"
    
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize configuration
    if err := config.InitViper(); err != nil {
        log.Fatalf("Failed to initialize configuration: %v", err)
    }
    
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Initialize dependency container
    container, err := container.NewContainer(cfg)
    if err != nil {
        log.Fatalf("Failed to initialize container: %v", err)
    }
    defer container.Close()
    
    // Setup router and routes
    router := setupRouter(container)
    
    // Create and start server
    srv := server.New(cfg.Server, router, container.Logger)
    
    // Start server in goroutine
    go func() {
        if err := srv.Start(); err != nil {
            container.Logger.Error("Server failed to start", "error", err)
        }
    }()
    
    // Wait for interrupt signal to gracefully shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    container.Logger.Info("Server shutting down...")
    
    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        container.Logger.Error("Server forced to shutdown", "error", err)
    }
    
    container.Logger.Info("Server exited")
}

func setupRouter(c *container.Container) *gin.Engine {
    // Set Gin mode based on environment
    if c.Config.Server.Env == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    router := gin.New()
    
    // Setup routes with injected dependencies
    routes.SetupAPIRoutes(router, &routes.Handlers{
        User:    c.UserHandler,
        Article: c.ArticleHandler,
        Auth:    c.AuthHandler,
        Logger:  c.Logger,
    })
    
    return router
}
```

### Route Setup with DI

```go
// internal/api/routes/routes.go
package routes

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/yantology/golang_template/internal/api/handlers"
    "github.com/yantology/golang_template/internal/api/middleware"
    "github.com/yantology/golang_template/internal/pkg/logger"
)

type Handlers struct {
    User    *handlers.UserHandler
    Article *handlers.ArticleHandler
    Auth    *handlers.AuthHandler
    Logger  logger.Logger
}

func SetupAPIRoutes(router *gin.Engine, h *Handlers) {
    // Global middleware
    router.Use(middleware.CORS())
    router.Use(middleware.Logger(h.Logger))
    router.Use(gin.Recovery())
    
    // Health endpoints
    router.GET("/health", healthCheck)
    router.GET("/ready", readinessCheck)
    
    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // Authentication routes
        auth := v1.Group("/auth")
        {
            auth.POST("/register", h.Auth.Register)
            auth.POST("/login", h.Auth.Login)
            
            // Protected auth routes
            protected := auth.Group("", middleware.AuthRequired())
            {
                protected.POST("/logout", h.Auth.Logout)
                protected.GET("/me", h.Auth.GetCurrentUser)
            }
        }
        
        // User routes
        users := v1.Group("/users", middleware.AuthRequired())
        {
            users.GET("", h.User.ListUsers)
            users.POST("", h.User.CreateUser)
            users.GET("/:id", h.User.GetUser)
            users.PUT("/:id", h.User.UpdateUser)
            users.DELETE("/:id", h.User.DeleteUser)
        }
        
        // Article routes
        articles := v1.Group("/articles")
        {
            // Public routes
            articles.GET("", h.Article.ListArticles)
            articles.GET("/:id", h.Article.GetArticle)
            
            // Protected routes
            protected := articles.Group("", middleware.AuthRequired())
            {
                protected.POST("", h.Article.CreateArticle)
                protected.PUT("/:id", h.Article.UpdateArticle)
                protected.DELETE("/:id", h.Article.DeleteArticle)
            }
        }
    }
}

func healthCheck(c *gin.Context) {
    c.JSON(200, gin.H{"status": "healthy"})
}

func readinessCheck(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ready"})
}
```

## 🧪 Testing with DI

### Mock Generation

```bash
# Install mockery
go install github.com/vektra/mockery/v2@latest

# Generate mocks for interfaces
mockery --all --output=tests/mocks
```

### Unit Test Example

```go
// internal/business/services/user_service_test.go
package services

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "github.com/yantology/golang_template/internal/data/models"
    "github.com/yantology/golang_template/tests/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockUserRepo := &mocks.UserRepository{}
    mockAuthService := &mocks.AuthService{}
    mockEmailService := &mocks.EmailService{}
    mockLogger := &mocks.Logger{}
    
    service := NewUserService(mockUserRepo, mockAuthService, mockEmailService, mockLogger)
    
    req := &CreateUserRequest{
        Email:     "test@example.com",
        FirstName: "John",
        LastName:  "Doe",
        Password:  "password123",
    }
    
    expectedUser := &models.User{
        ID:        1,
        Email:     req.Email,
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }
    
    // Setup mocks
    mockUserRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, errors.New("not found"))
    mockAuthService.On("HashPassword", req.Password).Return("hashed_password", nil)
    mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(expectedUser, nil)
    mockEmailService.On("SendWelcomeEmail", req.Email, req.FirstName).Return(nil)
    mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
    
    // Act
    result, err := service.CreateUser(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedUser.Email, result.Email)
    assert.Equal(t, expectedUser.FirstName, result.FirstName)
    
    // Verify all mocks were called as expected
    mockUserRepo.AssertExpectations(t)
    mockAuthService.AssertExpectations(t)
    mockEmailService.AssertExpectations(t)
}
```

### Integration Test Setup

```go
// tests/integration/setup.go
package integration

import (
    "database/sql"
    "testing"
    
    "github.com/yantology/golang_template/internal/config"
    "github.com/yantology/golang_template/internal/container"
)

func SetupTestContainer(t *testing.T) *container.Container {
    // Load test configuration
    cfg := &config.Config{
        Database: config.DatabaseConfig{
            Host:     "localhost",
            Port:     "5432",
            User:     "test_user",
            Password: "test_password",
            Name:     "test_db",
            SSLMode:  "disable",
        },
        Logger: config.LoggerConfig{
            Level:  "error", // Reduce noise in tests
            Format: "text",
        },
    }
    
    // Create container with test configuration
    container, err := container.NewContainer(cfg)
    if err != nil {
        t.Fatalf("Failed to create test container: %v", err)
    }
    
    // Setup test data if needed
    setupTestData(t, container.DB)
    
    return container
}

func CleanupTestContainer(t *testing.T, container *container.Container) {
    cleanupTestData(t, container.DB)
    container.Close()
}
```

## 🎯 Best Practices

### 1. Interface First
Always define interfaces before implementations:

```go
// Define interface first
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
}

// Then implement
type userService struct {
    // dependencies
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // implementation
}
```

### 2. Constructor Pattern
Use constructor functions for dependency injection:

```go
func NewUserService(
    userRepo UserRepository,
    authService AuthService,
    logger Logger,
) UserService {
    return &userService{
        userRepo:    userRepo,
        authService: authService,
        logger:      logger,
    }
}
```

### 3. Avoid Service Locator
Don't use global variables or service locators:

```go
// ❌ Avoid this
var GlobalUserService UserService

func GetUserService() UserService {
    return GlobalUserService
}

// ✅ Use this instead
func NewUserHandler(userService UserService) *UserHandler {
    return &UserHandler{userService: userService}
}
```

### 4. Minimize Dependencies
Keep the number of dependencies reasonable:

```go
// ❌ Too many dependencies
func NewUserService(
    repo1 Repo1, repo2 Repo2, repo3 Repo3, repo4 Repo4,
    service1 Service1, service2 Service2, service3 Service3,
    util1 Util1, util2 Util2,
) UserService

// ✅ Focused dependencies
func NewUserService(
    userRepo UserRepository,
    authService AuthService,
    logger Logger,
) UserService
```

### 5. Configuration Injection
Inject configuration objects, not individual values:

```go
// ❌ Too many parameters
func NewEmailService(smtpHost string, smtpPort int, username string, password string) EmailService

// ✅ Configuration object
type EmailConfig struct {
    SMTPHost string
    SMTPPort int
    Username string
    Password string
}

func NewEmailService(config EmailConfig) EmailService
```

## 🚀 Next Steps

- **See complete examples**: [Examples](../examples/)
- **Learn API development**: [Creating APIs](../api-development/creating-apis.md)
- **Understand testing**: [Testing APIs](../api-development/testing-apis.md)

---

Dependency injection enables clean, testable, and maintainable code by making dependencies explicit and easily replaceable.
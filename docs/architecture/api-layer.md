# API Layer Architecture

The API layer (`internal/api/`) handles HTTP requests, routing, and cross-cutting concerns. It serves as the entry point for all client interactions with the backend.

## 🎯 Layer Responsibilities

- **HTTP Request Handling**: Process incoming HTTP requests
- **Routing**: Map URLs to appropriate handler functions
- **Middleware**: Cross-cutting concerns (auth, logging, CORS)
- **Request/Response Transformation**: Convert between HTTP and business objects
- **Input Validation**: Validate request data format and structure
- **Error Handling**: Transform business errors to appropriate HTTP responses

## 📂 API Layer Structure

```
internal/api/
├── handlers/           # HTTP request handlers
│   └── handlers.go     # Handler implementations
├── middleware/         # Custom middleware
│   └── (middleware files)
└── routes/             # Route definitions
    └── routes.go       # Route setup
```

## 🔧 Handlers

Handlers process HTTP requests and coordinate with business services.

### Current Handler Implementation

Based on the existing code in `internal/api/handlers/handlers.go`:

```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    // Add service dependencies here as the project grows
    // userService   services.UserService
    // authService   services.AuthService
}

func NewHandler( /* service dependencies */ ) *Handler {
    return &Handler{
        // Initialize service dependencies
    }
}

// HealthCheck provides a basic health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "message": "Service is running",
    })
}
```

### Handler Pattern Template

When adding new handlers, follow this pattern:

```go
// Example: UserHandler for user management
type UserHandler struct {
    userService services.UserService
    logger      logger.Logger
}

func NewUserHandler(userService services.UserService, logger logger.Logger) *UserHandler {
    return &UserHandler{
        userService: userService,
        logger:      logger,
    }
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    user, err := h.userService.CreateUser(c.Request.Context(), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "data":    user,
    })
}

// GetUser handles GET /api/v1/users/{id}
func (h *UserHandler) GetUser(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid user ID",
        })
        return
    }

    user, err := h.userService.GetUserByID(c.Request.Context(), id)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Success",
        "data":    user,
    })
}

// handleError centralizes error handling
func (h *UserHandler) handleError(c *gin.Context, err error) {
    switch err := err.(type) {
    case *errors.ValidationError:
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Validation failed",
            "details": err.Error(),
        })
    case *errors.NotFoundError:
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Resource not found",
        })
    case *errors.UnauthorizedError:
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized",
        })
    default:
        h.logger.Error("Internal server error", "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Internal server error",
        })
    }
}
```

### Request/Response DTOs

Define clear data transfer objects for API communication:

```go
// Request DTOs
type CreateUserRequest struct {
    Email     string `json:"email" binding:"required,email"`
    FirstName string `json:"first_name" binding:"required,min=1,max=100"`
    LastName  string `json:"last_name" binding:"required,min=1,max=100"`
    Password  string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
    FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=1,max=100"`
    LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=1,max=100"`
}

type ListUsersParams struct {
    Page     int    `form:"page" binding:"omitempty,min=1"`
    Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
    Search   string `form:"search" binding:"omitempty,max=100"`
    SortBy   string `form:"sort_by" binding:"omitempty,oneof=name email created_at"`
    SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// Response DTOs
type UserResponse struct {
    ID        int64     `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ListUsersResponse struct {
    Users      []UserResponse `json:"users"`
    Pagination Pagination     `json:"pagination"`
}

type Pagination struct {
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    Total      int64 `json:"total"`
    TotalPages int64 `json:"total_pages"`
}
```

## 🛡️ Middleware

Middleware handles cross-cutting concerns across all routes.

### Authentication Middleware

```go
package middleware

import (
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/yantology/golang_template/internal/pkg/auth"
)

// AuthRequired validates JWT tokens and sets user context
func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header required",
            })
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization format",
            })
            c.Abort()
            return
        }

        claims, err := auth.ValidateJWT(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }

        // Set user context
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("user_role", claims.Role)
        c.Next()
    }
}

// OptionalAuth validates JWT if present but doesn't require it
func OptionalAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.Next()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString != authHeader {
            if claims, err := auth.ValidateJWT(tokenString); err == nil {
                c.Set("user_id", claims.UserID)
                c.Set("user_email", claims.Email)
                c.Set("user_role", claims.Role)
            }
        }
        c.Next()
    }
}
```

### CORS Middleware

```go
package middleware

import (
    "time"
    
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

// CORS middleware using gin-contrib/cors package
func CORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "https://yourdomain.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}

// CORSForDevelopment provides permissive CORS for development
func CORSForDevelopment() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"*"},
        AllowHeaders:     []string{"*"},
        AllowCredentials: false,
        MaxAge:           12 * time.Hour,
    })
}
```

### Logging Middleware

```go
package middleware

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/yantology/golang_template/internal/pkg/logger"
)

// Logger creates a logging middleware
func Logger(logger logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // Process request
        c.Next()

        // Log request details
        latency := time.Since(start)
        clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()

        if raw != "" {
            path = path + "?" + raw
        }

        logger.Info("HTTP Request",
            "status", statusCode,
            "latency", latency,
            "client_ip", clientIP,
            "method", method,
            "path", path,
        )
    }
}
```


## 🛣️ Routes

Routes define the API endpoints and wire them to handlers.

### Current Route Setup

Based on `internal/api/routes/routes.go`:

```go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/yantology/golang_template/internal/api/handlers"
)

func SetupRoutes(r *gin.RouterGroup, h *handlers.Handler) {
    r.GET("/ping", h.HealthCheck)
}
```

### Extended Route Setup

Here's how to extend the routing as the application grows:

```go
package routes

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/yantology/golang_template/internal/api/handlers"
    "github.com/yantology/golang_template/internal/api/middleware"
)

// SetupAPIRoutes configures all API routes
func SetupAPIRoutes(
    router *gin.Engine,
    handlers *handlers.Handlers, // Updated to include all handlers
) {
    // Global middleware
    router.Use(middleware.CORS())
    router.Use(middleware.Logger(handlers.Logger))
    router.Use(gin.Recovery())

    // Health check endpoint (public)
    router.GET("/health", handlers.HealthCheck)
    router.GET("/ready", handlers.ReadinessCheck)

    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // Public routes
        public := v1.Group("/public")
        {
            public.GET("/ping", handlers.Ping)
            public.GET("/version", handlers.Version)
        }

        // Authentication routes (no auth required)
        auth := v1.Group("/auth")
        {
            auth.POST("/register", handlers.Auth.Register)
            auth.POST("/login", handlers.Auth.Login)
            auth.POST("/forgot-password", handlers.Auth.ForgotPassword)
            auth.POST("/reset-password", handlers.Auth.ResetPassword)
            
            // Protected auth routes
            protected := auth.Group("", middleware.AuthRequired())
            {
                protected.POST("/logout", handlers.Auth.Logout)
                protected.POST("/refresh", handlers.Auth.RefreshToken)
                protected.GET("/me", handlers.Auth.GetCurrentUser)
                protected.PUT("/change-password", handlers.Auth.ChangePassword)
            }
        }

        // User management routes
        users := v1.Group("/users")
        users.Use(middleware.AuthRequired()) // All user routes require auth
        {
            users.GET("", handlers.User.ListUsers)
            users.POST("", handlers.User.CreateUser)
            users.GET("/:id", handlers.User.GetUser)
            users.PUT("/:id", handlers.User.UpdateUser)
            users.DELETE("/:id", handlers.User.DeleteUser)
            
            // Admin-only routes
            admin := users.Group("", middleware.AdminRequired())
            {
                admin.POST("/:id/activate", handlers.User.ActivateUser)
                admin.POST("/:id/deactivate", handlers.User.DeactivateUser)
            }
        }

        // Article management routes
        articles := v1.Group("/articles")
        {
            // Public article routes (with optional auth for user context)
            articles.Use(middleware.OptionalAuth())
            articles.GET("", handlers.Article.ListArticles)
            articles.GET("/:id", handlers.Article.GetArticle)
            articles.GET("/slug/:slug", handlers.Article.GetArticleBySlug)
            
            // Protected article routes
            protected := articles.Group("", middleware.AuthRequired())
            {
                protected.POST("", handlers.Article.CreateArticle)
                protected.PUT("/:id", handlers.Article.UpdateArticle)
                protected.DELETE("/:id", handlers.Article.DeleteArticle)
                
                // Admin/Author-only routes
                management := protected.Group("", middleware.AuthorOrAdminRequired())
                {
                    management.POST("/:id/publish", handlers.Article.PublishArticle)
                    management.POST("/:id/archive", handlers.Article.ArchiveArticle)
                }
            }
        }

        // Category management routes
        categories := v1.Group("/categories")
        {
            categories.GET("", handlers.Category.ListCategories)
            categories.GET("/:id", handlers.Category.GetCategory)
            
            // Admin-only category management
            admin := categories.Group("", middleware.AuthRequired(), middleware.AdminRequired())
            {
                admin.POST("", handlers.Category.CreateCategory)
                admin.PUT("/:id", handlers.Category.UpdateCategory)
                admin.DELETE("/:id", handlers.Category.DeleteCategory)
            }
        }
    }
}
```

## 🎯 API Layer Best Practices

### 1. Consistent Response Format

Use a standardized response format across all endpoints:

```go
// pkg/response/response.go
package response

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
    Response
    Pagination Pagination `json:"pagination"`
}

func SuccessJSON(c *gin.Context, status int, message string, data interface{}) {
    c.JSON(status, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func ErrorJSON(c *gin.Context, status int, message string, err error) {
    response := Response{
        Success: false,
        Message: message,
    }
    
    if err != nil {
        response.Error = err.Error()
    }
    
    c.JSON(status, response)
}

func SuccessJSONWithPagination(c *gin.Context, status int, message string, data interface{}, pagination Pagination) {
    c.JSON(status, PaginatedResponse{
        Response: Response{
            Success: true,
            Message: message,
            Data:    data,
        },
        Pagination: pagination,
    })
}
```

### 2. Input Validation

Use Gin's binding features with custom validation:

```go
// Custom validation tags
import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
    validate = validator.New()
    
    // Register custom validation
    validate.RegisterValidation("no_script", validateNoScript)
}

func validateNoScript(fl validator.FieldLevel) bool {
    return !strings.Contains(strings.ToLower(fl.Field().String()), "<script>")
}

// In handler
type CreateArticleRequest struct {
    Title   string `json:"title" binding:"required,min=5,max=200"`
    Content string `json:"content" binding:"required,min=10,no_script"`
    Tags    []string `json:"tags" binding:"dive,required,min=1,max=50"`
}
```

### 3. Error Handling

Implement comprehensive error handling:

```go
func (h *Handler) handleBusinessError(c *gin.Context, err error) {
    switch e := err.(type) {
    case *errors.ValidationError:
        response.ErrorJSON(c, http.StatusBadRequest, "Validation failed", e)
    case *errors.NotFoundError:
        response.ErrorJSON(c, http.StatusNotFound, "Resource not found", e)
    case *errors.UnauthorizedError:
        response.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized", e)
    case *errors.ForbiddenError:
        response.ErrorJSON(c, http.StatusForbidden, "Forbidden", e)
    case *errors.ConflictError:
        response.ErrorJSON(c, http.StatusConflict, "Conflict", e)
    default:
        h.logger.Error("Unhandled error", "error", err)
        response.ErrorJSON(c, http.StatusInternalServerError, "Internal server error", nil)
    }
}
```

### 4. Context Usage

Extract user information from request context:

```go
func (h *UserHandler) UpdateProfile(c *gin.Context) {
    // Extract user context set by auth middleware
    userID, exists := c.Get("user_id")
    if !exists {
        response.ErrorJSON(c, http.StatusUnauthorized, "User not authenticated", nil)
        return
    }
    
    userRole, _ := c.Get("user_role")
    
    // Use context in business logic
    ctx := context.WithValue(c.Request.Context(), "user_id", userID)
    ctx = context.WithValue(ctx, "user_role", userRole)
    
    // Continue with business logic...
}
```

## 🚀 Next Steps

- **Learn about business layer**: [Business Layer](./business-layer.md)
- **Understand data layer**: [Data Layer](./data-layer.md)
- **See practical examples**: [Creating APIs](../api-development/creating-apis.md)

---

The API layer provides a clean, well-structured interface for clients to interact with your business logic while maintaining proper separation of concerns and following REST principles.
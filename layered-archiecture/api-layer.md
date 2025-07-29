# 🌐 API Layer (`internal/api/`)

The API layer handles HTTP requests, routing, and middleware. It serves as the entry point for all client interactions with the backend.

## 📋 Layer Responsibilities

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
├── middleware/         # Custom middleware
└── routes/             # Route definitions
```

## 🔧 Handlers

Handlers process HTTP requests and coordinate with business services.

### 📄 handlers/[entity]_handler.go

```go
package handlers

import (
    "net/http"
    "strconv"
    
    "[module-name]/internal/business/services"
    "[module-name]/pkg/response"
    
    "github.com/gin-gonic/gin"
)

type [Entity]Handler struct {
    [entity]Service services.[Entity]Service
}

func New[Entity]Handler([entity]Service services.[Entity]Service) *[Entity]Handler {
    return &[Entity]Handler{
        [entity]Service: [entity]Service,
    }
}

// Create[Entity] handles POST /[entities]
func (h *[Entity]Handler) Create[Entity](c *gin.Context) {
    var req Create[Entity]Request
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ErrorJSON(c, http.StatusBadRequest, "Invalid request format", err)
        return
    }

    [entity], err := h.[entity]Service.Create[Entity](c.Request.Context(), &req)
    if err != nil {
        response.ErrorJSON(c, http.StatusInternalServerError, "Failed to create [entity]", err)
        return
    }

    response.SuccessJSON(c, http.StatusCreated, "[Entity] created successfully", [entity])
}

// Get[Entity] handles GET /[entities]/{id}
func (h *[Entity]Handler) Get[Entity](c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.ErrorJSON(c, http.StatusBadRequest, "Invalid [entity] ID", err)
        return
    }

    [entity], err := h.[entity]Service.Get[Entity]ByID(c.Request.Context(), id)
    if err != nil {
        response.ErrorJSON(c, http.StatusNotFound, "[Entity] not found", err)
        return
    }

    response.SuccessJSON(c, http.StatusOK, "Success", [entity])
}

// List[Entities] handles GET /[entities]
func (h *[Entity]Handler) List[Entities](c *gin.Context) {
    var params List[Entities]Params
    if err := c.ShouldBindQuery(&params); err != nil {
        response.ErrorJSON(c, http.StatusBadRequest, "Invalid query parameters", err)
        return
    }

    [entities], pagination, err := h.[entity]Service.List[Entities](c.Request.Context(), &params)
    if err != nil {
        response.ErrorJSON(c, http.StatusInternalServerError, "Failed to fetch [entities]", err)
        return
    }

    response.SuccessJSONWithPagination(c, http.StatusOK, "Success", [entities], pagination)
}

// Update[Entity] handles PUT /[entities]/{id}
func (h *[Entity]Handler) Update[Entity](c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.ErrorJSON(c, http.StatusBadRequest, "Invalid [entity] ID", err)
        return
    }

    var req Update[Entity]Request
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ErrorJSON(c, http.StatusBadRequest, "Invalid request format", err)
        return
    }

    [entity], err := h.[entity]Service.Update[Entity](c.Request.Context(), id, &req)
    if err != nil {
        response.ErrorJSON(c, http.StatusInternalServerError, "Failed to update [entity]", err)
        return
    }

    response.SuccessJSON(c, http.StatusOK, "[Entity] updated successfully", [entity])
}

// Delete[Entity] handles DELETE /[entities]/{id}
func (h *[Entity]Handler) Delete[Entity](c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.ErrorJSON(c, http.StatusBadRequest, "Invalid [entity] ID", err)
        return
    }

    err = h.[entity]Service.Delete[Entity](c.Request.Context(), id)
    if err != nil {
        response.ErrorJSON(c, http.StatusInternalServerError, "Failed to delete [entity]", err)
        return
    }

    response.SuccessJSON(c, http.StatusOK, "[Entity] deleted successfully", nil)
}
```

### 📄 Request/Response DTOs

```go
// Create[Entity]Request represents the request payload for creating an entity
type Create[Entity]Request struct {
    Title         string `json:"title" binding:"required,min=1,max=500"`
    Content       string `json:"content" binding:"required"`
    CategoryID    int64  `json:"category_id" binding:"required"`
    SubCategoryID int64  `json:"sub_category_id" binding:"omitempty"`
    FeaturedImage string `json:"featured_image" binding:"omitempty,url"`
}

// Update[Entity]Request represents the request payload for updating an entity
type Update[Entity]Request struct {
    Title         *string `json:"title" binding:"omitempty,min=1,max=500"`
    Content       *string `json:"content" binding:"omitempty"`
    CategoryID    *int64  `json:"category_id" binding:"omitempty"`
    SubCategoryID *int64  `json:"sub_category_id" binding:"omitempty"`
    FeaturedImage *string `json:"featured_image" binding:"omitempty,url"`
    Status        *string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

// List[Entities]Params represents query parameters for listing entities
type List[Entities]Params struct {
    Page         int    `form:"page" binding:"omitempty,min=1"`
    Limit        int    `form:"limit" binding:"omitempty,min=1,max=100"`
    Search       string `form:"search" binding:"omitempty,max=100"`
    CategoryID   int64  `form:"category_id" binding:"omitempty"`
    Status       string `form:"status" binding:"omitempty,oneof=draft published archived"`
    SortBy       string `form:"sort_by" binding:"omitempty,oneof=title created_at updated_at view_count"`
    SortOrder    string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// SetDefaults sets default values for list parameters
func (p *List[Entities]Params) SetDefaults() {
    if p.Page == 0 {
        p.Page = 1
    }
    if p.Limit == 0 {
        p.Limit = 20
    }
    if p.SortBy == "" {
        p.SortBy = "created_at"
    }
    if p.SortOrder == "" {
        p.SortOrder = "desc"
    }
}
```

## 🛡️ Middleware

Middleware handles cross-cutting concerns across all routes.

### 📄 middleware/auth.go

```go
package middleware

import (
    "net/http"
    "strings"
    
    "[module-name]/internal/pkg/auth"
    "[module-name]/pkg/response"
    
    "github.com/gin-gonic/gin"
)

// AuthRequired validates JWT tokens and sets user context
func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.ErrorJSON(c, http.StatusUnauthorized, "Authorization header required", nil)
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            response.ErrorJSON(c, http.StatusUnauthorized, "Invalid authorization format", nil)
            c.Abort()
            return
        }

        claims, err := auth.ValidateJWT(tokenString)
        if err != nil {
            response.ErrorJSON(c, http.StatusUnauthorized, "Invalid or expired token", err)
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

// AdminRequired checks if user has admin role
func AdminRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("user_role")
        if !exists || role != "admin" {
            response.ErrorJSON(c, http.StatusForbidden, "Admin access required", nil)
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 📄 middleware/cors.go

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

// CustomCORS allows dynamic CORS configuration
func CustomCORS(config cors.Config) gin.HandlerFunc {
    return cors.New(config)
}
```


## 🛣️ Routes

Routes define the API endpoints and wire them to handlers.

### 📄 routes/[entity]_routes.go

```go
package routes

import (
    "[module-name]/internal/api/handlers"
    "[module-name]/internal/api/middleware"
    
    "github.com/gin-gonic/gin"
)

func Setup[Entity]Routes(r *gin.RouterGroup, handler *handlers.[Entity]Handler) {
    [entities] := r.Group("/[entities]")
    
    
    {
        // Public routes (with optional auth)
        [entities].Use(middleware.OptionalAuth())
        [entities].GET("", handler.List[Entities])
        [entities].GET("/:id", handler.Get[Entity])
        [entities].GET("/slug/:slug", handler.Get[Entity]BySlug)
        
        // Protected routes (require authentication)
        protected := [entities].Group("", middleware.AuthRequired())
        {
            protected.POST("", handler.Create[Entity])
            protected.PUT("/:id", handler.Update[Entity])
            protected.DELETE("/:id", handler.Delete[Entity])
            
            // Admin only routes
            admin := protected.Group("", middleware.AdminRequired())
            {
                admin.POST("/:id/publish", handler.Publish[Entity])
                admin.POST("/:id/archive", handler.Archive[Entity])
            }
        }
    }
}
```

### 📄 routes/api.go - Main Routes Setup

```go
package routes

import (
    "net/http"
    "time"
    
    "[module-name]/internal/api/handlers"
    "[module-name]/internal/api/middleware"
    
    "github.com/gin-gonic/gin"
)

// SetupAPIRoutes configures all API routes
func SetupAPIRoutes(
    router *gin.Engine,
    authHandler *handlers.AuthHandler,
    userHandler *handlers.UserHandler,
    articleHandler *handlers.ArticleHandler,
    categoryHandler *handlers.CategoryHandler,
) {
    // Global middleware
    router.Use(middleware.CORS())
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
            "timestamp": time.Now().UTC(),
        })
    })

    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // Authentication routes (no auth required)
        SetupAuthRoutes(v1, authHandler)
        
        // User management routes
        SetupUserRoutes(v1, userHandler)
        
        // Content routes
        SetupArticleRoutes(v1, articleHandler)
        SetupCategoryRoutes(v1, categoryHandler)
    }
}
```

### 📄 routes/auth_routes.go

```go
package routes

import (
    "[module-name]/internal/api/handlers"
    "[module-name]/internal/api/middleware"
    
    "github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup, handler *handlers.AuthHandler) {
    auth := r.Group("/auth")
    
    
    {
        // Public authentication endpoints
        auth.POST("/register", handler.Register)
        auth.POST("/login", handler.Login)
        auth.POST("/forgot-password", handler.ForgotPassword)
        auth.POST("/reset-password", handler.ResetPassword)
        auth.POST("/verify-email", handler.VerifyEmail)
        auth.POST("/resend-verification", handler.ResendVerification)
        
        // Protected endpoints (require valid token)
        protected := auth.Group("", middleware.AuthRequired())
        {
            protected.POST("/logout", handler.Logout)
            protected.POST("/refresh", handler.RefreshToken)
            protected.GET("/me", handler.GetCurrentUser)
            protected.PUT("/change-password", handler.ChangePassword)
        }
    }
}
```

## 🎯 Handler Best Practices

### 1. **Error Handling**
```go
func (h *[Entity]Handler) Create[Entity](c *gin.Context) {
    var req Create[Entity]Request
    if err := c.ShouldBindJSON(&req); err != nil {
        // Log the error for debugging
        h.logger.Error("Invalid request format", "error", err, "path", c.Request.URL.Path)
        response.ErrorJSON(c, http.StatusBadRequest, "Invalid request format", err)
        return
    }

    [entity], err := h.[entity]Service.Create[Entity](c.Request.Context(), &req)
    if err != nil {
        // Different error types should map to different HTTP status codes
        switch err := err.(type) {
        case *errors.ValidationError:
            response.ErrorJSON(c, http.StatusBadRequest, "Validation failed", err)
        case *errors.NotFoundError:
            response.ErrorJSON(c, http.StatusNotFound, "Resource not found", err)
        case *errors.UnauthorizedError:
            response.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized", err)
        default:
            h.logger.Error("Failed to create [entity]", "error", err)
            response.ErrorJSON(c, http.StatusInternalServerError, "Internal server error", nil)
        }
        return
    }

    response.SuccessJSON(c, http.StatusCreated, "[Entity] created successfully", [entity])
}
```

### 2. **Input Validation**
```go
// Use struct tags with validator package for validation
type Create[Entity]Request struct {
    Title    string `json:"title" validate:"required,min=1,max=500"`
    Content  string `json:"content" validate:"required,min=10"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"required,min=18,max=120"`
    Website  string `json:"website" validate:"omitempty,url"`
    Tags     []string `json:"tags" validate:"dive,required,min=1,max=50"`
}

// Using go-playground/validator for advanced validation
import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
    validate = validator.New()
    
    // Register custom validation tags
    validate.RegisterValidation("no_script", validateNoScript)
}

// Custom validation function
func validateNoScript(fl validator.FieldLevel) bool {
    return !strings.Contains(strings.ToLower(fl.Field().String()), "<script>")
}

// Validate request with detailed error messages
func (r *Create[Entity]Request) Validate() error {
    if err := validate.Struct(r); err != nil {
        var validationErrors []string
        for _, err := range err.(validator.ValidationErrors) {
            validationErrors = append(validationErrors, 
                fmt.Sprintf("Field '%s' failed validation '%s'", err.Field(), err.Tag()))
        }
        return fmt.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
    }
    return nil
}

// Handler with proper validation
func (h *[Entity]Handler) Create[Entity](c *gin.Context) {
    var req Create[Entity]Request
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ErrorJSON(c, http.StatusBadRequest, "Invalid JSON format", err)
        return
    }
    
    // Validate request with custom validator
    if err := req.Validate(); err != nil {
        response.ErrorJSON(c, http.StatusBadRequest, "Validation failed", err)
        return
    }
    
    // Continue with business logic...
}
```

### 3. **Response Consistency**
```go
// Always use consistent response format
response.SuccessJSON(c, http.StatusOK, "Success message", data)
response.ErrorJSON(c, http.StatusBadRequest, "Error message", error)

// For paginated responses
response.SuccessJSONWithPagination(c, http.StatusOK, "Success", data, pagination)
```

### 4. **Context Usage**
```go
func (h *[Entity]Handler) Get[Entity](c *gin.Context) {
    // Extract user context set by middleware
    userID, exists := c.Get("user_id")
    if exists {
        // Use user context for authorization or logging
        h.logger.Info("User accessing [entity]", "user_id", userID, "[entity]_id", id)
    }

    // Pass gin context to service layer for cancellation
    [entity], err := h.[entity]Service.Get[Entity]ByID(c.Request.Context(), id)
    // ...
}
```

---

**Previous**: [← Project Structure](./01-project-structure.md) | **Next**: [Business Layer →](./03-business-layer.md)

**Last Updated:** [YYYY-MM-DD]
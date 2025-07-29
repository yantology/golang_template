# 🚀 Development Workflow - New API/Module Creation

This document provides a comprehensive, step-by-step workflow for creating new APIs or modules using the layered architecture pattern. Follow this guide to ensure consistency, maintainability, and adherence to established patterns.

**🚨 PREREQUISITE:** Before starting this workflow, complete the [Planning Design Checklist](./planing_design.md) to ensure all business and technical planning is ready.

## 📋 Overview

The development workflow follows a **bottom-up approach**:
1. **api-specification, database and third-party ready** → Validate planning checklist completion
2. **Data Layer** → Models, repositories, database operations
3. **Business Layer** → Services, validators, business logic
4. **API Layer** → Handlers, routes, middleware
5. **Testing** → Unit, integration, manual, and E2E tests
6. **Documentation & Deployment** → Final steps and deployment

---

## 🎯 Phase 1: api-specification, database and third-party ready

### 1.1 Pre-Implementation Validation

**📝 Checklist:**
- [ ] **Complete Planning Checklist**: Use [Planning Design Checklist](./planing_design.md) to validate:
  - [ ] PRD Template completed
  - [ ] User Flow templates mapped
  - [ ] Requirements Analysis finished
  - [ ] API Specification contracts defined
  - [ ] Database infrastructure ready
- [ ] **All Prerequisites Met**: Ensure all planning documents are approved and ready

**📄 Example: Article Management Module**
```
Purpose: Manage blog articles with categories and user authorship
Entities: Article, Category, User (existing)
Operations: Create, Read, Update, Delete, Publish, Archive, List with filters
Business Rules:
- Only draft articles can be published
- Users can only edit their own articles
- Published articles must have featured image and min content length
- Articles must belong to valid categories
```

### 1.2 Database Design

#### 1.2.1 Database Infrastructure Check

**📝 Checklist:**
- [ ] **Infrastructure Validated**: Database setup verified in planning checklist
- [ ] **Schema Design Ready**: Entity relationships and table schemas defined
- [ ] **Migration Scripts Prepared**: UP/DOWN migrations created and tested

**📖 Reference:** [Planning Design Checklist](./planing_design.md) - Database Infrastructure section

#### 1.2.2 Database Schema Design

**📝 Checklist:**
- [ ] Design entity-relationship diagram
- [ ] Define table schemas with proper indexes
- [ ] Plan foreign key relationships
- [ ] Consider soft deletes vs hard deletes
- [ ] Plan for data migration if needed

**📄 Database Schema Example:**
```sql
-- articles table
CREATE TABLE articles (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(500) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt VARCHAR(500),
    featured_image VARCHAR(500),
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    view_count BIGINT DEFAULT 0,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    sub_category_id BIGINT REFERENCES categories(id),
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- indexes
CREATE INDEX idx_articles_user_id ON articles(user_id);
CREATE INDEX idx_articles_category_id ON articles(category_id);
CREATE INDEX idx_articles_status ON articles(status);
CREATE INDEX idx_articles_published_at ON articles(published_at);
```

### 1.3 API Contract Definition

**📝 Checklist:**
- [ ] **API Specification Complete**: Contracts defined using `/templates/tecnical/api_specifiation/` templates
- [ ] **Endpoint Templates Selected**: Appropriate CRUD and custom operation templates chosen
- [ ] **Request/Response Schemas**: Defined with validation rules and examples
- [ ] **Error Handling Strategy**: HTTP status codes and error response formats documented

**📖 Reference:** [API Specification Templates](../../api_specifiation/README.md)

**📄 API Contract Example:**
```yaml
# Article API Endpoints
POST   /api/v1/articles              # Create article
GET    /api/v1/articles              # List articles (with pagination/filters)
GET    /api/v1/articles/{id}         # Get article by ID
GET    /api/v1/articles/slug/{slug}  # Get article by slug
PUT    /api/v1/articles/{id}         # Update article
DELETE /api/v1/articles/{id}         # Delete article
POST   /api/v1/articles/{id}/publish # Publish article (admin/author)
POST   /api/v1/articles/{id}/archive # Archive article (admin/author)

# Query Parameters for List:
?page=1&limit=20&search=keyword&category_id=1&status=published&sort_by=created_at&sort_order=desc
```

---

## 🗄️ Phase 2: Data Layer Implementation

### 2.1 Create Models

**📂 Location:** `internal/data/models/[entity].go`

**📝 Checklist:**
- [ ] Define struct with proper tags (json, db, validate)
- [ ] Include related entities for JOIN operations
- [ ] Add validation methods
- [ ] Add helper methods (SetDefaults, GenerateSlug, etc.)

**📄 Implementation:**
```go
// internal/data/models/article.go
package models

import (
    "fmt"
    "strings"
    "time"
)

type Article struct {
    ID            int64      `json:"id" db:"id"`
    Title         string     `json:"title" db:"title"`
    Slug          string     `json:"slug" db:"slug"`
    Content       string     `json:"content" db:"content"`
    Excerpt       string     `json:"excerpt" db:"excerpt"`
    FeaturedImage string     `json:"featured_image,omitempty" db:"featured_image"`
    Status        string     `json:"status" db:"status"`
    ViewCount     int64      `json:"view_count" db:"view_count"`
    UserID        int64      `json:"user_id" db:"user_id"`
    CategoryID    int64      `json:"category_id" db:"category_id"`
    SubCategoryID int64      `json:"sub_category_id" db:"sub_category_id"`
    PublishedAt   *time.Time `json:"published_at,omitempty" db:"published_at"`
    CreatedAt     time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
    
    // Related entities (loaded via JOIN queries)
    User        User        `json:"user,omitempty"`
    Category    Category    `json:"category,omitempty"`
    SubCategory SubCategory `json:"sub_category,omitempty"`
}

// GenerateSlug creates a URL-friendly slug from title
func (a *Article) GenerateSlug() {
    a.Slug = generateSlug(a.Title)
}

// SetDefaults sets default values for new articles
func (a *Article) SetDefaults() {
    if a.Status == "" {
        a.Status = "draft"
    }
    if a.ViewCount == 0 {
        a.ViewCount = 0
    }
    if a.Excerpt == "" && a.Content != "" {
        a.Excerpt = generateExcerpt(a.Content, 160)
    }
}

// Validate performs basic validation
func (a *Article) Validate() error {
    if a.Title == "" {
        return fmt.Errorf("title is required")
    }
    if a.UserID == 0 {
        return fmt.Errorf("user_id is required")
    }
    if a.CategoryID == 0 {
        return fmt.Errorf("category_id is required")
    }
    
    validStatuses := map[string]bool{
        "draft": true, "published": true, "archived": true,
    }
    if !validStatuses[a.Status] {
        return fmt.Errorf("invalid status: %s", a.Status)
    }
    
    return nil
}

// IsPublished returns true if the article is published
func (a *Article) IsPublished() bool {
    return a.Status == "published"
}

// Helper functions
func generateSlug(title string) string {
    slug := strings.ToLower(title)
    slug = strings.ReplaceAll(slug, " ", "-")
    // Additional slug processing...
    return slug
}

func generateExcerpt(content string, maxLength int) string {
    if len(content) <= maxLength {
        return content
    }
    return content[:maxLength] + "..."
}
```

### 2.2 Define Repository Interface

**📂 Location:** `internal/data/repositories/interfaces.go`

**📝 Checklist:**
- [ ] Define interface with all required methods
- [ ] Include parameter types for complex queries
- [ ] Consider context for cancellation
- [ ] Plan for both simple and complex operations

**📄 Implementation:**
```go
// Add to internal/data/repositories/interfaces.go
type ArticleRepository interface {
    Create(ctx context.Context, article *models.Article) (*models.Article, error)
    GetByID(ctx context.Context, id int64) (*models.Article, error)
    GetBySlug(ctx context.Context, slug string) (*models.Article, error)
    List(ctx context.Context, params *ListArticlesParams) ([]*models.Article, int64, error)
    Update(ctx context.Context, article *models.Article) (*models.Article, error)
    Delete(ctx context.Context, id int64) error
    IncrementViewCount(ctx context.Context, id int64) error
    GetByUserID(ctx context.Context, userID int64, params *ListArticlesParams) ([]*models.Article, int64, error)
}

type ListArticlesParams struct {
    Page         int
    Limit        int
    Search       string
    CategoryID   int64
    Status       string
    SortBy       string
    SortOrder    string
    UserID       int64
}
```

### 2.3 Implement Repository

**📂 Location:** `internal/data/repositories/article_repository.go`

**📝 Checklist:**
- [ ] Implement all interface methods
- [ ] Use query builder (squirrel) for complex queries
- [ ] Handle errors appropriately
- [ ] Optimize queries with proper JOINs
- [ ] Add proper filtering and pagination

**📄 Key Implementation Points:**
```go
// Key methods to implement:
// 1. Create - INSERT with RETURNING
// 2. GetByID - SELECT with JOINs for related data
// 3. List - Complex SELECT with filtering, pagination, sorting
// 4. Update - UPDATE with optimistic locking
// 5. Delete - DELETE or soft delete

// Example query structure:
func (r *articleRepository) List(ctx context.Context, params *ListArticlesParams) ([]*models.Article, int64, error) {
    baseQuery := r.qb.Select("COUNT(*)").From("articles a")
    selectQuery := r.qb.Select(
        "a.id", "a.title", "a.slug", "a.content", "a.excerpt", "a.featured_image",
        "a.status", "a.view_count", "a.user_id", "a.category_id", "a.sub_category_id",
        "a.published_at", "a.created_at", "a.updated_at",
        "u.email", "u.first_name", "u.last_name", "u.role",
        "c.name as category_name", "c.slug as category_slug",
        "sc.name as sub_category_name", "sc.slug as sub_category_slug",
    ).From("articles a").
        LeftJoin("users u ON u.id = a.user_id").
        LeftJoin("categories c ON c.id = a.category_id").
        LeftJoin("categories sc ON sc.id = a.sub_category_id")
    
    // Apply filters, sorting, pagination...
}
```

### 2.4 Database Migration (if needed)

**📂 Location:** `internal/data/migrations/`

**📝 Checklist:**
- [ ] Create migration file with timestamp
- [ ] Include both UP and DOWN migrations
- [ ] Test migrations on sample data
- [ ] Consider data migration for existing records

---

## 💼 Phase 3: Business Layer Implementation

### 3.1 Create Service Interface

**📂 Location:** `internal/business/services/article_service.go`

**📝 Checklist:**
- [ ] Define service interface with business operations
- [ ] Include request/response types
- [ ] Plan for complex business workflows
- [ ] Consider error handling and validation

**📄 Implementation:**
```go
type ArticleService interface {
    CreateArticle(ctx context.Context, req *CreateArticleRequest) (*models.Article, error)
    GetArticleByID(ctx context.Context, id int64) (*models.Article, error)
    GetArticleBySlug(ctx context.Context, slug string) (*models.Article, error)
    ListArticles(ctx context.Context, params *ListArticlesParams) ([]*models.Article, *Pagination, error)
    UpdateArticle(ctx context.Context, id int64, req *UpdateArticleRequest) (*models.Article, error)
    DeleteArticle(ctx context.Context, id int64) error
    PublishArticle(ctx context.Context, id int64, userID int64) error
    ArchiveArticle(ctx context.Context, id int64, userID int64) error
    IncrementViewCount(ctx context.Context, id int64) error
}

// Request types
type CreateArticleRequest struct {
    Title         string `json:"title"`
    Content       string `json:"content"`
    CategoryID    int64  `json:"category_id"`
    SubCategoryID int64  `json:"sub_category_id,omitempty"`
    FeaturedImage string `json:"featured_image,omitempty"`
    UserID        int64  `json:"-"` // Set from context
}

type UpdateArticleRequest struct {
    Title         *string `json:"title,omitempty"`
    Content       *string `json:"content,omitempty"`
    CategoryID    *int64  `json:"category_id,omitempty"`
    SubCategoryID *int64  `json:"sub_category_id,omitempty"`
    FeaturedImage *string `json:"featured_image,omitempty"`
    Status        *string `json:"status,omitempty"`
    UserID        int64   `json:"-"`
}
```

### 3.2 Implement Service

**📝 Checklist:**
- [ ] Implement constructor with dependencies
- [ ] Add business rule validation
- [ ] Implement authorization checks
- [ ] Handle transactions for complex operations
- [ ] Add proper error handling and logging

**📄 Key Implementation Points:**
```go
type articleService struct {
    articleRepo  repositories.ArticleRepository
    userRepo     repositories.UserRepository
    categoryRepo repositories.CategoryRepository
    validator    *ArticleValidator
    logger       Logger
}

// Key business logic to implement:
func (s *articleService) CreateArticle(ctx context.Context, req *CreateArticleRequest) (*models.Article, error) {
    // 1. Business validation
    if err := s.validator.ValidateCreateArticle(ctx, req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. Verify dependencies (category exists)
    if _, err := s.categoryRepo.GetByID(ctx, req.CategoryID); err != nil {
        return nil, errors.NewValidationError("invalid category_id")
    }

    // 3. Create model and set defaults
    article := &models.Article{
        Title:      req.Title,
        Content:    req.Content,
        CategoryID: req.CategoryID,
        UserID:     req.UserID,
        Status:     "draft",
    }
    article.GenerateSlug()
    article.SetDefaults()

    // 4. Save and return
    return s.articleRepo.Create(ctx, article)
}

func (s *articleService) PublishArticle(ctx context.Context, id int64, userID int64) error {
    // 1. Get article
    article, err := s.articleRepo.GetByID(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to get article: %w", err)
    }

    // 2. Check permissions
    if err := s.checkUpdatePermission(ctx, article, userID); err != nil {
        return err
    }

    // 3. Business validation for publishing
    if article.Status != "draft" {
        return errors.NewValidationError("can only publish draft articles")
    }
    
    if err := s.validator.ValidatePublishArticle(article); err != nil {
        return fmt.Errorf("publish validation failed: %w", err)
    }

    // 4. Update status and published date
    article.Status = "published"
    now := time.Now()
    article.PublishedAt = &now

    // 5. Save changes
    _, err = s.articleRepo.Update(ctx, article)
    return err
}
```

### 3.3 Create Validators

**📂 Location:** `internal/business/validators/article_validator.go`

**📝 Checklist:**
- [ ] Implement business rule validation
- [ ] Validate complex scenarios
- [ ] Add custom validation rules
- [ ] Return detailed error messages

**📄 Implementation:**
```go
type ArticleValidator struct{}

func NewArticleValidator() *ArticleValidator {
    return &ArticleValidator{}
}

func (v *ArticleValidator) ValidateCreateArticle(ctx context.Context, req *CreateArticleRequest) error {
    if err := v.validateTitle(req.Title); err != nil {
        return err
    }
    if err := v.validateContent(req.Content); err != nil {
        return err
    }
    return nil
}

func (v *ArticleValidator) ValidatePublishArticle(article *models.Article) error {
    // Business rule: article must have minimum content length for publishing
    if utf8.RuneCountInString(article.Content) < 100 {
        return errors.NewValidationError("article content must be at least 100 characters for publishing")
    }
    
    // Business rule: article must have featured image for publishing
    if article.FeaturedImage == "" {
        return errors.NewValidationError("featured image is required for publishing")
    }
    
    return nil
}

func (v *ArticleValidator) validateTitle(title string) error {
    if title == "" {
        return errors.NewValidationError("title is required")
    }
    if utf8.RuneCountInString(title) < 5 {
        return errors.NewValidationError("title must be at least 5 characters")
    }
    if utf8.RuneCountInString(title) > 500 {
        return errors.NewValidationError("title must not exceed 500 characters")
    }
    return nil
}
```

---

## 🌐 Phase 4: API Layer Implementation

### 4.1 Create Request/Response DTOs

**📂 Location:** `internal/api/handlers/article_types.go`

**📝 Checklist:**
- [ ] Define request structs with validation tags
- [ ] Define response structures
- [ ] Add JSON tags for proper serialization
- [ ] Include validation rules using struct tags

**📄 Implementation:**
```go
type CreateArticleRequest struct {
    Title         string `json:"title" binding:"required,min=1,max=500"`
    Content       string `json:"content" binding:"required"`
    CategoryID    int64  `json:"category_id" binding:"required"`
    SubCategoryID int64  `json:"sub_category_id" binding:"omitempty"`
    FeaturedImage string `json:"featured_image" binding:"omitempty,url"`
}

type UpdateArticleRequest struct {
    Title         *string `json:"title" binding:"omitempty,min=1,max=500"`
    Content       *string `json:"content" binding:"omitempty"`
    CategoryID    *int64  `json:"category_id" binding:"omitempty"`
    SubCategoryID *int64  `json:"sub_category_id" binding:"omitempty"`
    FeaturedImage *string `json:"featured_image" binding:"omitempty,url"`
    Status        *string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

type ListArticlesParams struct {
    Page         int    `form:"page" binding:"omitempty,min=1"`
    Limit        int    `form:"limit" binding:"omitempty,min=1,max=100"`
    Search       string `form:"search" binding:"omitempty,max=100"`
    CategoryID   int64  `form:"category_id" binding:"omitempty"`
    Status       string `form:"status" binding:"omitempty,oneof=draft published archived"`
    SortBy       string `form:"sort_by" binding:"omitempty,oneof=title created_at updated_at view_count"`
    SortOrder    string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}
```

### 4.2 Implement Handlers

**📂 Location:** `internal/api/handlers/article_handler.go`

**📝 Checklist:**
- [ ] Implement all CRUD handlers
- [ ] Add proper request validation
- [ ] Extract user context from middleware
- [ ] Return consistent response formats
- [ ] Handle errors appropriately

**📄 Implementation:**
```go
type ArticleHandler struct {
    articleService services.ArticleService
    logger         Logger
}

func NewArticleHandler(articleService services.ArticleService, logger Logger) *ArticleHandler {
    return &ArticleHandler{
        articleService: articleService,
        logger:         logger,
    }
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    var req CreateArticleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ErrorJSON(c, http.StatusBadRequest, "Invalid request format", err)
        return
    }

    // Extract user ID from context (set by auth middleware)
    userID, exists := c.Get("user_id")
    if !exists {
        response.ErrorJSON(c, http.StatusUnauthorized, "User not authenticated", nil)
        return
    }

    // Convert to service request
    serviceReq := &services.CreateArticleRequest{
        Title:         req.Title,
        Content:       req.Content,
        CategoryID:    req.CategoryID,
        SubCategoryID: req.SubCategoryID,
        FeaturedImage: req.FeaturedImage,
        UserID:        userID.(int64),
    }

    article, err := h.articleService.CreateArticle(c.Request.Context(), serviceReq)
    if err != nil {
        h.handleServiceError(c, err, "Failed to create article")
        return
    }

    response.SuccessJSON(c, http.StatusCreated, "Article created successfully", article)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.ErrorJSON(c, http.StatusBadRequest, "Invalid article ID", err)
        return
    }

    article, err := h.articleService.GetArticleByID(c.Request.Context(), id)
    if err != nil {
        h.handleServiceError(c, err, "Failed to get article")
        return
    }

    response.SuccessJSON(c, http.StatusOK, "Success", article)
}

func (h *ArticleHandler) handleServiceError(c *gin.Context, err error, message string) {
    switch err := err.(type) {
    case *errors.ValidationError:
        response.ErrorJSON(c, http.StatusBadRequest, "Validation failed", err)
    case *errors.NotFoundError:
        response.ErrorJSON(c, http.StatusNotFound, "Resource not found", err)
    case *errors.UnauthorizedError:
        response.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized", err)
    default:
        h.logger.Error(message, "error", err)
        response.ErrorJSON(c, http.StatusInternalServerError, "Internal server error", nil)
    }
}
```

### 4.3 Define Routes

**📂 Location:** `internal/api/routes/article_routes.go`

**📝 Checklist:**
- [ ] Define route groups
- [ ] Apply appropriate middleware
- [ ] Set up authentication/authorization

**📄 Implementation:**
```go
func SetupArticleRoutes(r *gin.RouterGroup, handler *handlers.ArticleHandler) {
    articles := r.Group("/articles")
    
    
    {
        // Public routes (with optional auth)
        articles.Use(middleware.OptionalAuth())
        articles.GET("", handler.ListArticles)
        articles.GET("/:id", handler.GetArticle)
        articles.GET("/slug/:slug", handler.GetArticleBySlug)
        
        // Protected routes (require authentication)
        protected := articles.Group("", middleware.AuthRequired())
        {
            protected.POST("", handler.CreateArticle)
            protected.PUT("/:id", handler.UpdateArticle)
            protected.DELETE("/:id", handler.DeleteArticle)
            
            // Special actions
            protected.POST("/:id/publish", handler.PublishArticle)
            protected.POST("/:id/archive", handler.ArchiveArticle)
        }
    }
}
```

---

## 🧪 Phase 5: Testing Implementation

### 5.1 Unit Tests

**📝 Checklist:**
- [ ] Test all service methods with mocks
- [ ] Test repository methods with SQL mocks
- [ ] Test handlers with HTTP test framework
- [ ] Test validation logic
- [ ] Achieve > 80% code coverage

**📄 Test Structure:**
```
tests/
├── unit/
│   ├── services/
│   │   └── article_service_test.go
│   ├── repositories/
│   │   └── article_repository_test.go
│   ├── handlers/
│   │   └── article_handler_test.go
│   └── validators/
│       └── article_validator_test.go
```

### 5.2 Integration Tests

**📝 Checklist:**
- [ ] Test API endpoints with real database
- [ ] Test database operations with real DB
- [ ] Test middleware integration
- [ ] Test authentication and authorization flows

### 5.3 End-to-End Tests

**📝 Checklist:**
- [ ] Test complete user workflows
- [ ] Test error scenarios
- [ ] Test performance under load
- [ ] Test security scenarios

**📄 Example E2E Test:**
```go
func (suite *ArticleWorkflowTestSuite) TestCompleteArticleLifecycle() {
    // 1. User registration and login
    // 2. Create draft article
    // 3. Update article with featured image
    // 4. Publish article
    // 5. Verify article is publicly visible
    // 6. Archive article
    // 7. Verify article is no longer public
}
```

### 5.4 Manual Testing

Manual testing is essential for validating real-world API behavior and catching issues that automated tests might miss.

**📝 Checklist:**
- [ ] Set up API testing tools (Postman, cURL)
- [ ] Test all CRUD operations manually
- [ ] Validate authentication and authorization flows
- [ ] Test business logic and edge cases
- [ ] Verify error handling and responses
- [ ] Test performance and security scenarios

#### 5.4.1 API Testing Tools Setup

**🔧 Postman Collection Setup:**

Create a Postman collection for your API with the following structure:

```
Article API Collection/
├── Authentication/
│   ├── Register User
│   ├── Login User
│   └── Get Current User
├── Articles/
│   ├── Create Article
│   ├── Get Article by ID
│   ├── Get Article by Slug
│   ├── List Articles
│   ├── Update Article
│   ├── Delete Article
│   ├── Publish Article
│   └── Archive Article
└── Categories/
    ├── List Categories
    └── Get Category by ID
```

**🌐 Environment Variables:**
```json
{
  "base_url": "http://localhost:8080/api/v1",
  "auth_token": "{{auth_token}}",
  "user_id": "{{user_id}}",
  "article_id": "{{article_id}}"
}
```

**📝 Pre-request Script for Authentication:**
```javascript
// Auto-set auth token from login response
if (pm.response.json() && pm.response.json().data && pm.response.json().data.access_token) {
    pm.environment.set("auth_token", pm.response.json().data.access_token);
}
```

#### 5.4.2 cURL Commands Reference

**🔐 Authentication Testing:**

```bash
# Register new user
curl -X POST "{{base_url}}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Login user
curl -X POST "{{base_url}}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Get current user (with token)
curl -X GET "{{base_url}}/auth/me" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**📝 Article CRUD Testing:**

```bash
# Create article
curl -X POST "{{base_url}}/articles" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "My Test Article",
    "content": "This is the content of my test article. It should be long enough for validation.",
    "category_id": 1,
    "featured_image": "https://example.com/image.jpg"
  }'

# Get article by ID
curl -X GET "{{base_url}}/articles/1"

# Get article by slug
curl -X GET "{{base_url}}/articles/slug/my-test-article"

# List articles with filters
curl -X GET "{{base_url}}/articles?page=1&limit=10&status=published&category_id=1&search=test"

# Update article
curl -X PUT "{{base_url}}/articles/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "Updated Article Title",
    "content": "Updated content with more details."
  }'

# Publish article
curl -X POST "{{base_url}}/articles/1/publish" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Archive article
curl -X POST "{{base_url}}/articles/1/archive" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Delete article
curl -X DELETE "{{base_url}}/articles/1" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

#### 5.4.3 Test Scenarios by Category

**✅ Happy Path Testing:**

1. **Complete Article Lifecycle:**
   ```
   ✓ Register user
   ✓ Login and get token
   ✓ Create draft article
   ✓ Update article with featured image
   ✓ Publish article
   ✓ View published article publicly
   ✓ Archive article
   ✓ Verify article no longer public
   ```

2. **CRUD Operations:**
   ```
   ✓ Create article with all fields
   ✓ Create article with minimal fields
   ✓ Read article by ID
   ✓ Read article by slug
   ✓ Update partial fields
   ✓ Update all fields
   ✓ Delete article
   ```

**❌ Error Scenario Testing:**

1. **Authentication Errors:**
   ```bash
   # Test without token
   curl -X POST "{{base_url}}/articles" \
     -H "Content-Type: application/json" \
     -d '{"title": "Test"}'
   # Expected: 401 Unauthorized

   # Test with invalid token
   curl -X POST "{{base_url}}/articles" \
     -H "Authorization: Bearer invalid_token" \
     -H "Content-Type: application/json" \
     -d '{"title": "Test"}'
   # Expected: 401 Unauthorized

   # Test with expired token
   curl -X POST "{{base_url}}/articles" \
     -H "Authorization: Bearer expired_token" \
     -H "Content-Type: application/json" \
     -d '{"title": "Test"}'
   # Expected: 401 Unauthorized
   ```

2. **Validation Errors:**
   ```bash
   # Test missing required fields
   curl -X POST "{{base_url}}/articles" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{}'
   # Expected: 400 Bad Request with validation errors

   # Test invalid field types
   curl -X POST "{{base_url}}/articles" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "title": "",
       "content": "x",
       "category_id": "invalid"
     }'
   # Expected: 400 Bad Request

   # Test field length limits
   curl -X POST "{{base_url}}/articles" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "title": "' + "x".repeat(501) + '",
       "content": "Valid content",
       "category_id": 1
     }'
   # Expected: 400 Bad Request
   ```

3. **Authorization Errors:**
   ```bash
   # Test updating another user's article
   curl -X PUT "{{base_url}}/articles/OTHER_USER_ARTICLE_ID" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"title": "Hacked"}'
   # Expected: 403 Forbidden

   # Test deleting another user's article
   curl -X DELETE "{{base_url}}/articles/OTHER_USER_ARTICLE_ID" \
     -H "Authorization: Bearer YOUR_TOKEN"
   # Expected: 403 Forbidden
   ```

4. **Business Logic Errors:**
   ```bash
   # Test publishing without featured image
   curl -X POST "{{base_url}}/articles" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "title": "Article Without Image",
       "content": "Short content",
       "category_id": 1
     }' | \
   curl -X POST "{{base_url}}/articles/ARTICLE_ID/publish" \
     -H "Authorization: Bearer YOUR_TOKEN"
   # Expected: 400 Bad Request - featured image required

   # Test publishing already published article
   curl -X POST "{{base_url}}/articles/PUBLISHED_ARTICLE_ID/publish" \
     -H "Authorization: Bearer YOUR_TOKEN"
   # Expected: 400 Bad Request - already published
   ```

#### 5.4.4 Input Validation Testing

**📝 Field Validation Tests:**

```bash
# Title validation
curl -X POST "{{base_url}}/articles" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "",
    "content": "Valid content",
    "category_id": 1
  }'
# Expected: 400 - title required

# Content validation
curl -X POST "{{base_url}}/articles" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Valid Title",
    "content": "",
    "category_id": 1
  }'
# Expected: 400 - content required

# Category validation
curl -X POST "{{base_url}}/articles" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Valid Title",
    "content": "Valid content",
    "category_id": 99999
  }'
# Expected: 400 - invalid category

# URL validation for featured image
curl -X POST "{{base_url}}/articles" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Valid Title",
    "content": "Valid content",
    "category_id": 1,
    "featured_image": "not-a-url"
  }'
# Expected: 400 - invalid URL format
```

#### 5.4.5 Security Testing

**🔒 Security Test Cases:**

```bash
# SQL Injection Attempt
curl -X GET "{{base_url}}/articles?search='; DROP TABLE articles; --"
# Expected: Safe handling, no SQL injection

# XSS Prevention
curl -X POST "{{base_url}}/articles" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "<script>alert(\"XSS\")</script>",
    "content": "Valid content",
    "category_id": 1
  }'
# Expected: Script tags should be sanitized or rejected

# CORS Testing
curl -X OPTIONS "{{base_url}}/articles" \
  -H "Origin: https://malicious-site.com"
# Expected: Proper CORS headers, restricted origins

```

#### 5.4.6 Performance Testing

**⚡ Performance Test Scenarios:**

```bash
# Response Time Testing
time curl -X GET "{{base_url}}/articles?limit=100"
# Expected: < 500ms for list operations

# Large Dataset Pagination
curl -X GET "{{base_url}}/articles?page=100&limit=50"
# Expected: Consistent response time regardless of page

# Concurrent Request Testing
for i in {1..10}; do
  (curl -X GET "{{base_url}}/articles" &)
done
wait
# Expected: All requests succeed within reasonable time

# Database Query Performance
curl -X GET "{{base_url}}/articles?search=keyword&category_id=1&status=published"
# Monitor database query logs for optimization opportunities
```

#### 5.4.7 Manual Testing Checklist

**📋 Pre-Testing Setup:**

- [ ] Test environment is running and accessible
- [ ] Database is seeded with test data
- [ ] Authentication service is working
- [ ] Postman/testing tools are configured
- [ ] Environment variables are set correctly

**📋 Core Functionality Testing:**

**Authentication & Authorization:**
- [ ] User registration works with valid data
- [ ] User registration fails with invalid data
- [ ] User login works with correct credentials
- [ ] User login fails with incorrect credentials
- [ ] Protected endpoints require authentication
- [ ] Users can only access their own resources
- [ ] Admin users have elevated permissions

**Article CRUD Operations:**
- [ ] Create article with all required fields
- [ ] Create article fails with missing fields
- [ ] Create article fails with invalid data
- [ ] Get article by ID returns correct data
- [ ] Get article by slug returns correct data
- [ ] Get non-existent article returns 404
- [ ] List articles with pagination works
- [ ] List articles with filters works
- [ ] List articles with search works
- [ ] Update article with valid data
- [ ] Update article fails with invalid data
- [ ] Delete article removes it from database
- [ ] Delete non-existent article returns 404

**Business Logic:**
- [ ] Draft articles are not publicly visible
- [ ] Published articles are publicly visible
- [ ] Archived articles are not publicly visible
- [ ] Publishing requires featured image
- [ ] Publishing requires minimum content length
- [ ] Slug generation works correctly
- [ ] Slug uniqueness is enforced
- [ ] View count increments on article access
- [ ] Category relationships work correctly

**Error Handling:**
- [ ] All error responses have consistent format
- [ ] Error messages are descriptive and helpful
- [ ] Internal errors don't expose sensitive information
- [ ] HTTP status codes are appropriate
- [ ] Validation errors include field-specific messages

**📋 Edge Cases & Security:**

- [ ] Very long input values are handled
- [ ] Special characters in input are handled
- [ ] Unicode content works correctly
- [ ] Malformed JSON requests are rejected
- [ ] SQL injection attempts are prevented
- [ ] XSS attempts are prevented
- [ ] CSRF protection works (if implemented)
- [ ] CORS policies are enforced

**📋 Performance & Load:**

- [ ] Response times are acceptable (< 500ms)
- [ ] Large dataset queries perform well
- [ ] Pagination scales properly
- [ ] Concurrent requests are handled
- [ ] Database connections are managed properly
- [ ] Memory usage is reasonable under load

**📋 Integration:**

- [ ] Email notifications work (if implemented)
- [ ] File uploads work (if implemented)
- [ ] External API calls work (if implemented)
- [ ] Database transactions work correctly
- [ ] Caching works as expected (if implemented)

#### 5.4.8 Test Data Management

**🗃️ Test Data Setup:**

```sql
-- Sample test data for manual testing
INSERT INTO users (email, password_hash, first_name, last_name, role) VALUES
('admin@example.com', '$2a$10$hash', 'Admin', 'User', 'admin'),
('user1@example.com', '$2a$10$hash', 'John', 'Doe', 'user'),
('user2@example.com', '$2a$10$hash', 'Jane', 'Smith', 'user');

INSERT INTO categories (name, slug, description) VALUES
('Technology', 'technology', 'Tech related articles'),
('Sports', 'sports', 'Sports and fitness'),
('Travel', 'travel', 'Travel and adventure');

INSERT INTO articles (title, slug, content, status, user_id, category_id, featured_image) VALUES
('Published Article', 'published-article', 'This is a published article content...', 'published', 2, 1, 'https://example.com/image1.jpg'),
('Draft Article', 'draft-article', 'This is a draft article content...', 'draft', 2, 1, NULL),
('Another User Article', 'another-user-article', 'This belongs to another user...', 'published', 3, 2, 'https://example.com/image2.jpg');
```

**🧹 Test Data Cleanup:**

```bash
# Clean up test data after testing
curl -X DELETE "{{base_url}}/articles/TEST_ARTICLE_ID" -H "Authorization: Bearer TOKEN"
# Or use database cleanup scripts
```

---

## 📚 Phase 6: Documentation & Deployment

### 6.1 API Documentation

**📝 Checklist:**
- [ ] Document all endpoints with examples
- [ ] Include request/response schemas
- [ ] Document error responses
- [ ] Add authentication requirements
- [ ] Provide usage examples

**📄 Tools:**
- Swagger/OpenAPI specification
- Postman collections
- API documentation generators

### 6.2 Code Documentation

**📝 Checklist:**
- [ ] Add package comments
- [ ] Document public functions and types
- [ ] Include usage examples in comments
- [ ] Document business rules and constraints

### 6.3 Deployment Checklist

**📝 Checklist:**
- [ ] Database migrations tested
- [ ] Environment configuration updated
- [ ] Security review completed
- [ ] Performance testing passed
- [ ] Monitoring and logging configured
- [ ] Error handling and recovery tested
- [ ] Load balancer configuration updated
- [ ] SSL/TLS certificates validated

### 6.4 Monitoring & Observability

**📝 Checklist:**
- [ ] Add metrics for key operations
- [ ] Set up health checks
- [ ] Configure alerting for errors
- [ ] Add distributed tracing
- [ ] Log important business events

---

## 🎯 Development Best Practices

### 1. **Code Organization**
- Follow the established layered architecture
- Keep each layer focused on its responsibilities
- Use dependency injection for loose coupling

### 2. **Error Handling**
- Use custom error types for different scenarios
- Log errors with sufficient context
- Don't expose internal errors to API consumers

### 3. **Security**
- Validate all inputs at multiple layers
- Implement proper authentication and authorization
- Never log sensitive information
- Use parameterized queries to prevent SQL injection

### 4. **Performance**
- Use database indexes strategically
- Implement pagination for list operations
- Use connection pooling
- Consider caching for read-heavy operations

### 5. **Testing**
- Write tests as you develop (TDD approach)
- Mock external dependencies in unit tests
- Use real databases for integration tests
- Maintain high test coverage

---

## 🚀 Quick Start Template

For rapid development, copy and modify these template files:

**📁 File Template Checklist:**
- [ ] Copy `article_*` files and rename for your entity
- [ ] Update struct fields and database schema
- [ ] Modify business logic and validation rules
- [ ] Update API endpoints and routes
- [ ] Adapt tests for your use case
- [ ] Update documentation

---

## 🔧 Common Pitfalls & Solutions

### 1. **Database Design Issues**
❌ **Problem:** Missing indexes on frequently queried columns
✅ **Solution:** Add indexes for foreign keys, status fields, and search columns

### 2. **Business Logic in Wrong Layer**
❌ **Problem:** Database validation in repository layer
✅ **Solution:** Keep data validation in models, business rules in services

### 3. **Inconsistent Error Handling**
❌ **Problem:** Different error formats across endpoints
✅ **Solution:** Use centralized error handling with custom error types

### 4. **Security Vulnerabilities**
❌ **Problem:** Missing authorization checks
✅ **Solution:** Implement authorization at service layer, not just API layer

### 5. **Poor Test Coverage**
❌ **Problem:** Testing only happy path scenarios
✅ **Solution:** Test error cases, edge cases, and security scenarios

---

## 📝 Development Checklist

**Phase 1: Planning ✅**
- [ ] Requirements defined
- [ ] Database infrastructure checked and configured
- [ ] Database schema designed
- [ ] API contract specified

**Phase 2: Data Layer ✅**
- [ ] Models created and validated
- [ ] Repository interface defined
- [ ] Repository implemented and tested
- [ ] Database migration created

**Phase 3: Business Layer ✅**
- [ ] Service interface defined
- [ ] Service implemented with business logic
- [ ] Validators created
- [ ] Authorization implemented

**Phase 4: API Layer ✅**
- [ ] Request/response DTOs created
- [ ] Handlers implemented
- [ ] Routes configured
- [ ] Middleware applied

**Phase 5: Testing ✅**
- [ ] Unit tests written and passing
- [ ] Integration tests implemented
- [ ] E2E tests covering workflows
- [ ] Manual testing completed
- [ ] Performance tests completed

**Phase 6: Documentation & Deployment ✅**
- [ ] API documentation updated
- [ ] Code documented
- [ ] Deployment checklist completed
- [ ] Monitoring configured

---

**Navigation:**
- **Previous**: [← Testing Strategy](./testing-strategy.md)
- **Next**: [Configuration →](../configuration/workflow.md)

**Last Updated:** 2025-01-28
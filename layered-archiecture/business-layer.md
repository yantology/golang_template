# 💼 Business Layer (`internal/business/`)

The business layer contains the core business logic, rules, and workflows of the application. It acts as the brain of the system, coordinating between the API layer and data layer.

## 📋 Layer Responsibilities

- **Business Logic**: Core domain logic and business rules
- **Workflow Orchestration**: Coordinate complex operations across multiple repositories
- **Data Validation**: Business rule validation beyond basic input validation
- **Transaction Management**: Handle complex transactions spanning multiple operations
- **External Service Integration**: Coordinate with external APIs and services
- **Business Events**: Emit and handle domain events

## 📂 Business Layer Structure

```
internal/business/
├── services/           # Business services
├── validators/         # Business rule validators
├── events/            # Domain events (optional)
└── workflows/         # Complex workflow orchestration (optional)
```

## 🎯 Services

Services implement business logic and orchestrate operations between multiple data sources.

### 📄 services/[entity]_service.go

```go
package services

import (
    "context"
    "fmt"
    "time"
    
    "[module-name]/internal/data/models"
    "[module-name]/internal/data/repositories"
    "[module-name]/pkg/errors"
)

// [Entity]Service interface defines business operations
type [Entity]Service interface {
    Create[Entity](ctx context.Context, req *Create[Entity]Request) (*models.[Entity], error)
    Get[Entity]ByID(ctx context.Context, id int64) (*models.[Entity], error)
    Get[Entity]BySlug(ctx context.Context, slug string) (*models.[Entity], error)
    List[Entities](ctx context.Context, params *List[Entities]Params) ([]*models.[Entity], *Pagination, error)
    Update[Entity](ctx context.Context, id int64, req *Update[Entity]Request) (*models.[Entity], error)
    Delete[Entity](ctx context.Context, id int64) error
    Publish[Entity](ctx context.Context, id int64, userID int64) error
    Archive[Entity](ctx context.Context, id int64, userID int64) error
    IncrementViewCount(ctx context.Context, id int64) error
}

// [entity]Service implements [Entity]Service interface
type [entity]Service struct {
    [entity]Repo     repositories.[Entity]Repository
    userRepo         repositories.UserRepository
    categoryRepo     repositories.CategoryRepository
    emailService     EmailService
    validator        *[Entity]Validator
    logger           Logger
}

// New[Entity]Service creates a new [entity] service instance
func New[Entity]Service(
    [entity]Repo repositories.[Entity]Repository,
    userRepo repositories.UserRepository,
    categoryRepo repositories.CategoryRepository,
    emailService EmailService,
    logger Logger,
) [Entity]Service {
    return &[entity]Service{
        [entity]Repo: [entity]Repo,
        userRepo:     userRepo,
        categoryRepo: categoryRepo,
        emailService: emailService,
        validator:    New[Entity]Validator(),
        logger:       logger,
    }
}

// Create[Entity] creates a new [entity] with business validation
func (s *[entity]Service) Create[Entity](ctx context.Context, req *Create[Entity]Request) (*models.[Entity], error) {
    // Business validation
    if err := s.validator.ValidateCreate[Entity](ctx, req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Verify category exists
    if _, err := s.categoryRepo.GetByID(ctx, req.CategoryID); err != nil {
        return nil, errors.NewValidationError("invalid category_id")
    }

    // Verify subcategory if provided
    if req.SubCategoryID != 0 {
        subCategory, err := s.categoryRepo.GetByID(ctx, req.SubCategoryID)
        if err != nil {
            return nil, errors.NewValidationError("invalid sub_category_id")
        }
        // Ensure subcategory belongs to the main category
        if subCategory.ParentID != req.CategoryID {
            return nil, errors.NewValidationError("subcategory does not belong to specified category")
        }
    }

    // Create model from request
    [entity] := &models.[Entity]{
        Title:         req.Title,
        Content:       req.Content,
        CategoryID:    req.CategoryID,
        SubCategoryID: req.SubCategoryID,
        Status:        "draft",
        UserID:        req.UserID,
    }

    // Generate slug from title
    [entity].GenerateSlug()

    // Set defaults
    [entity].SetDefaults()

    // Validate model
    if err := [entity].Validate(); err != nil {
        return nil, fmt.Errorf("model validation failed: %w", err)
    }

    // Save to repository
    created[Entity], err := s.[entity]Repo.Create(ctx, [entity])
    if err != nil {
        s.logger.Error("Failed to create [entity]", "error", err, "user_id", req.UserID)
        return nil, fmt.Errorf("failed to create [entity]: %w", err)
    }

    // Log success
    s.logger.Info("[Entity] created successfully", 
        "[entity]_id", created[Entity].ID, 
        "user_id", created[Entity].UserID,
        "title", created[Entity].Title)

    return created[Entity], nil
}

// Get[Entity]ByID retrieves an [entity] by ID with business logic
func (s *[entity]Service) Get[Entity]ByID(ctx context.Context, id int64) (*models.[Entity], error) {
    [entity], err := s.[entity]Repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get [entity]: %w", err)
    }

    // Apply business rules for visibility
    if err := s.checkViewPermission(ctx, [entity]); err != nil {
        return nil, err
    }

    return [entity], nil
}

// List[Entities] retrieves [entities] with pagination and filtering
func (s *[entity]Service) List[Entities](ctx context.Context, params *List[Entities]Params) ([]*models.[Entity], *Pagination, error) {
    // Set defaults
    params.SetDefaults()

    // Validate parameters
    if err := s.validator.ValidateList[Entities]Params(params); err != nil {
        return nil, nil, fmt.Errorf("invalid parameters: %w", err)
    }

    // Get total count and [entities]
    [entities], total, err := s.[entity]Repo.List(ctx, params)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to list [entities]: %w", err)
    }

    // Apply business rules for each [entity]
    filtered[Entities] := make([]*models.[Entity], 0, len([entities]))
    for _, [entity] := range [entities] {
        if err := s.checkViewPermission(ctx, [entity]); err == nil {
            filtered[Entities] = append(filtered[Entities], [entity])
        }
    }

    // Calculate pagination
    pagination := &Pagination{
        Page:       params.Page,
        Limit:      params.Limit,
        Total:      total,
        TotalPages: (total + int64(params.Limit) - 1) / int64(params.Limit),
    }

    return filtered[Entities], pagination, nil
}

// Update[Entity] updates an existing [entity] with business validation
func (s *[entity]Service) Update[Entity](ctx context.Context, id int64, req *Update[Entity]Request) (*models.[Entity], error) {
    // Get existing [entity]
    existing[Entity], err := s.[entity]Repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get [entity]: %w", err)
    }

    // Check update permission
    if err := s.checkUpdatePermission(ctx, existing[Entity], req.UserID); err != nil {
        return nil, err
    }

    // Validate update request
    if err := s.validator.ValidateUpdate[Entity](ctx, req, existing[Entity]); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Apply updates
    updated[Entity] := s.applyUpdates(existing[Entity], req)

    // Regenerate slug if title changed
    if req.Title != nil && *req.Title != existing[Entity].Title {
        updated[Entity].GenerateSlug()
    }

    // Validate updated model
    if err := updated[Entity].Validate(); err != nil {
        return nil, fmt.Errorf("updated model validation failed: %w", err)
    }

    // Save updates
    result, err := s.[entity]Repo.Update(ctx, updated[Entity])
    if err != nil {
        s.logger.Error("Failed to update [entity]", "error", err, "[entity]_id", id)
        return nil, fmt.Errorf("failed to update [entity]: %w", err)
    }

    s.logger.Info("[Entity] updated successfully", "[entity]_id", id, "user_id", req.UserID)
    return result, nil
}

// Publish[Entity] publishes a draft [entity]
func (s *[entity]Service) Publish[Entity](ctx context.Context, id int64, userID int64) error {
    [entity], err := s.[entity]Repo.GetByID(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to get [entity]: %w", err)
    }

    // Check permission
    if err := s.checkUpdatePermission(ctx, [entity], userID); err != nil {
        return err
    }

    // Business rule: can only publish draft [entities]
    if [entity].Status != "draft" {
        return errors.NewValidationError("can only publish draft [entities]")
    }

    // Validate [entity] is ready for publishing
    if err := s.validator.ValidatePublish[Entity]([entity]); err != nil {
        return fmt.Errorf("publish validation failed: %w", err)
    }

    // Update status and published date
    [entity].Status = "published"
    now := time.Now()
    [entity].PublishedAt = &now

    if _, err := s.[entity]Repo.Update(ctx, [entity]); err != nil {
        return fmt.Errorf("failed to publish [entity]: %w", err)
    }

    // Send notification email to subscribers (async)
    go func() {
        if err := s.emailService.NotifyNewPublication([entity]); err != nil {
            s.logger.Error("Failed to send publication notification", "error", err, "[entity]_id", id)
        }
    }()

    s.logger.Info("[Entity] published successfully", "[entity]_id", id, "user_id", userID)
    return nil
}

// IncrementViewCount increments the view count for an [entity]
func (s *[entity]Service) IncrementViewCount(ctx context.Context, id int64) error {
    return s.[entity]Repo.IncrementViewCount(ctx, id)
}

// Private helper methods

// checkViewPermission checks if the current user can view the [entity]
func (s *[entity]Service) checkViewPermission(ctx context.Context, [entity] *models.[Entity]) error {
    // Published [entities] are visible to everyone
    if [entity].Status == "published" {
        return nil
    }

    // Get current user from context
    userID, ok := ctx.Value("user_id").(int64)
    if !ok {
        // Anonymous users can only see published [entities]
        return errors.NewUnauthorizedError("not authorized to view this [entity]")
    }

    // Users can see their own [entities] regardless of status
    if [entity].UserID == userID {
        return nil
    }

    // Admins can see all [entities]
    userRole, ok := ctx.Value("user_role").(string)
    if ok && userRole == "admin" {
        return nil
    }

    return errors.NewUnauthorizedError("not authorized to view this [entity]")
}

// checkUpdatePermission checks if the user can update the [entity]
func (s *[entity]Service) checkUpdatePermission(ctx context.Context, [entity] *models.[Entity], userID int64) error {
    // Users can update their own [entities]
    if [entity].UserID == userID {
        return nil
    }

    // Admins can update any [entity]
    userRole, ok := ctx.Value("user_role").(string)
    if ok && userRole == "admin" {
        return nil
    }

    return errors.NewUnauthorizedError("not authorized to update this [entity]")
}

// applyUpdates applies the update request to the existing [entity]
func (s *[entity]Service) applyUpdates(existing *models.[Entity], req *Update[Entity]Request) *models.[Entity] {
    updated := *existing // Copy the struct

    if req.Title != nil {
        updated.Title = *req.Title
    }
    if req.Content != nil {
        updated.Content = *req.Content
    }
    if req.CategoryID != nil {
        updated.CategoryID = *req.CategoryID
    }
    if req.SubCategoryID != nil {
        updated.SubCategoryID = *req.SubCategoryID
    }
    if req.FeaturedImage != nil {
        updated.FeaturedImage = *req.FeaturedImage
    }
    if req.Status != nil {
        updated.Status = *req.Status
    }

    return &updated
}
```

## 🔍 Validators

Validators implement business rule validation that goes beyond basic input validation.

### 📄 validators/[entity]_validator.go

```go
package validators

import (
    "context"
    "fmt"
    "strings"
    "unicode/utf8"
    
    "[module-name]/internal/data/models"
    "[module-name]/pkg/errors"
)

type [Entity]Validator struct {
    // Add dependencies like repositories if needed for validation
}

func New[Entity]Validator() *[Entity]Validator {
    return &[Entity]Validator{}
}

// ValidateCreate[Entity] validates business rules for creating an [entity]
func (v *[Entity]Validator) ValidateCreate[Entity](ctx context.Context, req *Create[Entity]Request) error {
    if err := v.validateTitle(req.Title); err != nil {
        return err
    }

    if err := v.validateContent(req.Content); err != nil {
        return err
    }

    return nil
}

// ValidateUpdate[Entity] validates business rules for updating an [entity]
func (v *[Entity]Validator) ValidateUpdate[Entity](ctx context.Context, req *Update[Entity]Request, existing *models.[Entity]) error {
    if req.Title != nil {
        if err := v.validateTitle(*req.Title); err != nil {
            return err
        }
    }

    if req.Content != nil {
        if err := v.validateContent(*req.Content); err != nil {
            return err
        }
    }

    // Business rule: published [entities] can only have limited updates
    if existing.Status == "published" {
        if req.Status != nil && *req.Status == "draft" {
            return errors.NewValidationError("cannot change published [entity] back to draft")
        }
    }

    return nil
}

// ValidatePublish[Entity] validates business rules for publishing an [entity]
func (v *[Entity]Validator) ValidatePublish[Entity]([entity] *models.[Entity]) error {
    if err := v.validateTitle([entity].Title); err != nil {
        return fmt.Errorf("title validation failed: %w", err)
    }

    if err := v.validateContent([entity].Content); err != nil {
        return fmt.Errorf("content validation failed: %w", err)
    }

    // Business rule: [entity] must have minimum content length for publishing
    if utf8.RuneCountInString([entity].Content) < 100 {
        return errors.NewValidationError("[entity] content must be at least 100 characters for publishing")
    }

    // Business rule: [entity] must have featured image for publishing
    if [entity].FeaturedImage == "" {
        return errors.NewValidationError("featured image is required for publishing")
    }

    return nil
}

// ValidateList[Entities]Params validates parameters for listing [entities]
func (v *[Entity]Validator) ValidateList[Entities]Params(params *List[Entities]Params) error {
    if params.Page < 1 {
        return errors.NewValidationError("page must be greater than 0")
    }

    if params.Limit < 1 || params.Limit > 100 {
        return errors.NewValidationError("limit must be between 1 and 100")
    }

    if params.Search != "" {
        if len(params.Search) < 2 {
            return errors.NewValidationError("search query must be at least 2 characters")
        }
        if len(params.Search) > 100 {
            return errors.NewValidationError("search query must not exceed 100 characters")
        }
    }

    validSortFields := map[string]bool{
        "title":      true,
        "created_at": true,
        "updated_at": true,
        "view_count": true,
    }
    if !validSortFields[params.SortBy] {
        return errors.NewValidationError("invalid sort field")
    }

    validSortOrders := map[string]bool{
        "asc":  true,
        "desc": true,
    }
    if !validSortOrders[params.SortOrder] {
        return errors.NewValidationError("invalid sort order")
    }

    return nil
}

// Private validation methods

func (v *[Entity]Validator) validateTitle(title string) error {
    if title == "" {
        return errors.NewValidationError("title is required")
    }

    if utf8.RuneCountInString(title) < 5 {
        return errors.NewValidationError("title must be at least 5 characters")
    }

    if utf8.RuneCountInString(title) > 500 {
        return errors.NewValidationError("title must not exceed 500 characters")
    }

    // Business rule: title cannot contain certain words
    bannedWords := []string{"spam", "click here", "free money"}
    titleLower := strings.ToLower(title)
    for _, word := range bannedWords {
        if strings.Contains(titleLower, word) {
            return errors.NewValidationError(fmt.Sprintf("title cannot contain '%s'", word))
        }
    }

    // Business rule: title cannot be all caps
    if title == strings.ToUpper(title) && utf8.RuneCountInString(title) > 10 {
        return errors.NewValidationError("title cannot be all uppercase")
    }

    return nil
}

func (v *[Entity]Validator) validateContent(content string) error {
    if content == "" {
        return errors.NewValidationError("content is required")
    }

    if utf8.RuneCountInString(content) < 10 {
        return errors.NewValidationError("content must be at least 10 characters")
    }

    if utf8.RuneCountInString(content) > 50000 {
        return errors.NewValidationError("content must not exceed 50,000 characters")
    }

    // Business rule: content cannot contain malicious scripts
    if strings.Contains(strings.ToLower(content), "<script>") {
        return errors.NewValidationError("content cannot contain script tags")
    }

    return nil
}
```

## 📊 Supporting Types

### 📄 services/types.go

```go
package services

import "time"

// Request types for service operations
type Create[Entity]Request struct {
    Title         string `json:"title"`
    Content       string `json:"content"`
    CategoryID    int64  `json:"category_id"`
    SubCategoryID int64  `json:"sub_category_id,omitempty"`
    FeaturedImage string `json:"featured_image,omitempty"`
    UserID        int64  `json:"-"` // Set from context, not from request
}

type Update[Entity]Request struct {
    Title         *string `json:"title,omitempty"`
    Content       *string `json:"content,omitempty"`
    CategoryID    *int64  `json:"category_id,omitempty"`
    SubCategoryID *int64  `json:"sub_category_id,omitempty"`
    FeaturedImage *string `json:"featured_image,omitempty"`
    Status        *string `json:"status,omitempty"`
    UserID        int64   `json:"-"` // Set from context
}

type List[Entities]Params struct {
    Page         int    `form:"page"`
    Limit        int    `form:"limit"`
    Search       string `form:"search"`
    CategoryID   int64  `form:"category_id"`
    Status       string `form:"status"`
    SortBy       string `form:"sort_by"`
    SortOrder    string `form:"sort_order"`
}

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

// Pagination represents pagination information
type Pagination struct {
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    Total      int64 `json:"total"`
    TotalPages int64 `json:"total_pages"`
}

// Common interfaces
type Logger interface {
    Info(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Debug(msg string, args ...interface{})
}

type EmailService interface {
    SendWelcomeEmail(userEmail, userName string) error
    SendPasswordResetEmail(userEmail, resetToken string) error
    NotifyNewPublication([entity] *models.[Entity]) error
}

type HTTPClient interface {
    Get(url string) (*HTTPResponse, error)
    Post(url string, body interface{}) (*HTTPResponse, error)
    Put(url string, body interface{}) (*HTTPResponse, error)
    Delete(url string) (*HTTPResponse, error)
}

type HTTPResponse struct {
    StatusCode int
    Body       []byte
    Headers    map[string]string
}
```

## 🌐 External HTTP Client Services

### 📄 HTTP Client Implementation with Resty

```go
// internal/pkg/httpclient/resty_client.go
package httpclient

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/go-resty/resty/v2"
    "[module-name]/internal/config"
    "[module-name]/internal/pkg/logger"
)

type RestyClient struct {
    client *resty.Client
    logger logger.Logger
}

type ClientConfig struct {
    BaseURL        string
    Timeout        time.Duration
    RetryCount     int
    RetryWaitTime  time.Duration
    Debug          bool
    DefaultHeaders map[string]string
}

func NewRestyClient(cfg ClientConfig, logger logger.Logger) *RestyClient {
    client := resty.New()
    
    // Basic configuration
    client.SetBaseURL(cfg.BaseURL)
    client.SetTimeout(cfg.Timeout)
    client.SetRetryCount(cfg.RetryCount)
    client.SetRetryWaitTime(cfg.RetryWaitTime)
    
    // Set default headers
    if cfg.DefaultHeaders != nil {
        client.SetHeaders(cfg.DefaultHeaders)
    }
    
    // Enable debug mode
    if cfg.Debug {
        client.SetDebug(true)
    }
    
    // Add logging middleware
    client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
        logger.Info("HTTP Request",
            "method", req.Method,
            "url", req.URL,
            "headers", req.Header)
        return nil
    })
    
    client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
        logger.Info("HTTP Response",
            "status", resp.StatusCode(),
            "time", resp.Time(),
            "size", len(resp.Body()))
        return nil
    })
    
    return &RestyClient{
        client: client,
        logger: logger,
    }
}

// Get performs HTTP GET request
func (rc *RestyClient) Get(ctx context.Context, endpoint string, result interface{}) error {
    resp, err := rc.client.R().
        SetContext(ctx).
        SetResult(result).
        Get(endpoint)
    
    if err != nil {
        return fmt.Errorf("GET request failed: %w", err)
    }
    
    if resp.IsError() {
        return fmt.Errorf("GET request failed with status %d: %s", 
            resp.StatusCode(), resp.String())
    }
    
    return nil
}

// Post performs HTTP POST request
func (rc *RestyClient) Post(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
    resp, err := rc.client.R().
        SetContext(ctx).
        SetBody(body).
        SetResult(result).
        Post(endpoint)
    
    if err != nil {
        return fmt.Errorf("POST request failed: %w", err)
    }
    
    if resp.IsError() {
        return fmt.Errorf("POST request failed with status %d: %s", 
            resp.StatusCode(), resp.String())
    }
    
    return nil
}

// Put performs HTTP PUT request
func (rc *RestyClient) Put(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
    resp, err := rc.client.R().
        SetContext(ctx).
        SetBody(body).
        SetResult(result).
        Put(endpoint)
    
    if err != nil {
        return fmt.Errorf("PUT request failed: %w", err)
    }
    
    if resp.IsError() {
        return fmt.Errorf("PUT request failed with status %d: %s", 
            resp.StatusCode(), resp.String())
    }
    
    return nil
}

// Delete performs HTTP DELETE request
func (rc *RestyClient) Delete(ctx context.Context, endpoint string) error {
    resp, err := rc.client.R().
        SetContext(ctx).
        Delete(endpoint)
    
    if err != nil {
        return fmt.Errorf("DELETE request failed: %w", err)
    }
    
    if resp.IsError() {
        return fmt.Errorf("DELETE request failed with status %d: %s", 
            resp.StatusCode(), resp.String())
    }
    
    return nil
}

// PostWithAuth performs authenticated POST request
func (rc *RestyClient) PostWithAuth(ctx context.Context, endpoint string, body interface{}, token string, result interface{}) error {
    resp, err := rc.client.R().
        SetContext(ctx).
        SetAuthToken(token).
        SetBody(body).
        SetResult(result).
        Post(endpoint)
    
    if err != nil {
        return fmt.Errorf("authenticated POST request failed: %w", err)
    }
    
    if resp.IsError() {
        return fmt.Errorf("authenticated POST request failed with status %d: %s", 
            resp.StatusCode(), resp.String())
    }
    
    return nil
}
```

### 📄 External Service Example

```go
// internal/business/services/external_api_service.go
package services

import (
    "context"
    "fmt"
    "time"
    
    "[module-name]/internal/pkg/httpclient"
    "[module-name]/internal/pkg/logger"
)

type ExternalAPIService interface {
    GetUserProfile(ctx context.Context, userID string) (*UserProfile, error)
    SendNotification(ctx context.Context, notification *Notification) error
    ValidateEmail(ctx context.Context, email string) (*EmailValidation, error)
}

type externalAPIService struct {
    client *httpclient.RestyClient
    logger logger.Logger
    config ExternalAPIConfig
}

type ExternalAPIConfig struct {
    BaseURL    string
    APIKey     string
    Timeout    time.Duration
    RetryCount int
}

type UserProfile struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Verified bool   `json:"verified"`
}

type Notification struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
    Type    string `json:"type"`
}

type EmailValidation struct {
    Email   string `json:"email"`
    Valid   bool   `json:"valid"`
    Reason  string `json:"reason,omitempty"`
}

func NewExternalAPIService(config ExternalAPIConfig, logger logger.Logger) ExternalAPIService {
    clientConfig := httpclient.ClientConfig{
        BaseURL:       config.BaseURL,
        Timeout:       config.Timeout,
        RetryCount:    config.RetryCount,
        RetryWaitTime: 1 * time.Second,
        DefaultHeaders: map[string]string{
            "Content-Type":  "application/json",
            "Authorization": "Bearer " + config.APIKey,
            "User-Agent":    "YourApp/1.0",
        },
    }
    
    client := httpclient.NewRestyClient(clientConfig, logger)
    
    return &externalAPIService{
        client: client,
        logger: logger,
        config: config,
    }
}

func (s *externalAPIService) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
    var profile UserProfile
    
    endpoint := fmt.Sprintf("/users/%s/profile", userID)
    
    if err := s.client.Get(ctx, endpoint, &profile); err != nil {
        s.logger.Error("Failed to get user profile", "user_id", userID, "error", err)
        return nil, fmt.Errorf("failed to get user profile: %w", err)
    }
    
    return &profile, nil
}

func (s *externalAPIService) SendNotification(ctx context.Context, notification *Notification) error {
    var response map[string]interface{}
    
    if err := s.client.Post(ctx, "/notifications", notification, &response); err != nil {
        s.logger.Error("Failed to send notification", "notification", notification, "error", err)
        return fmt.Errorf("failed to send notification: %w", err)
    }
    
    s.logger.Info("Notification sent successfully", "to", notification.To, "type", notification.Type)
    return nil
}

func (s *externalAPIService) ValidateEmail(ctx context.Context, email string) (*EmailValidation, error) {
    var validation EmailValidation
    
    requestBody := map[string]string{"email": email}
    
    if err := s.client.Post(ctx, "/email/validate", requestBody, &validation); err != nil {
        s.logger.Error("Failed to validate email", "email", email, "error", err)
        return nil, fmt.Errorf("failed to validate email: %w", err)
    }
    
    return &validation, nil
}
```

## 🎯 Business Layer Best Practices

### 1. **Service Organization**
- One service per domain entity
- Services coordinate between multiple repositories
- Keep services focused on business logic, not data manipulation

### 2. **Error Handling**
- Use custom error types for different business scenarios
- Log errors with sufficient context
- Don't expose internal errors to API layer

### 3. **Transaction Management**
```go
func (s *[entity]Service) ComplexOperation(ctx context.Context, req *Request) error {
    return s.db.WithTransaction(ctx, func(tx *sql.Tx) error {
        // Multiple operations in a single transaction
        if err := s.repo1.Create(ctx, tx, data1); err != nil {
            return err
        }
        
        if err := s.repo2.Update(ctx, tx, data2); err != nil {
            return err
        }
        
        return nil
    })
}
```

### 4. **Context Usage**
- Always pass context through service methods
- Extract user information from context for authorization
- Use context for cancellation and timeouts

### 5. **Business Events**
```go
// Optional: Emit business events for complex workflows
func (s *[entity]Service) Create[Entity](ctx context.Context, req *Request) (*models.[Entity], error) {
    [entity], err := s.createLogic(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // Emit event for other services to react
    s.eventBus.Emit(&events.[Entity]Created{
        [Entity]ID: [entity].ID,
        UserID:     [entity].UserID,
        CreatedAt:  [entity].CreatedAt,
    })
    
    return [entity], nil
}
```

---

**Previous**: [← API Layer](./02-api-layer.md) | **Next**: [Data Layer →](./04-data-layer.md)

**Last Updated:** [YYYY-MM-DD]
# Business Layer Architecture

The business layer (`internal/business/`) contains the core business logic, rules, and workflows of the application. It acts as the brain of the system, coordinating between the API layer and data layer.

## 🎯 Layer Responsibilities

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

### Service Interface Pattern

```go
// internal/business/services/user_service.go
package services

import (
    "context"
    
    "github.com/yantology/golang_template/internal/data/models"
)

// UserService interface defines business operations
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error)
    GetUserByID(ctx context.Context, id int64) (*models.User, error)
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
    ListUsers(ctx context.Context, params *ListUsersParams) ([]*models.User, *Pagination, error)
    UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) (*models.User, error)
    DeleteUser(ctx context.Context, id int64) error
    ActivateUser(ctx context.Context, id int64) error
    DeactivateUser(ctx context.Context, id int64) error
    ChangePassword(ctx context.Context, userID int64, req *ChangePasswordRequest) error
}
```

### Service Implementation

```go
// userService implements UserService interface
type userService struct {
    userRepo     repositories.UserRepository
    authService  AuthService
    emailService EmailService
    validator    *UserValidator
    logger       Logger
}

// NewUserService creates a new user service instance
func NewUserService(
    userRepo repositories.UserRepository,
    authService AuthService,
    emailService EmailService,
    logger Logger,
) UserService {
    return &userService{
        userRepo:     userRepo,
        authService:  authService,
        emailService: emailService,
        validator:    NewUserValidator(),
        logger:       logger,
    }
}

// CreateUser creates a new user with business validation
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
    // Business validation
    if err := s.validator.ValidateCreateUser(ctx, req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Check if user already exists
    existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err == nil && existingUser != nil {
        return nil, errors.NewConflictError("user with this email already exists")
    }

    // Hash password
    hashedPassword, err := s.authService.HashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }

    // Create user model
    user := &models.User{
        Email:        req.Email,
        FirstName:    req.FirstName,
        LastName:     req.LastName,
        PasswordHash: hashedPassword,
        Role:         "user", // Default role
        Status:       "active",
    }

    // Set defaults and validate
    user.SetDefaults()
    if err := user.Validate(); err != nil {
        return nil, fmt.Errorf("model validation failed: %w", err)
    }

    // Save to repository
    createdUser, err := s.userRepo.Create(ctx, user)
    if err != nil {
        s.logger.Error("Failed to create user", "error", err, "email", req.Email)
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    // Send welcome email (async)
    go func() {
        if err := s.emailService.SendWelcomeEmail(createdUser.Email, createdUser.FirstName); err != nil {
            s.logger.Error("Failed to send welcome email", "error", err, "user_id", createdUser.ID)
        }
    }()

    // Log success
    s.logger.Info("User created successfully", 
        "user_id", createdUser.ID, 
        "email", createdUser.Email)

    return createdUser, nil
}

// GetUserByID retrieves a user by ID with business logic
func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    // Apply business rules for visibility
    if err := s.checkViewPermission(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

// ListUsers retrieves users with pagination and filtering
func (s *userService) ListUsers(ctx context.Context, params *ListUsersParams) ([]*models.User, *Pagination, error) {
    // Set defaults
    params.SetDefaults()

    // Validate parameters
    if err := s.validator.ValidateListUsersParams(params); err != nil {
        return nil, nil, fmt.Errorf("invalid parameters: %w", err)
    }

    // Check permissions
    if err := s.checkListPermission(ctx); err != nil {
        return nil, nil, err
    }

    // Get users and total count
    users, total, err := s.userRepo.List(ctx, params)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to list users: %w", err)
    }

    // Calculate pagination
    pagination := &Pagination{
        Page:       params.Page,
        Limit:      params.Limit,
        Total:      total,
        TotalPages: (total + int64(params.Limit) - 1) / int64(params.Limit),
    }

    return users, pagination, nil
}

// UpdateUser updates an existing user with business validation
func (s *userService) UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) (*models.User, error) {
    // Get existing user
    existingUser, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    // Check update permission
    if err := s.checkUpdatePermission(ctx, existingUser); err != nil {
        return nil, err
    }

    // Validate update request
    if err := s.validator.ValidateUpdateUser(ctx, req, existingUser); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Apply updates
    updatedUser := s.applyUpdates(existingUser, req)

    // Validate updated model
    if err := updatedUser.Validate(); err != nil {
        return nil, fmt.Errorf("updated model validation failed: %w", err)
    }

    // Save updates
    result, err := s.userRepo.Update(ctx, updatedUser)
    if err != nil {
        s.logger.Error("Failed to update user", "error", err, "user_id", id)
        return nil, fmt.Errorf("failed to update user: %w", err)
    }

    s.logger.Info("User updated successfully", "user_id", id)
    return result, nil
}

// ChangePassword changes a user's password
func (s *userService) ChangePassword(ctx context.Context, userID int64, req *ChangePasswordRequest) error {
    // Get user
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }

    // Verify current password
    if !s.authService.CheckPassword(req.CurrentPassword, user.PasswordHash) {
        return errors.NewValidationError("current password is incorrect")
    }

    // Validate new password
    if err := s.validator.ValidatePassword(req.NewPassword); err != nil {
        return fmt.Errorf("password validation failed: %w", err)
    }

    // Hash new password
    hashedPassword, err := s.authService.HashPassword(req.NewPassword)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    // Update password
    user.PasswordHash = hashedPassword
    if _, err := s.userRepo.Update(ctx, user); err != nil {
        return fmt.Errorf("failed to update password: %w", err)
    }

    s.logger.Info("Password changed successfully", "user_id", userID)
    return nil
}

// Private helper methods

// checkViewPermission checks if the current user can view the target user
func (s *userService) checkViewPermission(ctx context.Context, targetUser *models.User) error {
    currentUserID, ok := ctx.Value("user_id").(int64)
    if !ok {
        return errors.NewUnauthorizedError("user not authenticated")
    }

    // Users can view their own profile
    if targetUser.ID == currentUserID {
        return nil
    }

    // Admins can view all users
    userRole, ok := ctx.Value("user_role").(string)
    if ok && userRole == "admin" {
        return nil
    }

    return errors.NewForbiddenError("not authorized to view this user")
}

// checkUpdatePermission checks if the user can update the target user
func (s *userService) checkUpdatePermission(ctx context.Context, targetUser *models.User) error {
    currentUserID, ok := ctx.Value("user_id").(int64)
    if !ok {
        return errors.NewUnauthorizedError("user not authenticated")
    }

    // Users can update their own profile
    if targetUser.ID == currentUserID {
        return nil
    }

    // Admins can update any user
    userRole, ok := ctx.Value("user_role").(string)
    if ok && userRole == "admin" {
        return nil
    }

    return errors.NewForbiddenError("not authorized to update this user")
}

// checkListPermission checks if the user can list users
func (s *userService) checkListPermission(ctx context.Context) error {
    userRole, ok := ctx.Value("user_role").(string)
    if ok && userRole == "admin" {
        return nil
    }

    return errors.NewForbiddenError("admin access required to list users")
}

// applyUpdates applies the update request to the existing user
func (s *userService) applyUpdates(existing *models.User, req *UpdateUserRequest) *models.User {
    updated := *existing // Copy the struct

    if req.FirstName != nil {
        updated.FirstName = *req.FirstName
    }
    if req.LastName != nil {
        updated.LastName = *req.LastName
    }
    if req.Email != nil {
        updated.Email = *req.Email
    }

    return &updated
}
```

## 🔍 Validators

Validators implement business rule validation that goes beyond basic input validation.

### Validator Implementation

```go
// internal/business/validators/user_validator.go
package validators

import (
    "context"
    "fmt"
    "regexp"
    "strings"
    "unicode"
    
    "github.com/yantology/golang_template/internal/data/models"
    "github.com/yantology/golang_template/pkg/errors"
)

type UserValidator struct {
    // Add dependencies like repositories if needed for validation
}

func NewUserValidator() *UserValidator {
    return &UserValidator{}
}

// ValidateCreateUser validates business rules for creating a user
func (v *UserValidator) ValidateCreateUser(ctx context.Context, req *CreateUserRequest) error {
    if err := v.validateEmail(req.Email); err != nil {
        return err
    }

    if err := v.validatePassword(req.Password); err != nil {
        return err
    }

    if err := v.validateName(req.FirstName, "first name"); err != nil {
        return err
    }

    if err := v.validateName(req.LastName, "last name"); err != nil {
        return err
    }

    return nil
}

// ValidateUpdateUser validates business rules for updating a user
func (v *UserValidator) ValidateUpdateUser(ctx context.Context, req *UpdateUserRequest, existing *models.User) error {
    if req.Email != nil {
        if err := v.validateEmail(*req.Email); err != nil {
            return err
        }
    }

    if req.FirstName != nil {
        if err := v.validateName(*req.FirstName, "first name"); err != nil {
            return err
        }
    }

    if req.LastName != nil {
        if err := v.validateName(*req.LastName, "last name"); err != nil {
            return err
        }
    }

    return nil
}

// ValidateListUsersParams validates parameters for listing users
func (v *UserValidator) ValidateListUsersParams(params *ListUsersParams) error {
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
        "name":       true,
        "email":      true,
        "created_at": true,
        "updated_at": true,
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

// ValidatePassword validates password business rules
func (v *UserValidator) ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.NewValidationError("password must be at least 8 characters long")
    }

    if len(password) > 128 {
        return errors.NewValidationError("password must not exceed 128 characters")
    }

    var hasUpper, hasLower, hasDigit, hasSpecial bool
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasDigit = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }

    if !hasUpper {
        return errors.NewValidationError("password must contain at least one uppercase letter")
    }
    if !hasLower {
        return errors.NewValidationError("password must contain at least one lowercase letter")
    }
    if !hasDigit {
        return errors.NewValidationError("password must contain at least one digit")
    }
    if !hasSpecial {
        return errors.NewValidationError("password must contain at least one special character")
    }

    return nil
}

// Private validation methods

func (v *UserValidator) validateEmail(email string) error {
    if email == "" {
        return errors.NewValidationError("email is required")
    }

    if len(email) > 255 {
        return errors.NewValidationError("email must not exceed 255 characters")
    }

    // Basic email validation
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(email) {
        return errors.NewValidationError("invalid email format")
    }

    // Business rule: no disposable email providers
    disposableProviders := []string{"10minutemail.com", "tempmail.org", "guerrillamail.com"}
    emailLower := strings.ToLower(email)
    for _, provider := range disposableProviders {
        if strings.HasSuffix(emailLower, "@"+provider) {
            return errors.NewValidationError("disposable email addresses are not allowed")
        }
    }

    return nil
}

func (v *UserValidator) validateName(name, fieldName string) error {
    if name == "" {
        return errors.NewValidationError(fmt.Sprintf("%s is required", fieldName))
    }

    if len(name) < 2 {
        return errors.NewValidationError(fmt.Sprintf("%s must be at least 2 characters", fieldName))
    }

    if len(name) > 100 {
        return errors.NewValidationError(fmt.Sprintf("%s must not exceed 100 characters", fieldName))
    }

    // Business rule: no special characters in names
    nameRegex := regexp.MustCompile(`^[a-zA-Z\s'-]+$`)
    if !nameRegex.MatchString(name) {
        return errors.NewValidationError(fmt.Sprintf("%s can only contain letters, spaces, apostrophes, and hyphens", fieldName))
    }

    return nil
}
```

## 📊 Supporting Types

### Request Types

```go
// internal/business/services/types.go
package services

import "time"

// Request types for service operations
type CreateUserRequest struct {
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Password  string `json:"password"`
}

type UpdateUserRequest struct {
    FirstName *string `json:"first_name,omitempty"`
    LastName  *string `json:"last_name,omitempty"`
    Email     *string `json:"email,omitempty"`
}

type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password"`
    NewPassword     string `json:"new_password"`
}

type ListUsersParams struct {
    Page      int    `form:"page"`
    Limit     int    `form:"limit"`
    Search    string `form:"search"`
    Role      string `form:"role"`
    Status    string `form:"status"`
    SortBy    string `form:"sort_by"`
    SortOrder string `form:"sort_order"`
}

func (p *ListUsersParams) SetDefaults() {
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
    SendAccountActivationEmail(userEmail, activationLink string) error
}

type AuthService interface {
    HashPassword(password string) (string, error)
    CheckPassword(password, hash string) bool
    GenerateJWT(userID int64, email, role string) (string, error)
    ValidateJWT(token string) (*JWTClaims, error)
}

type JWTClaims struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
}
```

## 🔄 Transaction Management

### Transaction Service Pattern

```go
// internal/business/services/transaction_service.go
package services

import (
    "context"
    "database/sql"
    
    "github.com/yantology/golang_template/internal/data/repositories"
)

type TransactionService struct {
    db       *sql.DB
    userRepo repositories.UserRepository
    logger   Logger
}

func NewTransactionService(db *sql.DB, userRepo repositories.UserRepository, logger Logger) *TransactionService {
    return &TransactionService{
        db:       db,
        userRepo: userRepo,
        logger:   logger,
    }
}

// ComplexUserOperation demonstrates transaction management
func (s *TransactionService) ComplexUserOperation(ctx context.Context, req *ComplexRequest) error {
    // Begin transaction
    tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelReadCommitted,
    })
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()

    // Perform multiple operations within transaction
    if err := s.performOperation1(ctx, tx, req); err != nil {
        tx.Rollback()
        return fmt.Errorf("operation 1 failed: %w", err)
    }

    if err := s.performOperation2(ctx, tx, req); err != nil {
        tx.Rollback()
        return fmt.Errorf("operation 2 failed: %w", err)
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    s.logger.Info("Complex operation completed successfully")
    return nil
}
```

## 🎯 Business Layer Best Practices

### 1. Service Organization
- One service per domain entity
- Services coordinate between multiple repositories
- Keep services focused on business logic, not data manipulation

### 2. Error Handling
- Use custom error types for different business scenarios
- Log errors with sufficient context
- Don't expose internal errors to API layer

### 3. Context Usage
- Always pass context through service methods
- Extract user information from context for authorization
- Use context for cancellation and timeouts

### 4. Dependency Injection
- Use interfaces for all dependencies
- Make dependencies explicit in constructor
- Enable easy testing and mocking

### 5. Business Events (Optional)
```go
// Optional: Emit business events for complex workflows
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
    user, err := s.createUserLogic(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // Emit event for other services to react
    s.eventBus.Emit(&events.UserCreated{
        UserID:    user.ID,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
    })
    
    return user, nil
}
```

## 🚀 Next Steps

- **Understand data layer**: [Data Layer](./data-layer.md)
- **Learn dependency injection**: [Dependency Injection](./dependency-injection.md)
- **See practical examples**: [Examples](../examples/)

---

The business layer encapsulates your domain logic and business rules, ensuring they remain independent of external frameworks and infrastructure concerns.
# Data Layer Architecture

The data layer (`internal/data/`) handles data persistence, database operations, and data modeling. It provides an abstraction between the business layer and the underlying database.

## 🎯 Layer Responsibilities

- **Data Modeling**: Define database entity structures and relationships
- **Data Access**: Implement repository patterns for database operations
- **Query Building**: Construct optimized database queries
- **Transaction Management**: Handle database transactions
- **Data Validation**: Ensure data integrity at the database level
- **Schema Management**: Database migrations and schema versioning

## 📂 Data Layer Structure

```
internal/data/
├── models/             # Database entity models
├── repositories/       # Data access layer implementation
└── migrations/         # Database migration files
```

## 🗂️ Models

Models define the structure of database entities and their relationships.

### Base Model

```go
// internal/data/models/base.go
package models

import (
    "time"
)

// BaseModel contains common fields for all entities
type BaseModel struct {
    ID        int64     `json:"id" db:"id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SoftDeleteModel extends BaseModel with soft delete functionality
type SoftDeleteModel struct {
    BaseModel
    DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// IsDeleted returns true if the entity is soft deleted
func (m *SoftDeleteModel) IsDeleted() bool {
    return m.DeletedAt != nil
}

// TimestampModel provides automatic timestamp management
type TimestampModel struct {
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Touch updates the UpdatedAt timestamp
func (m *TimestampModel) Touch() {
    m.UpdatedAt = time.Now()
}
```

### User Model Example

```go
// internal/data/models/user.go
package models

import (
    "fmt"
    "regexp"
    "strings"
    "time"
)

type User struct {
    ID               int64      `json:"id" db:"id"`
    Email            string     `json:"email" db:"email"`
    PasswordHash     string     `json:"-" db:"password_hash"`
    FirstName        string     `json:"first_name" db:"first_name"`
    LastName         string     `json:"last_name" db:"last_name"`
    Role             string     `json:"role" db:"role"`
    Status           string     `json:"status" db:"status"`
    EmailVerifiedAt  *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
    LastLoginAt      *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
    CreatedAt        time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// SetDefaults sets default values for new users
func (u *User) SetDefaults() {
    if u.Role == "" {
        u.Role = "user"
    }
    if u.Status == "" {
        u.Status = "active"
    }
}

// Validate performs basic validation on the user
func (u *User) Validate() error {
    if u.Email == "" {
        return fmt.Errorf("email is required")
    }
    if u.FirstName == "" {
        return fmt.Errorf("first name is required")
    }
    if u.LastName == "" {
        return fmt.Errorf("last name is required")
    }
    if u.PasswordHash == "" {
        return fmt.Errorf("password hash is required")
    }
    
    validRoles := map[string]bool{
        "user":  true,
        "admin": true,
    }
    if !validRoles[u.Role] {
        return fmt.Errorf("invalid role: %s", u.Role)
    }
    
    validStatuses := map[string]bool{
        "active":   true,
        "inactive": true,
        "suspended": true,
    }
    if !validStatuses[u.Status] {
        return fmt.Errorf("invalid status: %s", u.Status)
    }
    
    // Validate email format
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(u.Email) {
        return fmt.Errorf("invalid email format")
    }
    
    return nil
}

// IsActive returns true if the user is active
func (u *User) IsActive() bool {
    return u.Status == "active"
}

// IsAdmin returns true if the user has admin role
func (u *User) IsAdmin() bool {
    return u.Role == "admin"
}

// IsEmailVerified returns true if the user's email is verified
func (u *User) IsEmailVerified() bool {
    return u.EmailVerifiedAt != nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
    return strings.TrimSpace(u.FirstName + " " + u.LastName)
}
```

### Article Model Example

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
    PublishedAt   *time.Time `json:"published_at,omitempty" db:"published_at"`
    CreatedAt     time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
    
    // Related entities (loaded via JOIN queries)
    User     User     `json:"user,omitempty"`
    Category Category `json:"category,omitempty"`
}

type Category struct {
    ID          int64  `json:"id" db:"id"`
    Name        string `json:"name" db:"name"`
    Slug        string `json:"slug" db:"slug"`
    Description string `json:"description,omitempty" db:"description"`
    ParentID    *int64 `json:"parent_id,omitempty" db:"parent_id"`
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

// Validate performs basic validation on the article
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
        "draft":     true,
        "published": true,
        "archived":  true,
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

// CanBeViewed checks if the article can be viewed by public
func (a *Article) CanBeViewed() bool {
    return a.Status == "published" && a.PublishedAt != nil
}

// Helper functions

// generateSlug creates a URL-friendly slug from title
func generateSlug(title string) string {
    // Basic slug generation - in production, use a proper slug library
    // like github.com/gosimple/slug
    slug := strings.ToLower(title)
    slug = strings.ReplaceAll(slug, " ", "-")
    slug = strings.ReplaceAll(slug, "'", "")
    // Remove special characters and keep only alphanumeric and hyphens
    var result strings.Builder
    for _, r := range slug {
        if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
            result.WriteRune(r)
        }
    }
    return strings.Trim(result.String(), "-")
}

// generateExcerpt creates an excerpt from content
func generateExcerpt(content string, maxLength int) string {
    if len(content) <= maxLength {
        return content
    }
    
    // Find the last space before maxLength to avoid cutting words
    excerpt := content[:maxLength]
    if lastSpace := strings.LastIndex(excerpt, " "); lastSpace > 0 {
        excerpt = content[:lastSpace]
    }
    
    return excerpt + "..."
}
```

## 🏪 Repositories

Repositories implement the data access layer using the Repository pattern.

### Repository Interfaces

```go
// internal/data/repositories/interfaces.go
package repositories

import (
    "context"
    
    "github.com/yantology/golang_template/internal/data/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
    Create(ctx context.Context, user *models.User) (*models.User, error)
    GetByID(ctx context.Context, id int64) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    List(ctx context.Context, params *ListUsersParams) ([]*models.User, int64, error)
    Update(ctx context.Context, user *models.User) (*models.User, error)
    Delete(ctx context.Context, id int64) error
    UpdateLastLogin(ctx context.Context, userID int64) error
}

// ArticleRepository defines the interface for article data operations
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

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
    Create(ctx context.Context, category *models.Category) (*models.Category, error)
    GetByID(ctx context.Context, id int64) (*models.Category, error)
    GetBySlug(ctx context.Context, slug string) (*models.Category, error)
    List(ctx context.Context) ([]*models.Category, error)
    GetSubCategories(ctx context.Context, parentID int64) ([]*models.Category, error)
    Update(ctx context.Context, category *models.Category) (*models.Category, error)
    Delete(ctx context.Context, id int64) error
}

// Query parameter types
type ListUsersParams struct {
    Page      int
    Limit     int
    Search    string
    Role      string
    Status    string
    SortBy    string
    SortOrder string
}

type ListArticlesParams struct {
    Page         int
    Limit        int
    Search       string
    CategoryID   int64
    Status       string
    SortBy       string
    SortOrder    string
    UserID       int64 // For filtering by user
}
```

### Repository Implementation

```go
// internal/data/repositories/user_repository.go
package repositories

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    
    "github.com/yantology/golang_template/internal/data/models"
    
    "github.com/Masterminds/squirrel"
)

type userRepository struct {
    db *sql.DB
    qb squirrel.StatementBuilderType
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{
        db: db,
        qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar), // PostgreSQL style
    }
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
    query := r.qb.Insert("users").
        Columns("email", "password_hash", "first_name", "last_name", "role", "status").
        Values(user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Role, user.Status).
        Suffix("RETURNING id, created_at, updated_at")
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build insert query: %w", err)
    }
    
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
    query := r.qb.Select(
        "id", "email", "password_hash", "first_name", "last_name", "role", "status",
        "email_verified_at", "last_login_at", "created_at", "updated_at",
    ).
        From("users").
        Where(squirrel.Eq{"id": id})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build select query: %w", err)
    }
    
    user := &models.User{}
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(
        &user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
        &user.Role, &user.Status, &user.EmailVerifiedAt, &user.LastLoginAt,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    query := r.qb.Select(
        "id", "email", "password_hash", "first_name", "last_name", "role", "status",
        "email_verified_at", "last_login_at", "created_at", "updated_at",
    ).
        From("users").
        Where(squirrel.Eq{"email": email})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build select query: %w", err)
    }
    
    user := &models.User{}
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(
        &user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
        &user.Role, &user.Status, &user.EmailVerifiedAt, &user.LastLoginAt,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}

func (r *userRepository) List(ctx context.Context, params *ListUsersParams) ([]*models.User, int64, error) {
    // Build base query for counting
    baseQuery := r.qb.Select("COUNT(*)").From("users")
    
    // Build main select query
    selectQuery := r.qb.Select(
        "id", "email", "first_name", "last_name", "role", "status",
        "email_verified_at", "last_login_at", "created_at", "updated_at",
    ).From("users")
    
    // Apply filters
    conditions := r.buildFilterConditions(params)
    if len(conditions) > 0 {
        baseQuery = baseQuery.Where(conditions)
        selectQuery = selectQuery.Where(conditions)
    }
    
    // Count total records
    total, err := r.executeCountQuery(ctx, baseQuery)
    if err != nil {
        return nil, 0, err
    }
    
    // Apply sorting
    selectQuery = r.applySorting(selectQuery, params)
    
    // Apply pagination
    if params.Limit > 0 {
        offset := (params.Page - 1) * params.Limit
        selectQuery = selectQuery.Limit(uint64(params.Limit)).Offset(uint64(offset))
    }
    
    // Execute main query
    users, err := r.executeSelectQuery(ctx, selectQuery)
    if err != nil {
        return nil, 0, err
    }
    
    return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
    query := r.qb.Update("users").
        Set("email", user.Email).
        Set("first_name", user.FirstName).
        Set("last_name", user.LastName).
        Set("role", user.Role).
        Set("status", user.Status).
        Set("password_hash", user.PasswordHash).
        Set("email_verified_at", user.EmailVerifiedAt).
        Set("updated_at", "NOW()").
        Where(squirrel.Eq{"id": user.ID}).
        Suffix("RETURNING updated_at")
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build update query: %w", err)
    }
    
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&user.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to update user: %w", err)
    }
    
    return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
    query := r.qb.Delete("users").Where(squirrel.Eq{"id": id})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return fmt.Errorf("failed to build delete query: %w", err)
    }
    
    result, err := r.db.ExecContext(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }
    
    return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
    query := r.qb.Update("users").
        Set("last_login_at", "NOW()").
        Where(squirrel.Eq{"id": userID})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return fmt.Errorf("failed to build update query: %w", err)
    }
    
    _, err = r.db.ExecContext(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("failed to update last login: %w", err)
    }
    
    return nil
}

// Private helper methods

func (r *userRepository) buildFilterConditions(params *ListUsersParams) squirrel.And {
    var conditions squirrel.And
    
    if params.Role != "" {
        conditions = append(conditions, squirrel.Eq{"role": params.Role})
    }
    
    if params.Status != "" {
        conditions = append(conditions, squirrel.Eq{"status": params.Status})
    }
    
    if params.Search != "" {
        searchPattern := "%" + params.Search + "%"
        conditions = append(conditions, squirrel.Or{
            squirrel.ILike{"first_name": searchPattern},
            squirrel.ILike{"last_name": searchPattern},
            squirrel.ILike{"email": searchPattern},
        })
    }
    
    return conditions
}

func (r *userRepository) applySorting(query squirrel.SelectBuilder, params *ListUsersParams) squirrel.SelectBuilder {
    validSortFields := map[string]string{
        "name":       "first_name",
        "email":      "email",
        "created_at": "created_at",
        "updated_at": "updated_at",
    }
    
    sortField, ok := validSortFields[params.SortBy]
    if !ok {
        sortField = "created_at" // Default sort
    }
    
    sortOrder := "DESC"
    if strings.ToUpper(params.SortOrder) == "ASC" {
        sortOrder = "ASC"
    }
    
    return query.OrderBy(fmt.Sprintf("%s %s", sortField, sortOrder))
}

func (r *userRepository) executeCountQuery(ctx context.Context, query squirrel.SelectBuilder) (int64, error) {
    sql, args, err := query.ToSql()
    if err != nil {
        return 0, fmt.Errorf("failed to build count query: %w", err)
    }
    
    var total int64
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&total)
    if err != nil {
        return 0, fmt.Errorf("failed to count users: %w", err)
    }
    
    return total, nil
}

func (r *userRepository) executeSelectQuery(ctx context.Context, query squirrel.SelectBuilder) ([]*models.User, error) {
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build select query: %w", err)
    }
    
    rows, err := r.db.QueryContext(ctx, sql, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %w", err)
    }
    defer rows.Close()
    
    var users []*models.User
    for rows.Next() {
        user, err := r.scanRowToUser(rows)
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }
        users = append(users, user)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }
    
    return users, nil
}

func (r *userRepository) scanRowToUser(rows *sql.Rows) (*models.User, error) {
    user := &models.User{}
    
    err := rows.Scan(
        &user.ID, &user.Email, &user.FirstName, &user.LastName,
        &user.Role, &user.Status, &user.EmailVerifiedAt, &user.LastLoginAt,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

## 🗃️ Database Migrations

### Migration Structure

```
internal/data/migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_categories_table.up.sql
├── 000002_create_categories_table.down.sql
├── 000003_create_articles_table.up.sql
└── 000003_create_articles_table.down.sql
```

### Example Migrations

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) DEFAULT 'user' NOT NULL,
    status VARCHAR(50) DEFAULT 'active' NOT NULL,
    email_verified_at TIMESTAMP,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Function to update updated_at automatically
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for users table
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

```sql
-- 000001_create_users_table.down.sql
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS users;
```

```sql
-- 000003_create_articles_table.up.sql
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    parent_id BIGINT REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE articles (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(500) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    featured_image VARCHAR(500),
    status VARCHAR(50) DEFAULT 'draft' NOT NULL,
    view_count BIGINT DEFAULT 0 NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- Indexes
CREATE INDEX idx_articles_slug ON articles(slug);
CREATE INDEX idx_articles_status ON articles(status);
CREATE INDEX idx_articles_user_id ON articles(user_id);
CREATE INDEX idx_articles_category_id ON articles(category_id);
CREATE INDEX idx_articles_published_at ON articles(published_at);
CREATE INDEX idx_articles_created_at ON articles(created_at);

-- Full-text search index
CREATE INDEX idx_articles_search ON articles USING gin(to_tsvector('english', title || ' ' || content));

-- Triggers
CREATE TRIGGER update_articles_updated_at BEFORE UPDATE ON articles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### Migration Management

```bash
# Create a new migration
make migrate-create NAME=add_user_profiles

# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down

# Check migration status
make migrate-status

# Force migration version
make migrate-force VERSION=1
```

## 🎯 Data Layer Best Practices

### 1. Repository Pattern
- One repository per aggregate root
- Keep repositories focused on data access, not business logic
- Use interfaces to allow for easy testing and swapping implementations

### 2. Query Optimization
```go
// Use JOINs to avoid N+1 queries
func (r *articleRepository) GetWithRelations(ctx context.Context, id int64) (*models.Article, error) {
    // Single query with JOINs instead of multiple queries
    query := r.qb.Select("a.*, u.first_name, u.last_name, c.name").
        From("articles a").
        LeftJoin("users u ON u.id = a.user_id").
        LeftJoin("categories c ON c.id = a.category_id").
        Where(squirrel.Eq{"a.id": id})
    // ...
}
```

### 3. Error Handling
```go
// Distinguish between different types of database errors
func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
    // ...
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&user.ID)
    if err != nil {
        // Check for constraint violations, duplicate keys, etc.
        if isDuplicateKeyError(err) {
            return nil, fmt.Errorf("user with this email already exists")
        }
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    // ...
}
```

### 4. Context Usage
- Always use context for cancellation and timeouts
- Pass context to all database operations
- Handle context cancellation gracefully

### 5. Connection Management
```go
// Configure connection pool properly
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

## 🚀 Next Steps

- **Learn dependency injection**: [Dependency Injection](./dependency-injection.md)
- **See practical examples**: [Database Guide](../database/)
- **Understand complete workflow**: [Examples](../examples/)

---

The data layer provides a clean abstraction over your database operations while maintaining performance and data integrity through proper modeling and query optimization.
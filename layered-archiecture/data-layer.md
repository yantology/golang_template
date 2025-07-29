# 💾 Data Layer (`internal/data/`)

The data layer handles data persistence, database operations, and data modeling. It provides an abstraction between the business layer and the underlying database.

## 📋 Layer Responsibilities

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

### 📄 models/[entity].go

```go
package models

import (
    "fmt"
    "strings"
    "time"
)

type [Entity] struct {
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

type User struct {
    ID        int64  `json:"id" db:"id"`
    Email     string `json:"email" db:"email"`
    FirstName string `json:"first_name" db:"first_name"`
    LastName  string `json:"last_name" db:"last_name"`
    Role      string `json:"role" db:"role"`
}

type Category struct {
    ID          int64  `json:"id" db:"id"`
    Name        string `json:"name" db:"name"`
    Slug        string `json:"slug" db:"slug"`
    Description string `json:"description,omitempty" db:"description"`
    ParentID    *int64 `json:"parent_id,omitempty" db:"parent_id"`
}

type SubCategory struct {
    ID       int64  `json:"id" db:"id"`
    Name     string `json:"name" db:"name"`
    Slug     string `json:"slug" db:"slug"`
    ParentID int64  `json:"parent_id" db:"parent_id"`
}

// GenerateSlug creates a URL-friendly slug from title
func (e *[Entity]) GenerateSlug() {
    e.Slug = generateSlug(e.Title)
}

// SetDefaults sets default values for new entities
func (e *[Entity]) SetDefaults() {
    if e.Status == "" {
        e.Status = "draft"
    }
    if e.ViewCount == 0 {
        e.ViewCount = 0
    }
    if e.Excerpt == "" && e.Content != "" {
        e.Excerpt = generateExcerpt(e.Content, 160)
    }
}

// Validate performs basic validation on the entity
func (e *[Entity]) Validate() error {
    if e.Title == "" {
        return fmt.Errorf("title is required")
    }
    if e.UserID == 0 {
        return fmt.Errorf("user_id is required")
    }
    if e.CategoryID == 0 {
        return fmt.Errorf("category_id is required")
    }
    
    validStatuses := map[string]bool{
        "draft":     true,
        "published": true,
        "archived":  true,
    }
    if !validStatuses[e.Status] {
        return fmt.Errorf("invalid status: %s", e.Status)
    }
    
    return nil
}

// IsPublished returns true if the entity is published
func (e *[Entity]) IsPublished() bool {
    return e.Status == "published"
}

// CanBeViewed checks if the entity can be viewed by public
func (e *[Entity]) CanBeViewed() bool {
    return e.Status == "published" && e.PublishedAt != nil
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

### 📄 models/base.go

```go
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

## 🏪 Repositories

Repositories implement the data access layer using the Repository pattern.

### 📄 repositories/interfaces.go

```go
package repositories

import (
    "context"
    
    "[module-name]/internal/data/models"
)

// [Entity]Repository defines the interface for [entity] data operations
type [Entity]Repository interface {
    Create(ctx context.Context, [entity] *models.[Entity]) (*models.[Entity], error)
    GetByID(ctx context.Context, id int64) (*models.[Entity], error)
    GetBySlug(ctx context.Context, slug string) (*models.[Entity], error)
    List(ctx context.Context, params *List[Entities]Params) ([]*models.[Entity], int64, error)
    Update(ctx context.Context, [entity] *models.[Entity]) (*models.[Entity], error)
    Delete(ctx context.Context, id int64) error
    IncrementViewCount(ctx context.Context, id int64) error
    GetByUserID(ctx context.Context, userID int64, params *List[Entities]Params) ([]*models.[Entity], int64, error)
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
    Create(ctx context.Context, user *models.User) (*models.User, error)
    GetByID(ctx context.Context, id int64) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) (*models.User, error)
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, params *ListUsersParams) ([]*models.User, int64, error)
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
type List[Entities]Params struct {
    Page         int
    Limit        int
    Search       string
    CategoryID   int64
    Status       string
    SortBy       string
    SortOrder    string
    UserID       int64 // For filtering by user
}

type ListUsersParams struct {
    Page      int
    Limit     int
    Search    string
    Role      string
    SortBy    string
    SortOrder string
}
```

### 📄 repositories/[entity]_repository.go

```go
package repositories

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    
    "[module-name]/internal/data/models"
    
    "github.com/Masterminds/squirrel"
)

type [entity]Repository struct {
    db *sql.DB
    qb squirrel.StatementBuilderType
}

func New[Entity]Repository(db *sql.DB) [Entity]Repository {
    return &[entity]Repository{
        db: db,
        qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar), // PostgreSQL style
    }
}

func (r *[entity]Repository) Create(ctx context.Context, [entity] *models.[Entity]) (*models.[Entity], error) {
    query := r.qb.Insert("[entities]").
        Columns("title", "slug", "content", "excerpt", "featured_image", "status", "user_id", "category_id", "sub_category_id").
        Values([entity].Title, [entity].Slug, [entity].Content, [entity].Excerpt, [entity].FeaturedImage, [entity].Status, [entity].UserID, [entity].CategoryID, [entity].SubCategoryID).
        Suffix("RETURNING id, created_at, updated_at")
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build insert query: %w", err)
    }
    
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&[entity].ID, &[entity].CreatedAt, &[entity].UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to create [entity]: %w", err)
    }
    
    return [entity], nil
}

func (r *[entity]Repository) GetByID(ctx context.Context, id int64) (*models.[Entity], error) {
    query := r.qb.Select(
        "e.id", "e.title", "e.slug", "e.content", "e.excerpt", "e.featured_image", 
        "e.status", "e.view_count", "e.user_id", "e.category_id", "e.sub_category_id", 
        "e.published_at", "e.created_at", "e.updated_at",
        "u.id", "u.email", "u.first_name", "u.last_name", "u.role",
        "c.id", "c.name", "c.slug", "c.description",
        "sc.id", "sc.name", "sc.slug",
    ).
        From("[entities] e").
        LeftJoin("users u ON u.id = e.user_id").
        LeftJoin("categories c ON c.id = e.category_id").
        LeftJoin("categories sc ON sc.id = e.sub_category_id").
        Where(squirrel.Eq{"e.id": id})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build select query: %w", err)
    }
    
    [entity] := &models.[Entity]{}
    var user models.User
    var category models.Category
    var subCategory models.SubCategory
    var categoryDesc sql.NullString
    
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(
        &[entity].ID, &[entity].Title, &[entity].Slug, &[entity].Content, &[entity].Excerpt, 
        &[entity].FeaturedImage, &[entity].Status, &[entity].ViewCount, &[entity].UserID, 
        &[entity].CategoryID, &[entity].SubCategoryID, &[entity].PublishedAt, 
        &[entity].CreatedAt, &[entity].UpdatedAt,
        &user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Role,
        &category.ID, &category.Name, &category.Slug, &categoryDesc,
        &subCategory.ID, &subCategory.Name, &subCategory.Slug,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("[entity] not found")
        }
        return nil, fmt.Errorf("failed to get [entity]: %w", err)
    }
    
    // Set related entities
    [entity].User = user
    [entity].Category = category
    if categoryDesc.Valid {
        [entity].Category.Description = categoryDesc.String
    }
    [entity].SubCategory = subCategory
    
    return [entity], nil
}

func (r *[entity]Repository) GetBySlug(ctx context.Context, slug string) (*models.[Entity], error) {
    // Similar to GetByID but filter by slug
    query := r.qb.Select(
        "e.id", "e.title", "e.slug", "e.content", "e.excerpt", "e.featured_image", 
        "e.status", "e.view_count", "e.user_id", "e.category_id", "e.sub_category_id", 
        "e.published_at", "e.created_at", "e.updated_at",
        "u.id", "u.email", "u.first_name", "u.last_name", "u.role",
        "c.id", "c.name", "c.slug", "c.description",
        "sc.id", "sc.name", "sc.slug",
    ).
        From("[entities] e").
        LeftJoin("users u ON u.id = e.user_id").
        LeftJoin("categories c ON c.id = e.category_id").
        LeftJoin("categories sc ON sc.id = e.sub_category_id").
        Where(squirrel.Eq{"e.slug": slug})
    
    // Implementation similar to GetByID...
    // [Rest of implementation omitted for brevity]
}

func (r *[entity]Repository) List(ctx context.Context, params *List[Entities]Params) ([]*models.[Entity], int64, error) {
    // Build base query for counting
    baseQuery := r.qb.Select("COUNT(*)").From("[entities] e")
    
    // Build main select query
    selectQuery := r.qb.Select(
        "e.id", "e.title", "e.slug", "e.content", "e.excerpt", "e.featured_image", 
        "e.status", "e.view_count", "e.user_id", "e.category_id", "e.sub_category_id", 
        "e.published_at", "e.created_at", "e.updated_at",
        "u.email", "u.first_name", "u.last_name", "u.role",
        "c.name as category_name", "c.slug as category_slug",
        "sc.name as sub_category_name", "sc.slug as sub_category_slug",
    ).
        From("[entities] e").
        LeftJoin("users u ON u.id = e.user_id").
        LeftJoin("categories c ON c.id = e.category_id").
        LeftJoin("categories sc ON sc.id = e.sub_category_id")
    
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
    [entities], err := r.executeSelectQuery(ctx, selectQuery)
    if err != nil {
        return nil, 0, err
    }
    
    return [entities], total, nil
}

func (r *[entity]Repository) Update(ctx context.Context, [entity] *models.[Entity]) (*models.[Entity], error) {
    query := r.qb.Update("[entities]").
        Set("title", [entity].Title).
        Set("slug", [entity].Slug).
        Set("content", [entity].Content).
        Set("excerpt", [entity].Excerpt).
        Set("featured_image", [entity].FeaturedImage).
        Set("status", [entity].Status).
        Set("category_id", [entity].CategoryID).
        Set("sub_category_id", [entity].SubCategoryID).
        Set("published_at", [entity].PublishedAt).
        Set("updated_at", "NOW()").
        Where(squirrel.Eq{"id": [entity].ID}).
        Suffix("RETURNING updated_at")
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build update query: %w", err)
    }
    
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&[entity].UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to update [entity]: %w", err)
    }
    
    return [entity], nil
}

func (r *[entity]Repository) Delete(ctx context.Context, id int64) error {
    query := r.qb.Delete("[entities]").Where(squirrel.Eq{"id": id})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return fmt.Errorf("failed to build delete query: %w", err)
    }
    
    result, err := r.db.ExecContext(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("failed to delete [entity]: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("[entity] not found")
    }
    
    return nil
}

func (r *[entity]Repository) IncrementViewCount(ctx context.Context, id int64) error {
    query := r.qb.Update("[entities]").
        Set("view_count", squirrel.Expr("view_count + 1")).
        Where(squirrel.Eq{"id": id})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return fmt.Errorf("failed to build update query: %w", err)
    }
    
    _, err = r.db.ExecContext(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("failed to increment view count: %w", err)
    }
    
    return nil
}

func (r *[entity]Repository) GetByUserID(ctx context.Context, userID int64, params *List[Entities]Params) ([]*models.[Entity], int64, error) {
    // Add user_id filter to existing List implementation
    userParams := *params
    userParams.UserID = userID
    return r.List(ctx, &userParams)
}

// Private helper methods

func (r *[entity]Repository) buildFilterConditions(params *List[Entities]Params) squirrel.And {
    var conditions squirrel.And
    
    if params.CategoryID != 0 {
        conditions = append(conditions, squirrel.Eq{"e.category_id": params.CategoryID})
    }
    
    if params.Status != "" {
        conditions = append(conditions, squirrel.Eq{"e.status": params.Status})
    }
    
    if params.UserID != 0 {
        conditions = append(conditions, squirrel.Eq{"e.user_id": params.UserID})
    }
    
    if params.Search != "" {
        searchPattern := "%" + params.Search + "%"
        conditions = append(conditions, squirrel.Or{
            squirrel.ILike{"e.title": searchPattern},
            squirrel.ILike{"e.content": searchPattern},
            squirrel.ILike{"e.excerpt": searchPattern},
        })
    }
    
    return conditions
}

func (r *[entity]Repository) applySorting(query squirrel.SelectBuilder, params *List[Entities]Params) squirrel.SelectBuilder {
    validSortFields := map[string]string{
        "title":      "e.title",
        "created_at": "e.created_at",
        "updated_at": "e.updated_at",
        "view_count": "e.view_count",
    }
    
    sortField, ok := validSortFields[params.SortBy]
    if !ok {
        sortField = "e.created_at" // Default sort
    }
    
    sortOrder := "DESC"
    if strings.ToUpper(params.SortOrder) == "ASC" {
        sortOrder = "ASC"
    }
    
    return query.OrderBy(fmt.Sprintf("%s %s", sortField, sortOrder))
}

func (r *[entity]Repository) executeCountQuery(ctx context.Context, query squirrel.SelectBuilder) (int64, error) {
    sql, args, err := query.ToSql()
    if err != nil {
        return 0, fmt.Errorf("failed to build count query: %w", err)
    }
    
    var total int64
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&total)
    if err != nil {
        return 0, fmt.Errorf("failed to count [entities]: %w", err)
    }
    
    return total, nil
}

func (r *[entity]Repository) executeSelectQuery(ctx context.Context, query squirrel.SelectBuilder) ([]*models.[Entity], error) {
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build select query: %w", err)
    }
    
    rows, err := r.db.QueryContext(ctx, sql, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %w", err)
    }
    defer rows.Close()
    
    var [entities] []*models.[Entity]
    for rows.Next() {
        [entity], err := r.scanRowToEntity(rows)
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }
        [entities] = append([entities], [entity])
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }
    
    return [entities], nil
}

func (r *[entity]Repository) scanRowToEntity(rows *sql.Rows) (*models.[Entity], error) {
    [entity] := &models.[Entity]{}
    var userEmail, userFirstName, userLastName, userRole sql.NullString
    var categoryName, categorySlug, subCategoryName, subCategorySlug sql.NullString
    
    err := rows.Scan(
        &[entity].ID, &[entity].Title, &[entity].Slug, &[entity].Content, &[entity].Excerpt,
        &[entity].FeaturedImage, &[entity].Status, &[entity].ViewCount, &[entity].UserID,
        &[entity].CategoryID, &[entity].SubCategoryID, &[entity].PublishedAt,
        &[entity].CreatedAt, &[entity].UpdatedAt,
        &userEmail, &userFirstName, &userLastName, &userRole,
        &categoryName, &categorySlug,
        &subCategoryName, &subCategorySlug,
    )
    if err != nil {
        return nil, err
    }
    
    // Set related entities if they exist
    if userEmail.Valid {
        [entity].User = models.User{
            ID:        [entity].UserID,
            Email:     userEmail.String,
            FirstName: userFirstName.String,
            LastName:  userLastName.String,
            Role:      userRole.String,
        }
    }
    
    if categoryName.Valid {
        [entity].Category = models.Category{
            ID:   [entity].CategoryID,
            Name: categoryName.String,
            Slug: categorySlug.String,
        }
    }
    
    if subCategoryName.Valid {
        [entity].SubCategory = models.SubCategory{
            ID:   [entity].SubCategoryID,
            Name: subCategoryName.String,
            Slug: subCategorySlug.String,
        }
    }
    
    return [entity], nil
}
```

## 🗃️ Database Transactions

### 📄 repositories/transaction.go

```go
package repositories

import (
    "context"
    "database/sql"
    "fmt"
)

// TxManager handles database transactions
type TxManager struct {
    db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
    return &TxManager{db: db}
}

// WithTransaction executes a function within a database transaction
func (tm *TxManager) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
    tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{
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
    
    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, err)
        }
        return err
    }
    
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}

// WithReadOnlyTransaction executes a function within a read-only transaction
func (tm *TxManager) WithReadOnlyTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
    tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelReadCommitted,
        ReadOnly:  true,
    })
    if err != nil {
        return fmt.Errorf("failed to begin read-only transaction: %w", err)
    }
    
    defer tx.Rollback() // Read-only transactions can always be rolled back
    
    return fn(tx)
}
```

## 🎯 Data Layer Best Practices

### 1. **Repository Pattern**
- One repository per aggregate root
- Keep repositories focused on data access, not business logic
- Use interfaces to allow for easy testing and swapping implementations

### 2. **Query Optimization**
```go
// Use JOINs to avoid N+1 queries
func (r *[entity]Repository) GetWithRelations(ctx context.Context, id int64) (*models.[Entity], error) {
    // Single query with JOINs instead of multiple queries
    query := r.qb.Select("e.*, u.name, c.name").
        From("[entities] e").
        LeftJoin("users u ON u.id = e.user_id").
        LeftJoin("categories c ON c.id = e.category_id").
        Where(squirrel.Eq{"e.id": id})
    // ...
}
```

### 3. **Error Handling**
```go
// Distinguish between different types of database errors
func (r *[entity]Repository) Create(ctx context.Context, [entity] *models.[Entity]) (*models.[Entity], error) {
    // ...
    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&[entity].ID)
    if err != nil {
        // Check for constraint violations, duplicate keys, etc.
        if isDuplicateKeyError(err) {
            return nil, fmt.Errorf("[entity] with this slug already exists")
        }
        return nil, fmt.Errorf("failed to create [entity]: %w", err)
    }
    // ...
}
```

### 4. **Context Usage**
- Always use context for cancellation and timeouts
- Pass context to all database operations
- Handle context cancellation gracefully

### 5. **Connection Management**
```go
// Configure connection pool properly
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

---

**Previous**: [← Business Layer](./business-layer.md) | **Next**: [Configuration →](../configuration/README.md)

**Last Updated:** [YYYY-MM-DD]
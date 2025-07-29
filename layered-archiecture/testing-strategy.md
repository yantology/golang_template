# 🧪 Testing Strategy

Comprehensive testing strategy covering unit tests, integration tests, and end-to-end testing for the Golang backend application.

## 📋 Testing Pyramid

```
    /\
   /  \     E2E Tests (Few)
  /____\    - Full application tests
 /      \   - API contract tests
/________\  Integration Tests (Some)
          \ - Database integration
           \- External service tests
            \
             \____________
              Unit Tests (Many)
              - Business logic
              - Individual functions
              - Repository patterns
```

## 🎯 Testing Levels

### 1. **Unit Tests** (70% of tests)
- Test individual functions and methods
- Mock external dependencies
- Fast execution (< 1ms per test)
- High code coverage

### 2. **Integration Tests** (20% of tests)
- Test component interactions
- Real database connections
- External service integrations
- Moderate execution time

### 3. **End-to-End Tests** (10% of tests)
- Test complete user workflows
- Full application stack
- Slower execution
- Critical path coverage

## 📂 Testing Structure

```
tests/
├── unit/                   # Unit tests
│   ├── handlers/          # API handler tests
│   ├── services/          # Business logic tests
│   ├── repositories/      # Repository tests with mocks
│   └── utils/             # Utility function tests
├── integration/           # Integration tests
│   ├── api/              # API integration tests
│   ├── database/         # Database integration tests
│   └── external/         # External service tests
├── e2e/                  # End-to-end tests
│   ├── user_flows/       # User workflow tests
│   └── admin_flows/      # Admin workflow tests
├── fixtures/             # Test data fixtures
│   ├── json/            # JSON test data
│   └── sql/             # SQL test data
├── mocks/               # Generated mocks
└── helpers/             # Test helper functions
```

## 🔧 Unit Testing

### 📄 Service Unit Tests

```go
// tests/unit/services/article_service_test.go
package services_test

import (
    "context"
    "testing"
    "time"
    
    "[module-name]/internal/business/services"
    "[module-name]/internal/data/models"
    "[module-name]/pkg/errors"
    "[module-name]/tests/mocks"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

func TestArticleService_CreateArticle_Success(t *testing.T) {
    // Arrange
    ctx := context.Background()
    
    mockArticleRepo := mocks.NewMockArticleRepository(t)
    mockUserRepo := mocks.NewMockUserRepository(t)
    mockCategoryRepo := mocks.NewMockCategoryRepository(t)
    mockEmailService := mocks.NewMockEmailService(t)
    mockLogger := mocks.NewMockLogger(t)
    
    service := services.NewArticleService(
        mockArticleRepo,
        mockUserRepo,
        mockCategoryRepo,
        mockEmailService,
        mockLogger,
    )
    
    req := &services.CreateArticleRequest{
        Title:      "Test Article",
        Content:    "This is test content for the article",
        CategoryID: 1,
        UserID:     1,
    }
    
    expectedCategory := &models.Category{
        ID:   1,
        Name: "Technology",
        Slug: "technology",
    }
    
    expectedArticle := &models.Article{
        ID:         1,
        Title:      req.Title,
        Content:    req.Content,
        CategoryID: req.CategoryID,
        UserID:     req.UserID,
        Status:     "draft",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    
    // Setup mocks
    mockCategoryRepo.On("GetByID", ctx, int64(1)).Return(expectedCategory, nil)
    mockArticleRepo.On("Create", ctx, mock.AnythingOfType("*models.Article")).Return(expectedArticle, nil)
    mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
    
    // Act
    result, err := service.CreateArticle(ctx, req)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expectedArticle.Title, result.Title)
    assert.Equal(t, expectedArticle.Status, result.Status)
    assert.Equal(t, expectedArticle.UserID, result.UserID)
    
    // Verify all expectations
    mockCategoryRepo.AssertExpectations(t)
    mockArticleRepo.AssertExpectations(t)
    mockLogger.AssertExpectations(t)
}

func TestArticleService_CreateArticle_InvalidCategory(t *testing.T) {
    // Arrange
    ctx := context.Background()
    
    mockArticleRepo := mocks.NewMockArticleRepository(t)
    mockUserRepo := mocks.NewMockUserRepository(t)
    mockCategoryRepo := mocks.NewMockCategoryRepository(t)
    mockEmailService := mocks.NewMockEmailService(t)
    mockLogger := mocks.NewMockLogger(t)
    
    service := services.NewArticleService(
        mockArticleRepo,
        mockUserRepo,
        mockCategoryRepo,
        mockEmailService,
        mockLogger,
    )
    
    req := &services.CreateArticleRequest{
        Title:      "Test Article",
        Content:    "This is test content",
        CategoryID: 999, // Non-existent category
        UserID:     1,
    }
    
    // Setup mocks
    mockCategoryRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.NewNotFoundError("category not found"))
    
    // Act
    result, err := service.CreateArticle(ctx, req)
    
    // Assert
    require.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "invalid category_id")
    
    mockCategoryRepo.AssertExpectations(t)
}

func TestArticleService_PublishArticle_Success(t *testing.T) {
    // Arrange
    ctx := context.Background()
    userID := int64(1)
    articleID := int64(1)
    
    mockArticleRepo := mocks.NewMockArticleRepository(t)
    mockUserRepo := mocks.NewMockUserRepository(t)
    mockCategoryRepo := mocks.NewMockCategoryRepository(t)
    mockEmailService := mocks.NewMockEmailService(t)
    mockLogger := mocks.NewMockLogger(t)
    
    service := services.NewArticleService(
        mockArticleRepo,
        mockUserRepo,
        mockCategoryRepo,
        mockEmailService,
        mockLogger,
    )
    
    article := &models.Article{
        ID:            articleID,
        Title:         "Test Article",
        Content:       "This is a long enough content for publishing",
        FeaturedImage: "https://example.com/image.jpg",
        Status:        "draft",
        UserID:        userID,
    }
    
    // Setup mocks
    mockArticleRepo.On("GetByID", ctx, articleID).Return(article, nil)
    mockArticleRepo.On("Update", ctx, mock.AnythingOfType("*models.Article")).Return(article, nil)
    mockEmailService.On("NotifyNewPublication", mock.AnythingOfType("*models.Article")).Return(nil)
    mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
    
    // Act
    err := service.PublishArticle(ctx, articleID, userID)
    
    // Assert
    require.NoError(t, err)
    
    // Verify the article was updated with published status
    publishCall := mockArticleRepo.Calls[1] // Second call is Update
    updatedArticle := publishCall.Arguments[1].(*models.Article)
    assert.Equal(t, "published", updatedArticle.Status)
    assert.NotNil(t, updatedArticle.PublishedAt)
    
    mockArticleRepo.AssertExpectations(t)
    mockEmailService.AssertExpectations(t)
    mockLogger.AssertExpectations(t)
}
```

### 📄 Repository Unit Tests (with Mocks)

```go
// tests/unit/repositories/article_repository_test.go
package repositories_test

import (
    "context"
    "testing"
    "time"
    
    "[module-name]/internal/data/models"
    "[module-name]/internal/data/repositories"
    "[module-name]/tests/helpers"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestArticleRepository_Create_Success(t *testing.T) {
    // Arrange
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()
    
    repo := repositories.NewArticleRepository(db)
    ctx := context.Background()
    
    article := &models.Article{
        Title:      "Test Article",
        Slug:       "test-article",
        Content:    "Test content",
        Status:     "draft",
        UserID:     1,
        CategoryID: 1,
    }
    
    now := time.Now()
    
    // Setup mock expectations
    mock.ExpectQuery(`INSERT INTO articles`).
        WithArgs(article.Title, article.Slug, article.Content, article.Excerpt, 
                article.FeaturedImage, article.Status, article.UserID, 
                article.CategoryID, article.SubCategoryID).
        WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
            AddRow(1, now, now))
    
    // Act
    result, err := repo.Create(ctx, article)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, int64(1), result.ID)
    assert.Equal(t, article.Title, result.Title)
    assert.WithinDuration(t, now, result.CreatedAt, time.Second)
    
    // Verify all expectations were met
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestArticleRepository_GetByID_Success(t *testing.T) {
    // Arrange
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()
    
    repo := repositories.NewArticleRepository(db)
    ctx := context.Background()
    articleID := int64(1)
    
    now := time.Now()
    
    // Setup mock expectations
    rows := sqlmock.NewRows([]string{
        "e.id", "e.title", "e.slug", "e.content", "e.excerpt", "e.featured_image",
        "e.status", "e.view_count", "e.user_id", "e.category_id", "e.sub_category_id",
        "e.published_at", "e.created_at", "e.updated_at",
        "u.id", "u.email", "u.first_name", "u.last_name", "u.role",
        "c.id", "c.name", "c.slug", "c.description",
        "sc.id", "sc.name", "sc.slug",
    }).AddRow(
        1, "Test Article", "test-article", "Test content", "Test excerpt", "image.jpg",
        "published", 100, 1, 1, 2,
        now, now, now,
        1, "user@example.com", "John", "Doe", "user",
        1, "Technology", "technology", "Tech category",
        2, "Programming", "programming",
    )
    
    mock.ExpectQuery(`SELECT .+ FROM articles e`).
        WithArgs(articleID).
        WillReturnRows(rows)
    
    // Act
    result, err := repo.GetByID(ctx, articleID)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, articleID, result.ID)
    assert.Equal(t, "Test Article", result.Title)
    assert.Equal(t, "user@example.com", result.User.Email)
    assert.Equal(t, "Technology", result.Category.Name)
    assert.Equal(t, "Programming", result.SubCategory.Name)
    
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### 📄 Handler Unit Tests

```go
// tests/unit/handlers/article_handler_test.go
package handlers_test

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    
    "[module-name]/internal/api/handlers"
    "[module-name]/internal/data/models"
    "[module-name]/tests/mocks"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

func TestArticleHandler_CreateArticle_Success(t *testing.T) {
    // Arrange
    gin.SetMode(gin.TestMode)
    
    mockService := mocks.NewMockArticleService(t)
    mockLogger := mocks.NewMockLogger(t)
    
    handler := handlers.NewArticleHandler(mockService, mockLogger)
    
    requestBody := handlers.CreateArticleRequest{
        Title:      "Test Article",
        Content:    "Test content",
        CategoryID: 1,
    }
    
    expectedArticle := &models.Article{
        ID:         1,
        Title:      requestBody.Title,
        Content:    requestBody.Content,
        CategoryID: requestBody.CategoryID,
        UserID:     1,
        Status:     "draft",
        CreatedAt:  time.Now(),
    }
    
    // Setup mocks
    mockService.On("CreateArticle", mock.Anything, mock.AnythingOfType("*services.CreateArticleRequest")).
        Return(expectedArticle, nil)
    
    // Create request
    jsonBody, err := json.Marshal(requestBody)
    require.NoError(t, err)
    
    req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    
    // Create Gin context with user information
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    c.Set("user_id", int64(1))
    
    // Act
    handler.CreateArticle(c)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    err = json.Unmarshal(w.Body.Bytes(), &response)
    require.NoError(t, err)
    
    assert.Equal(t, "Article created successfully", response["message"])
    assert.NotNil(t, response["data"])
    
    mockService.AssertExpectations(t)
}

func TestArticleHandler_CreateArticle_ValidationError(t *testing.T) {
    // Arrange
    gin.SetMode(gin.TestMode)
    
    mockService := mocks.NewMockArticleService(t)
    mockLogger := mocks.NewMockLogger(t)
    
    handler := handlers.NewArticleHandler(mockService, mockLogger)
    
    // Invalid request - missing required fields
    requestBody := handlers.CreateArticleRequest{
        Title: "", // Empty title should fail validation
    }
    
    jsonBody, err := json.Marshal(requestBody)
    require.NoError(t, err)
    
    req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    c.Set("user_id", int64(1))
    
    // Act
    handler.CreateArticle(c)
    
    // Assert
    assert.Equal(t, http.StatusBadRequest, w.Code)
    
    var response map[string]interface{}
    err = json.Unmarshal(w.Body.Bytes(), &response)
    require.NoError(t, err)
    
    assert.Equal(t, "Invalid request format", response["message"])
    assert.Contains(t, response["error"], "required")
    
    // Service should not be called
    mockService.AssertNotCalled(t, "CreateArticle")
}
```

## 🔗 Integration Testing

### 📄 API Integration Tests

```go
// tests/integration/api/article_api_test.go
package api_test

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "testing"
    "time"
    
    "[module-name]/tests/helpers"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

type ArticleAPITestSuite struct {
    suite.Suite
    testServer *helpers.TestServer
    authToken  string
}

func (suite *ArticleAPITestSuite) SetupSuite() {
    // Initialize test server with real database
    suite.testServer = helpers.NewTestServer(suite.T())
    
    // Create test user and get auth token
    user := suite.testServer.CreateTestUser("test@example.com", "password123")
    suite.authToken = suite.testServer.LoginUser(user.Email, "password123")
}

func (suite *ArticleAPITestSuite) TearDownSuite() {
    suite.testServer.Close()
}

func (suite *ArticleAPITestSuite) SetupTest() {
    // Clean database before each test
    suite.testServer.CleanDatabase()
}

func (suite *ArticleAPITestSuite) TestCreateArticle_Success() {
    // Arrange
    category := suite.testServer.CreateTestCategory("Technology")
    
    requestBody := map[string]interface{}{
        "title":       "Integration Test Article",
        "content":     "This is content for integration testing",
        "category_id": category.ID,
    }
    
    jsonBody, err := json.Marshal(requestBody)
    require.NoError(suite.T(), err)
    
    // Act
    resp, err := suite.testServer.PostWithAuth("/api/v1/articles", bytes.NewBuffer(jsonBody), suite.authToken)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    // Assert
    assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
    
    var response map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&response)
    require.NoError(suite.T(), err)
    
    assert.Equal(suite.T(), "Article created successfully", response["message"])
    
    data := response["data"].(map[string]interface{})
    assert.Equal(suite.T(), requestBody["title"], data["title"])
    assert.Equal(suite.T(), "draft", data["status"])
    assert.NotNil(suite.T(), data["id"])
}

func (suite *ArticleAPITestSuite) TestGetArticle_Success() {
    // Arrange
    article := suite.testServer.CreateTestArticle("Test Article", "published")
    
    // Act
    resp, err := suite.testServer.Get(fmt.Sprintf("/api/v1/articles/%d", article.ID))
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    // Assert
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
    
    var response map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&response)
    require.NoError(suite.T(), err)
    
    data := response["data"].(map[string]interface{})
    assert.Equal(suite.T(), article.Title, data["title"])
    assert.Equal(suite.T(), "published", data["status"])
}

func (suite *ArticleAPITestSuite) TestListArticles_WithPagination() {
    // Arrange
    // Create multiple test articles
    for i := 1; i <= 15; i++ {
        suite.testServer.CreateTestArticle(fmt.Sprintf("Article %d", i), "published")
    }
    
    // Act
    resp, err := suite.testServer.Get("/api/v1/articles?page=1&limit=10")
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    // Assert
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
    
    var response map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&response)
    require.NoError(suite.T(), err)
    
    data := response["data"].([]interface{})
    assert.Len(suite.T(), data, 10) // Should return 10 items
    
    pagination := response["pagination"].(map[string]interface{})
    assert.Equal(suite.T(), float64(1), pagination["page"])
    assert.Equal(suite.T(), float64(10), pagination["limit"])
    assert.Equal(suite.T(), float64(15), pagination["total"])
    assert.Equal(suite.T(), float64(2), pagination["total_pages"])
}

func (suite *ArticleAPITestSuite) TestUpdateArticle_Unauthorized() {
    // Arrange
    otherUser := suite.testServer.CreateTestUser("other@example.com", "password123")
    article := suite.testServer.CreateTestArticleForUser("Test Article", "draft", otherUser.ID)
    
    requestBody := map[string]interface{}{
        "title": "Updated Title",
    }
    jsonBody, err := json.Marshal(requestBody)
    require.NoError(suite.T(), err)
    
    // Act - Try to update article owned by another user
    resp, err := suite.testServer.PutWithAuth(
        fmt.Sprintf("/api/v1/articles/%d", article.ID), 
        bytes.NewBuffer(jsonBody), 
        suite.authToken,
    )
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    // Assert
    assert.Equal(suite.T(), http.StatusForbidden, resp.StatusCode)
}

func TestArticleAPITestSuite(t *testing.T) {
    suite.Run(t, new(ArticleAPITestSuite))
}
```

### 📄 Database Integration Tests

```go
// tests/integration/database/article_repository_test.go
package database_test

import (
    "context"
    "testing"
    "time"
    
    "[module-name]/internal/data/models"
    "[module-name]/internal/data/repositories"
    "[module-name]/tests/helpers"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

type ArticleRepositoryTestSuite struct {
    suite.Suite
    testDB *helpers.TestDatabase
    repo   repositories.ArticleRepository
}

func (suite *ArticleRepositoryTestSuite) SetupSuite() {
    suite.testDB = helpers.NewTestDatabase(suite.T())
    suite.repo = repositories.NewArticleRepository(suite.testDB.DB)
}

func (suite *ArticleRepositoryTestSuite) TearDownSuite() {
    suite.testDB.Close()
}

func (suite *ArticleRepositoryTestSuite) SetupTest() {
    suite.testDB.CleanTables("articles", "users", "categories")
}

func (suite *ArticleRepositoryTestSuite) TestCreate_Success() {
    // Arrange
    ctx := context.Background()
    
    // Create test dependencies
    user := suite.testDB.CreateTestUser("test@example.com")
    category := suite.testDB.CreateTestCategory("Technology")
    
    article := &models.Article{
        Title:      "Test Article",
        Slug:       "test-article",
        Content:    "This is test content",
        Status:     "draft",
        UserID:     user.ID,
        CategoryID: category.ID,
    }
    
    // Act
    result, err := suite.repo.Create(ctx, article)
    
    // Assert
    require.NoError(suite.T(), err)
    assert.NotNil(suite.T(), result)
    assert.NotZero(suite.T(), result.ID)
    assert.Equal(suite.T(), article.Title, result.Title)
    assert.Equal(suite.T(), article.UserID, result.UserID)
    assert.WithinDuration(suite.T(), time.Now(), result.CreatedAt, time.Second)
    assert.WithinDuration(suite.T(), time.Now(), result.UpdatedAt, time.Second)
}

func (suite *ArticleRepositoryTestSuite) TestGetByID_WithRelations() {
    // Arrange
    ctx := context.Background()
    
    user := suite.testDB.CreateTestUser("test@example.com")
    category := suite.testDB.CreateTestCategory("Technology")
    article := suite.testDB.CreateTestArticle("Test Article", user.ID, category.ID)
    
    // Act
    result, err := suite.repo.GetByID(ctx, article.ID)
    
    // Assert
    require.NoError(suite.T(), err)
    assert.NotNil(suite.T(), result)
    assert.Equal(suite.T(), article.ID, result.ID)
    assert.Equal(suite.T(), article.Title, result.Title)
    
    // Check relations are loaded
    assert.Equal(suite.T(), user.Email, result.User.Email)
    assert.Equal(suite.T(), category.Name, result.Category.Name)
}

func (suite *ArticleRepositoryTestSuite) TestList_WithFiltering() {
    // Arrange
    ctx := context.Background()
    
    user := suite.testDB.CreateTestUser("test@example.com")
    techCategory := suite.testDB.CreateTestCategory("Technology")
    sportsCategory := suite.testDB.CreateTestCategory("Sports")
    
    // Create test articles
    suite.testDB.CreateTestArticleWithCategory("Tech Article 1", user.ID, techCategory.ID, "published")
    suite.testDB.CreateTestArticleWithCategory("Tech Article 2", user.ID, techCategory.ID, "draft")
    suite.testDB.CreateTestArticleWithCategory("Sports Article", user.ID, sportsCategory.ID, "published")
    
    params := &repositories.ListArticlesParams{
        Page:       1,
        Limit:      10,
        CategoryID: techCategory.ID,
        Status:     "published",
    }
    
    // Act
    articles, total, err := suite.repo.List(ctx, params)
    
    // Assert
    require.NoError(suite.T(), err)
    assert.Len(suite.T(), articles, 1) // Only one published tech article
    assert.Equal(suite.T(), int64(1), total)
    assert.Equal(suite.T(), "Tech Article 1", articles[0].Title)
}

func (suite *ArticleRepositoryTestSuite) TestUpdate_Success() {
    // Arrange
    ctx := context.Background()
    
    user := suite.testDB.CreateTestUser("test@example.com")
    category := suite.testDB.CreateTestCategory("Technology")
    article := suite.testDB.CreateTestArticle("Original Title", user.ID, category.ID)
    
    // Modify article
    article.Title = "Updated Title"
    article.Content = "Updated content"
    article.Status = "published"
    
    // Act
    result, err := suite.repo.Update(ctx, article)
    
    // Assert
    require.NoError(suite.T(), err)
    assert.NotNil(suite.T(), result)
    assert.Equal(suite.T(), "Updated Title", result.Title)
    assert.Equal(suite.T(), "published", result.Status)
    
    // Verify in database
    retrieved, err := suite.repo.GetByID(ctx, article.ID)
    require.NoError(suite.T(), err)
    assert.Equal(suite.T(), "Updated Title", retrieved.Title)
    assert.Equal(suite.T(), "published", retrieved.Status)
}

func TestArticleRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(ArticleRepositoryTestSuite))
}
```

## 🌐 End-to-End Testing

### 📄 User Flow E2E Tests

```go
// tests/e2e/user_flows/article_workflow_test.go
package e2e_test

import (
    "testing"
    "time"
    
    "[module-name]/tests/helpers"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

type ArticleWorkflowTestSuite struct {
    suite.Suite
    testServer *helpers.TestServer
}

func (suite *ArticleWorkflowTestSuite) SetupSuite() {
    suite.testServer = helpers.NewTestServer(suite.T())
}

func (suite *ArticleWorkflowTestSuite) TearDownSuite() {
    suite.testServer.Close()
}

func (suite *ArticleWorkflowTestSuite) SetupTest() {
    suite.testServer.CleanDatabase()
}

func (suite *ArticleWorkflowTestSuite) TestCompleteArticleLifecycle() {
    t := suite.T()
    
    // Step 1: User registration
    userResp := suite.testServer.RegisterUser(map[string]interface{}{
        "email":      "author@example.com",
        "password":   "password123",
        "first_name": "John",
        "last_name":  "Doe",
    })
    require.Equal(t, 201, userResp.StatusCode)
    
    // Step 2: User login
    authToken := suite.testServer.LoginUser("author@example.com", "password123")
    require.NotEmpty(t, authToken)
    
    // Step 3: Create category (admin action)
    adminToken := suite.testServer.GetAdminToken()
    categoryResp := suite.testServer.CreateCategoryWithAuth(map[string]interface{}{
        "name": "Technology",
        "slug": "technology",
    }, adminToken)
    require.Equal(t, 201, categoryResp.StatusCode)
    
    categoryData := suite.testServer.ParseResponse(categoryResp)
    categoryID := int64(categoryData["data"].(map[string]interface{})["id"].(float64))
    
    // Step 4: Create draft article
    articleResp := suite.testServer.CreateArticleWithAuth(map[string]interface{}{
        "title":       "My First Article",
        "content":     "This is the content of my first article. It's long enough for publishing.",
        "category_id": categoryID,
    }, authToken)
    require.Equal(t, 201, articleResp.StatusCode)
    
    articleData := suite.testServer.ParseResponse(articleResp)
    articleID := int64(articleData["data"].(map[string]interface{})["id"].(float64))
    
    // Verify article is created as draft
    article := articleData["data"].(map[string]interface{})
    assert.Equal(t, "My First Article", article["title"])
    assert.Equal(t, "draft", article["status"])
    
    // Step 5: Update article with featured image
    updateResp := suite.testServer.UpdateArticleWithAuth(articleID, map[string]interface{}{
        "featured_image": "https://example.com/featured-image.jpg",
    }, authToken)
    require.Equal(t, 200, updateResp.StatusCode)
    
    // Step 6: Publish article
    publishResp := suite.testServer.PublishArticleWithAuth(articleID, authToken)
    require.Equal(t, 200, publishResp.StatusCode)
    
    // Step 7: Verify article is published and visible publicly
    getResp := suite.testServer.GetArticle(articleID)
    require.Equal(t, 200, getResp.StatusCode)
    
    publishedData := suite.testServer.ParseResponse(getResp)
    publishedArticle := publishedData["data"].(map[string]interface{})
    assert.Equal(t, "published", publishedArticle["status"])
    assert.NotNil(t, publishedArticle["published_at"])
    
    // Step 8: Verify article appears in public listing
    listResp := suite.testServer.ListArticles(map[string]string{
        "status": "published",
    })
    require.Equal(t, 200, listResp.StatusCode)
    
    listData := suite.testServer.ParseResponse(listResp)
    articles := listData["data"].([]interface{})
    assert.Len(t, articles, 1)
    assert.Equal(t, "My First Article", articles[0].(map[string]interface{})["title"])
    
    // Step 9: Simulate view count increment
    // Multiple requests to the article should increment view count
    for i := 0; i < 5; i++ {
        suite.testServer.GetArticle(articleID)
        time.Sleep(100 * time.Millisecond) // Small delay between requests
    }
    
    // Verify view count increased
    finalResp := suite.testServer.GetArticle(articleID)
    finalData := suite.testServer.ParseResponse(finalResp)
    finalArticle := finalData["data"].(map[string]interface{})
    viewCount := int(finalArticle["view_count"].(float64))
    assert.GreaterOrEqual(t, viewCount, 5)
    
    // Step 10: Archive article
    archiveResp := suite.testServer.ArchiveArticleWithAuth(articleID, authToken)
    require.Equal(t, 200, archiveResp.StatusCode)
    
    // Step 11: Verify archived article is not in public listing
    archivedListResp := suite.testServer.ListArticles(map[string]string{
        "status": "published",
    })
    require.Equal(t, 200, archivedListResp.StatusCode)
    
    archivedListData := suite.testServer.ParseResponse(archivedListResp)
    archivedArticles := archivedListData["data"].([]interface{})
    assert.Len(t, archivedArticles, 0) // No published articles
}

func (suite *ArticleWorkflowTestSuite) TestUnauthorizedAccess() {
    t := suite.T()
    
    // Create a user and article
    userToken := suite.testServer.CreateUserAndGetToken("user@example.com")
    articleID := suite.testServer.CreateTestArticleForToken("Test Article", userToken)
    
    // Create another user
    otherUserToken := suite.testServer.CreateUserAndGetToken("other@example.com")
    
    // Try to update article with different user
    updateResp := suite.testServer.UpdateArticleWithAuth(articleID, map[string]interface{}{
        "title": "Hacked Title",
    }, otherUserToken)
    assert.Equal(t, 403, updateResp.StatusCode)
    
    // Try to delete article with different user
    deleteResp := suite.testServer.DeleteArticleWithAuth(articleID, otherUserToken)
    assert.Equal(t, 403, deleteResp.StatusCode)
    
    // Verify original article is unchanged
    getResp := suite.testServer.GetArticle(articleID)
    require.Equal(t, 200, getResp.StatusCode)
    
    data := suite.testServer.ParseResponse(getResp)
    article := data["data"].(map[string]interface{})
    assert.Equal(t, "Test Article", article["title"])
}

func TestArticleWorkflowTestSuite(t *testing.T) {
    suite.Run(t, new(ArticleWorkflowTestSuite))
}
```

## 🛠️ Test Helpers

### 📄 Test Server Helper

```go
// tests/helpers/test_server.go
package helpers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "[module-name]/cmd/api"
    "[module-name]/internal/config"
    
    "github.com/stretchr/testify/require"
)

type TestServer struct {
    server *httptest.Server
    app    *main.Application
    db     *TestDatabase
    t      *testing.T
}

func NewTestServer(t *testing.T) *TestServer {
    // Load test configuration
    cfg := &config.Config{
        Server: config.ServerConfig{
            Port: "0", // Random port for testing
            Host: "localhost",
            Env:  "test",
        },
        Database: config.DatabaseConfig{
            Type:     "postgres",
            Host:     "localhost",
            Port:     "5433", // Different port for test database
            User:     "postgres",
            Password: "postgres",
            Name:     "golang_template_test",
            SSLMode:  "disable",
        },
        // ... other test config
    }
    
    // Create test database
    testDB := NewTestDatabase(t)
    
    // Create application with test config
    app, err := main.NewApplicationWithDB(cfg, testDB.DB)
    require.NoError(t, err)
    
    // Create test server
    server := httptest.NewServer(app.Server.Handler)
    
    return &TestServer{
        server: server,
        app:    app,
        db:     testDB,
        t:      t,
    }
}

func (ts *TestServer) Close() {
    ts.server.Close()
    ts.db.Close()
}

func (ts *TestServer) CleanDatabase() {
    ts.db.CleanAllTables()
}

func (ts *TestServer) RegisterUser(userData map[string]interface{}) *http.Response {
    jsonBody, err := json.Marshal(userData)
    require.NoError(ts.t, err)
    
    resp, err := http.Post(
        ts.server.URL+"/api/v1/auth/register",
        "application/json",
        bytes.NewBuffer(jsonBody),
    )
    require.NoError(ts.t, err)
    
    return resp
}

func (ts *TestServer) LoginUser(email, password string) string {
    loginData := map[string]interface{}{
        "email":    email,
        "password": password,
    }
    
    jsonBody, err := json.Marshal(loginData)
    require.NoError(ts.t, err)
    
    resp, err := http.Post(
        ts.server.URL+"/api/v1/auth/login",
        "application/json",
        bytes.NewBuffer(jsonBody),
    )
    require.NoError(ts.t, err)
    defer resp.Body.Close()
    
    var response map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&response)
    require.NoError(ts.t, err)
    
    data := response["data"].(map[string]interface{})
    return data["access_token"].(string)
}

func (ts *TestServer) PostWithAuth(path string, body *bytes.Buffer, token string) (*http.Response, error) {
    req, err := http.NewRequest(http.MethodPost, ts.server.URL+path, body)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)
    
    return http.DefaultClient.Do(req)
}

func (ts *TestServer) Get(path string) (*http.Response, error) {
    return http.Get(ts.server.URL + path)
}

func (ts *TestServer) ParseResponse(resp *http.Response) map[string]interface{} {
    defer resp.Body.Close()
    
    var data map[string]interface{}
    err := json.NewDecoder(resp.Body).Decode(&data)
    require.NoError(ts.t, err)
    
    return data
}

// Additional helper methods for creating test data...
```

## 🎯 Testing Best Practices

### 1. **Test Organization**
- Follow the AAA pattern (Arrange, Act, Assert)
- Use descriptive test names that explain the scenario
- Group related tests using test suites

### 2. **Test Data Management**
```go
// Use test fixtures for consistent test data
func loadTestData(t *testing.T, filename string) map[string]interface{} {
    data, err := os.ReadFile(fmt.Sprintf("fixtures/%s.json", filename))
    require.NoError(t, err)
    
    var result map[string]interface{}
    err = json.Unmarshal(data, &result)
    require.NoError(t, err)
    
    return result
}
```

### 3. **Mock Management**
```go
// Use interfaces for easy mocking
//go:generate mockery --name=ArticleService --output=../mocks
type ArticleService interface {
    CreateArticle(ctx context.Context, req *CreateArticleRequest) (*models.Article, error)
    // ... other methods
}
```

### 4. **Test Isolation**
- Each test should be independent
- Clean database state between tests
- Use separate test databases

### 5. **Performance Testing**
```go
func BenchmarkArticleList(b *testing.B) {
    // Setup
    service := setupTestService()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _, err := service.ListArticles(context.Background(), &ListParams{
            Page:  1,
            Limit: 20,
        })
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

---

**Previous**: [← Application Entry](./06-application-entry.md) | **Next**: [Environment Setup →](./08-environment-setup.md)

**Last Updated:** [YYYY-MM-DD]
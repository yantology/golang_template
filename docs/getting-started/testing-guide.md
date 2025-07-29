# Testing Guide

Panduan praktis untuk testing aplikasi Go. Cocok untuk pemula yang mau belajar testing dengan benar.

## 🧪 Apa itu Testing?

**Testing** adalah cara memastikan code kita bekerja dengan benar.

**Analogi sederhana:**
- Code = Resep masakan
- Test = Mencicipi masakan untuk pastikan rasanya enak
- Automated test = Robot yang bisa mencicipi otomatis

**Kenapa testing penting?**
- ✅ Pastikan fitur bekerja dengan benar
- ✅ Cegah bug masuk ke production
- ✅ Confidence saat ubah code
- ✅ Dokumentasi bagaimana code seharusnya bekerja

## 🏗️ Jenis-Jenis Test

### 1. Unit Test (Test Kecil)
- Test function/method individual
- Cepat dijalankan (milidetik)
- Tidak akses database/network
- Pakai mock untuk dependency

**Kapan pakai:** Test business logic, helper functions

### 2. Integration Test (Test Medium)  
- Test beberapa component berinteraksi
- Akses database/external service
- Lebih lambat dari unit test
- Setup/cleanup diperlukan

**Kapan pakai:** Test API endpoints, database operations

### 3. End-to-End Test (Test Besar)
- Test full user journey
- Pakai browser/HTTP client
- Paling lambat
- Most realistic

**Kapan pakai:** Test critical user flows

## 🚀 Running Tests

### Basic Commands
```bash
# Run semua tests
make test

# Run dengan coverage report
make test-coverage

# Run specific test types
make test-unit        # Unit tests only
make test-integration # Integration tests only
make test-e2e        # End-to-end tests only
```

### Advanced Commands
```bash
# Run tests dengan verbose output
go test -v ./...

# Run specific test function
go test -v ./internal/business/services -run TestUserService_CreateUser

# Run tests dengan race condition detection
go test -race ./...

# Run tests dengan timeout
go test -timeout 30s ./...
```

## 📁 Test Organization

```
tests/
├── unit/           # Unit tests (fast, isolated)
├── integration/    # Integration tests (database, APIs)
├── e2e/           # End-to-end tests (full scenarios)
├── fixtures/      # Test data (JSON, SQL, etc.)
└── helpers/       # Test helper functions
```

**Convention:**
- Unit tests: `*_test.go` dalam package yang sama
- Integration tests: di folder `tests/integration/`
- E2E tests: di folder `tests/e2e/`

## ✍️ Writing Unit Tests

### 1. Simple Function Test

**File: `internal/pkg/utils/string_test.go`**
```go
package utils

import (
    \"testing\"
    \"github.com/stretchr/testify/assert\"
)

func TestCapitalize(t *testing.T) {
    // Test cases
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {\"empty string\", \"\", \"\"},
        {\"single char\", \"a\", \"A\"},
        {\"normal word\", \"hello\", \"Hello\"},
        {\"already capitalized\", \"Hello\", \"Hello\"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Capitalize(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 2. Service Layer Test (dengan Mock)

**File: `internal/business/services/user_service_test.go`**
```go
package services

import (
    \"context\"
    \"testing\"
    \"github.com/stretchr/testify/assert\"
    \"github.com/stretchr/testify/mock\"
)

// Mock repository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
    args := m.Called(ctx, user)
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    args := m.Called(ctx, email)
    return args.Get(0).(*models.User), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
    // Setup
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo, nil)
    
    user := &models.User{
        Email:     \"test@example.com\",
        FirstName: \"John\",
        LastName:  \"Doe\",
    }
    
    // Mock expectations
    mockRepo.On(\"GetByEmail\", mock.Anything, user.Email).Return(nil, errors.New(\"not found\"))
    mockRepo.On(\"Create\", mock.Anything, mock.AnythingOfType(\"*models.User\")).Return(user, nil)
    
    // Execute
    result, err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, user.Email, result.Email)
    assert.NotEmpty(t, result.ID)
    
    // Verify mock was called
    mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_EmailExists(t *testing.T) {
    // Setup
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo, nil)
    
    existingUser := &models.User{Email: \"test@example.com\"}
    newUser := &models.User{Email: \"test@example.com\"}
    
    // Mock expectations - email already exists
    mockRepo.On(\"GetByEmail\", mock.Anything, newUser.Email).Return(existingUser, nil)
    
    // Execute
    result, err := service.CreateUser(context.Background(), newUser)
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), \"email already exists\")
    
    mockRepo.AssertExpectations(t)
}
```

## 🔗 Integration Tests

### 1. Database Test

**File: `tests/integration/user_repository_test.go`**
```go
package integration

import (
    \"context\"
    \"testing\"
    \"github.com/stretchr/testify/assert\"
    \"github.com/stretchr/testify/require\"
)

func TestUserRepository_Create(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := repositories.NewUserRepository(db)
    
    user := &models.User{
        Email:     \"test@example.com\",
        FirstName: \"John\",
        LastName:  \"Doe\",
    }
    
    // Execute
    result, err := repo.Create(context.Background(), user)
    
    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, result.ID)
    assert.Equal(t, user.Email, result.Email)
    assert.NotZero(t, result.CreatedAt)
    
    // Verify in database
    var count int
    err = db.QueryRow(\"SELECT COUNT(*) FROM users WHERE email = $1\", user.Email).Scan(&count)
    require.NoError(t, err)
    assert.Equal(t, 1, count)
}

// Helper functions
func setupTestDB(t *testing.T) *sql.DB {
    // Connect to test database
    db, err := sql.Open(\"postgres\", \"postgres://postgres:password@localhost/test_db\")
    require.NoError(t, err)
    
    // Run migrations
    migrator := migrate.New(\"file://../../internal/data/migrations\", \"postgres://...\")
    err = migrator.Up()
    require.NoError(t, err)
    
    return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
    // Clean up test data
    _, err := db.Exec(\"TRUNCATE TABLE users CASCADE\")
    require.NoError(t, err)
    
    db.Close()
}
```

### 2. API Integration Test

**File: `tests/integration/user_api_test.go`**
```go
package integration

import (
    \"bytes\"
    \"encoding/json\"
    \"net/http\"
    \"net/http/httptest\"
    \"testing\"
    \"github.com/stretchr/testify/assert\"
    \"github.com/stretchr/testify/require\"
)

func TestUserAPI_CreateUser(t *testing.T) {
    // Setup test app
    app := setupTestApp(t)
    defer cleanupTestApp(t, app)
    
    // Test data
    payload := map[string]interface{}{
        \"email\":      \"test@example.com\",
        \"first_name\": \"John\",
        \"last_name\":  \"Doe\",
        \"password\":   \"password123\",
    }
    
    jsonPayload, _ := json.Marshal(payload)
    
    // Make request
    req := httptest.NewRequest(\"POST\", \"/api/v1/users\", bytes.NewReader(jsonPayload))
    req.Header.Set(\"Content-Type\", \"application/json\")
    
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    // Assert response
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    require.NoError(t, err)
    
    assert.True(t, response[\"success\"].(bool))
    assert.Equal(t, \"User created successfully\", response[\"message\"])
    
    // Check user data
    userData := response[\"data\"].(map[string]interface{})
    assert.Equal(t, payload[\"email\"], userData[\"email\"])
    assert.NotEmpty(t, userData[\"id\"])
}

func TestUserAPI_CreateUser_ValidationError(t *testing.T) {
    app := setupTestApp(t)
    defer cleanupTestApp(t, app)
    
    // Invalid payload (missing email)
    payload := map[string]interface{}{
        \"first_name\": \"John\",
        \"last_name\":  \"Doe\",
    }
    
    jsonPayload, _ := json.Marshal(payload)
    
    req := httptest.NewRequest(\"POST\", \"/api/v1/users\", bytes.NewReader(jsonPayload))
    req.Header.Set(\"Content-Type\", \"application/json\")
    
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    // Should return validation error
    assert.Equal(t, http.StatusBadRequest, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    
    assert.False(t, response[\"success\"].(bool))
    assert.Contains(t, response[\"message\"].(string), \"validation\")
}
```

## 🌐 End-to-End Tests

**File: `tests/e2e/user_flow_test.go`**
```go
package e2e

import (
    \"bytes\"
    \"encoding/json\"
    \"fmt\"
    \"net/http\"
    \"testing\"
    \"github.com/stretchr/testify/assert\"
    \"github.com/stretchr/testify/require\"
)

func TestUserRegistrationFlow(t *testing.T) {
    baseURL := \"http://localhost:8080\" // Assuming app is running
    
    // Step 1: Register user
    registerPayload := map[string]interface{}{
        \"email\":      \"e2e@example.com\",
        \"first_name\": \"E2E\",
        \"last_name\":  \"Test\",
        \"password\":   \"password123\",
    }
    
    registerResp := makeHTTPRequest(t, \"POST\", baseURL+\"/api/v1/auth/register\", registerPayload)
    assert.Equal(t, http.StatusCreated, registerResp.StatusCode)
    
    // Step 2: Login
    loginPayload := map[string]interface{}{
        \"email\":    \"e2e@example.com\",
        \"password\": \"password123\",
    }
    
    loginResp := makeHTTPRequest(t, \"POST\", baseURL+\"/api/v1/auth/login\", loginPayload)
    assert.Equal(t, http.StatusOK, loginResp.StatusCode)
    
    // Extract token
    var loginData map[string]interface{}
    json.NewDecoder(loginResp.Body).Decode(&loginData)
    token := loginData[\"data\"].(map[string]interface{})[\"access_token\"].(string)
    
    // Step 3: Access protected endpoint
    req, _ := http.NewRequest(\"GET\", baseURL+\"/api/v1/auth/me\", nil)
    req.Header.Set(\"Authorization\", \"Bearer \"+token)
    
    client := &http.Client{}
    meResp, err := client.Do(req)
    require.NoError(t, err)
    defer meResp.Body.Close()
    
    assert.Equal(t, http.StatusOK, meResp.StatusCode)
    
    var meData map[string]interface{}
    json.NewDecoder(meResp.Body).Decode(&meData)
    userData := meData[\"data\"].(map[string]interface{})
    assert.Equal(t, \"e2e@example.com\", userData[\"email\"])
}

func makeHTTPRequest(t *testing.T, method, url string, payload map[string]interface{}) *http.Response {
    jsonPayload, _ := json.Marshal(payload)
    
    req, err := http.NewRequest(method, url, bytes.NewReader(jsonPayload))
    require.NoError(t, err)
    req.Header.Set(\"Content-Type\", \"application/json\")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    require.NoError(t, err)
    
    return resp
}
```

## 🔧 Test Helpers & Fixtures

### Test Fixtures

**File: `tests/fixtures/users.json`**
```json
{
  \"valid_user\": {
    \"email\": \"john@example.com\",
    \"first_name\": \"John\",
    \"last_name\": \"Doe\",
    \"password\": \"password123\"
  },
  \"admin_user\": {
    \"email\": \"admin@example.com\", 
    \"first_name\": \"Admin\",
    \"last_name\": \"User\",
    \"role\": \"admin\"
  }
}
```

### Test Helpers

**File: `tests/helpers/database.go`**
```go
package helpers

import (
    \"database/sql\"
    \"testing\"
    \"github.com/stretchr/testify/require\"
)

func SetupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open(\"postgres\", GetTestDatabaseURL())
    require.NoError(t, err)
    
    // Run migrations
    RunMigrations(t, db)
    
    return db
}

func CleanupTestDB(t *testing.T, db *sql.DB) {
    // Clean all tables
    tables := []string{\"users\", \"products\", \"orders\"}
    for _, table := range tables {
        _, err := db.Exec(fmt.Sprintf(\"TRUNCATE TABLE %s CASCADE\", table))
        require.NoError(t, err)
    }
    
    db.Close()
}

func GetTestDatabaseURL() string {
    return \"postgres://postgres:password@localhost/test_db?sslmode=disable\"
}
```

## 📊 Test Coverage

### Check Coverage
```bash
# Generate coverage report
make test-coverage

# View HTML report
go tool cover -html=coverage.out
```

### Coverage Goals
- **Unit tests:** 80%+ coverage
- **Critical paths:** 95%+ coverage  
- **Edge cases:** Should be covered

### Coverage Tips
- Focus on business logic
- Don't obsess over 100% coverage
- Test error cases too
- Mock external dependencies

## 🎯 Testing Best Practices

### ✅ DO (Lakukan)
- **Write tests for business logic**
- **Test error cases and edge cases**
- **Use descriptive test names**
- **Keep tests independent**
- **Use table-driven tests for multiple cases**
- **Clean up test data**
- **Mock external dependencies**

### ❌ DON'T (Jangan)
- **Don't test framework code**
- **Don't make tests depend on each other**
- **Don't use real external services in tests**
- **Don't ignore flaky tests**
- **Don't over-mock (mock only what you need)**

### Test Naming Convention
```go
// Pattern: TestMethodName_Scenario_ExpectedBehavior
func TestUserService_CreateUser_WithValidData_ReturnsUser(t *testing.T) {}
func TestUserService_CreateUser_WithDuplicateEmail_ReturnsError(t *testing.T) {}
func TestUserService_CreateUser_WithInvalidEmail_ReturnsValidationError(t *testing.T) {}
```

## 🚨 Common Testing Mistakes

### 1. Not Cleaning Up Test Data
```go
// ❌ Bad: Test data affects other tests
func TestCreateUser(t *testing.T) {
    user := createTestUser(\"test@example.com\")
    // ... test logic ...
    // No cleanup!
}

// ✅ Good: Always cleanup
func TestCreateUser(t *testing.T) {
    user := createTestUser(\"test@example.com\")
    defer deleteTestUser(user.ID) // Cleanup
    // ... test logic ...
}
```

### 2. Testing Implementation Instead of Behavior
```go
// ❌ Bad: Testing internal implementation
func TestUserService_CreateUser(t *testing.T) {
    service := NewUserService(repo)
    service.validator.Validate() // Testing internal calls
}

// ✅ Good: Testing behavior
func TestUserService_CreateUser(t *testing.T) {
    result, err := service.CreateUser(validUser)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 3. Not Testing Error Cases
```go
// ❌ Bad: Only test happy path
func TestCreateUser(t *testing.T) {
    user, err := service.CreateUser(validUser)
    assert.NoError(t, err)
}

// ✅ Good: Test error cases too
func TestCreateUser_InvalidEmail(t *testing.T) {
    user, err := service.CreateUser(invalidUser)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), \"invalid email\")
}
```

## 🔗 Next Steps

- **Ada bug di tests?** → [Troubleshooting Guide](troubleshooting.md)
- **Mau setup CI/CD testing?** → [CI/CD Guide](../deployment/cicd.md)
- **Code quality issues?** → [Code Quality Guide](code-quality.md)

---

**💡 Pro Tips:**
- Run tests sering, jangan nunggu sampai banyak changes
- Test seharusnya cepat (< 1 detik per test)
- Write test dulu, baru write code (TDD)
- Keep tests simple dan readable
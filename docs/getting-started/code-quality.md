# Code Quality Guide

Panduan praktis untuk menjaga kualitas code Go. Tools dan best practices untuk code yang bersih dan maintainable.

## 🎯 Apa itu Code Quality?

**Code Quality** adalah seberapa baik code kita ditulis dari segi:
- **Readability**: Mudah dibaca dan dipahami
- **Maintainability**: Mudah diubah dan diperbaiki  
- **Consistency**: Mengikuti standard yang sama
- **Performance**: Efficient dan tidak boros resource

**Analogi sederhana:**
- Code = Tulisan tangan
- Code quality = Tulisan rapi, jelas, dan mudah dibaca orang lain
- Tools = Penggaris dan correction tape untuk rapikan tulisan

## 🛠️ Code Quality Tools

### 1. Formatting (go fmt)
Membuat code format consistent.

```bash
# Format semua Go files
make fmt

# Atau langsung
go fmt ./...

# Format specific file
go fmt internal/api/handlers/user_handler.go
```

**Apa yang dilakukan:**
- Indentation yang konsisten
- Spacing yang tepat
- Import grouping
- Bracket placement

**Before:**
```go
package main
import(
"fmt"
    "os"
)
func main(){
fmt.Println("Hello")
}
```

**After:**
```go
package main

import (
    \"fmt\"
    \"os\"
)

func main() {
    fmt.Println(\"Hello\")
}
```

### 2. Vetting (go vet)
Menemukan potential bugs dan suspicious code.

```bash
# Run vet pada semua packages
make vet

# Atau langsung
go vet ./...

# Vet specific package
go vet ./internal/api/handlers
```

**Yang dicek go vet:**
- Unreachable code
- Wrong printf format
- Incorrect struct tags
- Shadowed variables
- Unused variables

**Example Issues:**
```go
// ❌ go vet akan warning: printf format issue
fmt.Printf(\"User ID: %s\", userID) // userID is int, should be %d

// ❌ go vet akan warning: unreachable code
return user
fmt.Println(\"This will never execute\")

// ❌ go vet akan warning: struct tag issue
type User struct {
    ID   int64  `json:\"id\" db:\"user_id\"`  // Missing comma
    Name string `json\"name\"`                // Missing colon
}
```

### 3. Linting (golangci-lint)
Advanced code analysis dan style checking.

```bash
# Run linter
make lint

# Atau langsung (jika golangci-lint installed)
golangci-lint run

# Run specific linters
golangci-lint run --enable=gosec,gocritic

# Fix auto-fixable issues
golangci-lint run --fix
```

**Yang dicek linter:**
- Code style violations
- Security issues
- Performance issues
- Best practice violations
- Dead code
- Cyclomatic complexity

**Common Issues:**
```go
// ❌ Linter warning: variable name should be camelCase
var user_name string

// ✅ Fixed
var userName string

// ❌ Linter warning: should use constants for magic numbers
if len(password) < 8 {
    return errors.New(\"password too short\")
}

// ✅ Fixed
const MinPasswordLength = 8
if len(password) < MinPasswordLength {
    return errors.New(\"password too short\")
}
```

### 4. Module Tidying
Clean up dependencies.

```bash
# Tidy go modules
make tidy

# Atau langsung
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify
```

## 🚀 Quality Commands

### Basic Commands
```bash
# Format code
make fmt

# Run vet
make vet

# Run linter
make lint

# Tidy modules
make tidy

# Run all quality checks
make verify
```

### Advanced Commands
```bash
# Check for security issues
golangci-lint run --enable=gosec

# Check for performance issues
golangci-lint run --enable=prealloc,maligned

# Generate complexity report
gocyclo -over 10 .

# Check test coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## ⚙️ Configuration

### 1. Golangci-lint Config

**File: `.golangci.yml`**
```yaml
linters-settings:
  govet:
    check-shadowing: true
    settings:
      printf:
        funcs:
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Infof
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Warnf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Errorf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Fatalf
  
  revive:
    min-confidence: 0
  
  goimports:
    local-prefixes: github.com/yantology/golang_template
  
  gocyclo:
    min-complexity: 15
  
  maligned:
    suggest-new: true
  
  dupl:
    threshold: 100
  
  goconst:
    min-len: 2
    min-occurrences: 2

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - misspell
    - nakedret
    - rowserrcheck
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace

issues:
  exclude-rules:
    - path: _test\\.go
      linters:
        - gomnd
        - gocritic
    
    - path: internal/data/migrations/
      linters:
        - golint
        - stylecheck

run:
  deadline: 1m
  issues-exit-code: 1
  tests: true
```

### 2. Editor Integration

#### VS Code
**File: `.vscode/settings.json`**
```json
{
    \"go.formatTool\": \"goimports\",
    \"go.lintTool\": \"golangci-lint\",
    \"go.lintOnSave\": \"package\",
    \"go.formatOnSave\": true,
    \"go.useCodeSnippetsOnFunctionSuggest\": false,
    \"go.vetOnSave\": \"package\",
    \"go.buildOnSave\": \"package\",
    \"go.testOnSave\": false
}
```

#### Vim/Neovim
```vim
\" Install vim-go plugin
Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }

\" Auto format on save
let g:go_fmt_autosave = 1
let g:go_fmt_command = \"goimports\"

\" Run linters
let g:go_metalinter_autosave = 1
let g:go_metalinter_command = \"golangci-lint\"
```

## 📋 Pre-commit Setup

### Git Hooks
**File: `.git/hooks/pre-commit`**
```bash
#!/bin/sh

echo \"Running pre-commit quality checks...\"

# Run quality checks
make verify

if [ $? -ne 0 ]; then
    echo \"❌ Pre-commit checks failed. Please fix the issues before committing.\"
    exit 1
fi

echo \"✅ All pre-commit checks passed!\"
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

### Pre-commit with specific checks
```bash
#!/bin/sh

echo \"🔍 Running pre-commit checks...\"

# 1. Format code
echo \"📝 Formatting code...\"
make fmt
if [ $? -ne 0 ]; then
    echo \"❌ Code formatting failed\"
    exit 1
fi

# 2. Run vet
echo \"🔎 Running go vet...\"
make vet
if [ $? -ne 0 ]; then
    echo \"❌ go vet failed\"
    exit 1
fi

# 3. Run linter
echo \"🧹 Running linter...\"
make lint
if [ $? -ne 0 ]; then
    echo \"❌ Linter failed\"
    exit 1
fi

# 4. Run tests
echo \"🧪 Running tests...\"
make test
if [ $? -ne 0 ]; then
    echo \"❌ Tests failed\"
    exit 1
fi

echo \"✅ All pre-commit checks passed!\"
```

### Auto-fix Script
**File: `scripts/fix-quality.sh`**
```bash
#!/bin/bash

echo \"🔧 Auto-fixing code quality issues...\"

# Format code
echo \"📝 Formatting code...\"
go fmt ./...

# Fix imports
echo \"📦 Fixing imports...\"
goimports -w .

# Fix auto-fixable linter issues
echo \"🧹 Fixing linter issues...\"
golangci-lint run --fix

# Tidy modules
echo \"📚 Tidying modules...\"
go mod tidy

echo \"✅ Auto-fix completed!\"
echo \"⚠️  Please review changes and run 'make verify' to check remaining issues.\"
```

Make it executable:
```bash
chmod +x scripts/fix-quality.sh
```

## 🎯 Quality Best Practices

### 1. Naming Conventions
```go
// ✅ Good naming
type UserService struct {}
var userName string
const MaxRetryAttempts = 3

func (s *UserService) CreateUser(ctx context.Context, user *User) (*User, error) {}

// ❌ Bad naming
type userservice struct {}        // Should be UserService
var user_name string             // Should be userName
const MAX_RETRY_ATTEMPTS = 3     // Should be MaxRetryAttempts

func (s *userservice) create_user(user *User) (*User, error) {} // Should be CreateUser
```

### 2. Error Handling
```go
// ✅ Good error handling
func GetUser(id int64) (*User, error) {
    if id <= 0 {
        return nil, fmt.Errorf(\"invalid user ID: %d\", id)
    }
    
    user, err := repo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf(\"failed to get user: %w\", err)
    }
    
    return user, nil
}

// ❌ Bad error handling
func GetUser(id int64) *User {
    user, _ := repo.GetByID(id)  // Ignoring error
    return user
}
```

### 3. Function Size
```go
// ✅ Good: Small, focused function
func ValidateEmail(email string) error {
    if email == \"\" {
        return errors.New(\"email is required\")
    }
    
    if !emailRegex.MatchString(email) {
        return errors.New(\"invalid email format\")
    }
    
    return nil
}

// ❌ Bad: Too long, doing too many things
func ProcessUser(data map[string]interface{}) error {
    // 50+ lines of validation, processing, database operations, etc.
    // Should be broken down into smaller functions
}
```

### 4. Constants vs Magic Numbers
```go
// ✅ Good: Use constants
const (
    MinPasswordLength = 8
    MaxPasswordLength = 72
    DefaultPageSize   = 20
    MaxPageSize      = 100
)

func ValidatePassword(password string) error {
    if len(password) < MinPasswordLength {
        return fmt.Errorf(\"password must be at least %d characters\", MinPasswordLength)
    }
    return nil
}

// ❌ Bad: Magic numbers
func ValidatePassword(password string) error {
    if len(password) < 8 {  // What is 8? Why 8?
        return errors.New(\"password too short\")
    }
    return nil
}
```

## 📊 Quality Metrics

### Code Coverage
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage percentage
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

**Coverage Goals:**
- Unit tests: 80%+
- Critical business logic: 95%+
- Integration tests: 70%+

### Cyclomatic Complexity
```bash
# Install gocyclo
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# Check complexity (threshold: 10)
gocyclo -over 10 .

# Generate complexity report
gocyclo -avg .
```

**Complexity Guidelines:**
- 1-10: Simple, easy to test
- 11-20: Complex, harder to test
- 21+: Very complex, refactor recommended

### Technical Debt
```bash
# Check for TODO/FIXME comments
grep -r \"TODO\\|FIXME\\|HACK\" --exclude-dir=vendor .

# Check for deprecated functions
golangci-lint run --enable=staticcheck
```

## 🚨 Common Quality Issues

### 1. Inconsistent Error Handling
```go
// ❌ Inconsistent
func GetUser(id int64) (*User, error) {
    if id <= 0 {
        return nil, errors.New(\"Invalid ID\")  // Capital I
    }
    
    user, err := repo.GetByID(id)
    if err != nil {
        return nil, errors.New(\"user not found\")  // Lowercase u
    }
    
    return user, nil
}

// ✅ Consistent
func GetUser(id int64) (*User, error) {
    if id <= 0 {
        return nil, errors.New(\"invalid user ID\")  // Consistent lowercase
    }
    
    user, err := repo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf(\"failed to get user: %w\", err)  // Wrap error
    }
    
    return user, nil
}
```

### 2. Not Handling Context
```go
// ❌ Ignoring context
func GetUsers() ([]*User, error) {
    return repo.GetAll()  // No context, can't be cancelled
}

// ✅ Using context
func GetUsers(ctx context.Context) ([]*User, error) {
    return repo.GetAll(ctx)  // Can be cancelled, has timeout
}
```

### 3. Poor Variable Names
```go
// ❌ Poor names
func ProcessUsers(u []*User) error {
    for _, x := range u {
        if x.A {
            err := x.DoSomething()
            if err != nil {
                return err
            }
        }
    }
    return nil
}

// ✅ Clear names
func ActivateUsers(users []*User) error {
    for _, user := range users {
        if user.IsActive {
            err := user.SendWelcomeEmail()
            if err != nil {
                return fmt.Errorf(\"failed to send welcome email to user %d: %w\", user.ID, err)
            }
        }
    }
    return nil
}
```

## 🔧 Quality Automation

### GitHub Actions
**File: `.github/workflows/quality.yml`**
```yaml
name: Code Quality

on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Format
      run: |
        gofmt -s -w .
        git diff --exit-code
    
    - name: Vet
      run: go vet ./...
    
    - name: Lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
    
    - name: Test
      run: go test -race -coverprofile=coverage.out ./...
    
    - name: Coverage
      run: go tool cover -html=coverage.out -o coverage.html
```

### Makefile Integration
Update Makefile untuk include quality targets:
```makefile
.PHONY: quality
quality: fmt vet lint tidy test ## Run all quality checks

.PHONY: fix
fix: ## Auto-fix quality issues
	@./scripts/fix-quality.sh

.PHONY: quality-report
quality-report: ## Generate quality report
	@echo \"🔍 Generating quality report...\"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o reports/coverage.html
	@gocyclo -avg . > reports/complexity.txt
	@golangci-lint run --out-format=html > reports/lint.html
	@echo \"📊 Quality report generated in reports/\"
```

## 🔗 Next Steps

- **Setup CI/CD?** → [CI/CD Guide](../deployment/cicd.md)
- **Performance issues?** → [Performance Guide](../optimization/performance.md)
- **Need troubleshooting?** → [Troubleshooting Guide](troubleshooting.md)

---

**💡 Pro Tips:**
- Setup editor integration untuk real-time feedback
- Run `make verify` sebelum commit
- Use pre-commit hooks untuk enforce quality
- Review metrics regularly dan set improvement goals
# Architecture Overview

This document provides a comprehensive overview of the Go Backend Template's architecture, which follows Clean Architecture principles to ensure maintainability, testability, and scalability.

## 🏛️ Clean Architecture Principles

The template implements Clean Architecture with the following core principles:

### 1. Dependency Inversion
- **Inner layers define interfaces** that outer layers implement
- **Dependencies point inward** toward business logic
- **Business logic is independent** of external frameworks and databases

### 2. Layer Separation
- **Clear boundaries** between different layers
- **Single responsibility** for each layer
- **Interface-based communication** between layers

### 3. Testability
- **Easy to mock** dependencies using interfaces
- **Independent testing** of each layer
- **Fast unit tests** without external dependencies

## 📊 Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    🌐 API Layer                             │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Handlers   │  │ Middleware  │  │   Routes    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│                   💼 Business Layer                         │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Services   │  │ Validators  │  │  Workflows  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│                    💾 Data Layer                            │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Models    │  │Repositories │  │ Migrations  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
              │                           │
              ▼                           ▼
    ┌─────────────────┐        ┌─────────────────┐
    │   🌐 HTTP       │        │  🗄️ Database    │
    │   (Gin)         │        │  (PostgreSQL)   │
    └─────────────────┘        └─────────────────┘
```

## 📂 Project Structure

```
golang_template/
├── cmd/api/                    # 🚀 Application Entry Point
│   └── main.go                 #   Main application bootstrap
├── internal/                   # 🔒 Private Application Code
│   ├── api/                   # 🌐 API Layer
│   │   ├── handlers/          #   HTTP request handlers
│   │   ├── middleware/        #   Custom middleware
│   │   └── routes/            #   Route definitions
│   ├── business/              # 💼 Business Layer
│   │   ├── services/          #   Business logic services
│   │   └── validators/        #   Business rule validation
│   ├── data/                  # 💾 Data Layer
│   │   ├── models/            #   Database entity models
│   │   ├── repositories/      #   Data access repositories
│   │   └── migrations/        #   Database migrations
│   ├── config/                # ⚙️ Configuration
│   │   ├── config.go          #   Configuration structures
│   │   ├── database.go        #   Database configuration
│   │   ├── jwt.go             #   JWT configuration
│   │   ├── logger.go          #   Logger configuration
│   │   ├── server.go          #   Server configuration
│   │   └── viper.go           #   Viper setup
│   ├── pkg/                   # 📦 Internal Packages
│   │   ├── auth/              #   Authentication utilities
│   │   ├── database/          #   Database connection
│   │   ├── logger/            #   Logging utilities
│   │   └── utils/             #   Common utilities
│   └── server/                # 🖥️ Server Setup
│       └── server.go          #   HTTP server configuration
├── pkg/                       # 📦 Public Packages
│   ├── errors/                #   Custom error types
│   ├── response/              #   HTTP response utilities
│   └── validator/             #   Validation utilities
├── tests/                     # 🧪 Test Files
│   ├── e2e/                   #   End-to-end tests
│   ├── fixtures/              #   Test fixtures
│   ├── helpers/               #   Test helpers
│   ├── integration/           #   Integration tests
│   └── unit/                  #   Unit tests
└── docs/                      # 📚 Documentation
```

## 🌐 API Layer (`internal/api/`)

**Purpose**: Handle HTTP requests and responses, routing, and cross-cutting concerns.

### Responsibilities
- **HTTP Request/Response Handling**
- **Request Validation and Parsing**
- **Authentication and Authorization**
- **CORS**
- **Error Response Formatting**

### Components

#### Handlers (`internal/api/handlers/`)
- Process HTTP requests
- Coordinate with business services
- Transform business objects to HTTP responses
- Handle input validation

#### Middleware (`internal/api/middleware/`)
- Authentication/Authorization
- Request logging
- CORS handling
- Error recovery

#### Routes (`internal/api/routes/`)
- Define API endpoints
- Wire handlers to routes
- Apply middleware
- Group related endpoints

### Key Patterns
- **Dependency Injection** for services
- **Interface-based** service dependencies
- **Consistent response formats**
- **Proper error handling**

## 💼 Business Layer (`internal/business/`)

**Purpose**: Implement business logic, rules, and workflows.

### Responsibilities
- **Business Logic Implementation**
- **Business Rule Validation**
- **Workflow Orchestration**
- **Transaction Management**
- **External Service Integration**

### Components

#### Services (`internal/business/services/`)
- Implement business operations
- Orchestrate between repositories
- Handle complex business workflows
- Manage transactions

#### Validators (`internal/business/validators/`)
- Business rule validation
- Domain-specific validation logic
- Cross-field validation
- Business constraint checking

### Key Patterns
- **Service interfaces** for dependency inversion
- **Repository pattern** for data access
- **Transaction management** for data consistency
- **Error handling** with custom error types

## 💾 Data Layer (`internal/data/`)

**Purpose**: Handle data persistence, database operations, and data modeling.

### Responsibilities
- **Data Persistence**
- **Database Query Execution**
- **Data Model Definition**
- **Schema Management**
- **Database Transactions**

### Components

#### Models (`internal/data/models/`)
- Define database entities
- Model relationships
- Provide validation methods
- Handle data transformations

#### Repositories (`internal/data/repositories/`)
- Implement data access patterns
- Execute database queries
- Handle query parameters
- Manage database connections

#### Migrations (`internal/data/migrations/`)
- Database schema versioning
- Schema change management
- Data migration scripts
- Database initialization

### Key Patterns
- **Repository interfaces** for testability
- **Query builders** for dynamic queries
- **Connection pooling** for performance
- **Migration-based** schema management

## 🔧 Configuration Layer (`internal/config/`)

**Purpose**: Manage application configuration across different environments.

### Components
- **Configuration structs** for type safety
- **Viper integration** for flexible config sources
- **Environment variable mapping**
- **Default value management**
- **Configuration validation**

## 📦 Internal Packages (`internal/pkg/`)

**Purpose**: Provide shared utilities and common functionality.

### Components
- **Authentication utilities** (JWT, password hashing)
- **Database connection management**
- **Logging infrastructure**
- **Common utilities and helpers**

## 🏗️ Dependency Flow

```
HTTP Request → Handler → Service → Repository → Database
              ↓
         Middleware
              ↓
         Response ← Business Logic ← Data Access ← Query
```

### Interface Dependencies

```go
// Business layer defines interfaces
type UserService interface {
    CreateUser(ctx context.Context, user *User) (*User, error)
    GetUser(ctx context.Context, id int64) (*User, error)
}

// Data layer implements business interfaces
type UserRepository interface {
    Create(ctx context.Context, user *User) (*User, error)
    GetByID(ctx context.Context, id int64) (*User, error)
}

// API layer depends on business interfaces
type UserHandler struct {
    userService UserService // Interface, not concrete type
}
```

## 🧪 Testing Strategy

### Unit Tests
- **Business logic testing** without external dependencies
- **Repository testing** with database mocks
- **Handler testing** with service mocks

### Integration Tests
- **Database integration** with test database
- **Service integration** with real repositories
- **API integration** with test server

### End-to-End Tests
- **Full application testing** with real dependencies
- **User workflow testing**
- **API contract testing**

## 🔄 Data Flow Examples

### Creating a User

```
1. POST /api/v1/users
2. UserHandler.CreateUser()
3. UserService.CreateUser()
4. UserValidator.ValidateCreateUser()
5. UserRepository.Create()
6. Database INSERT
7. Return User object up the chain
8. JSON response to client
```

### Fetching Users with Pagination

```
1. GET /api/v1/users?page=1&limit=20
2. UserHandler.ListUsers()
3. UserService.ListUsers()
4. UserValidator.ValidateListParams()
5. UserRepository.List()
6. Database SELECT with LIMIT/OFFSET
7. Return Users + Pagination info
8. JSON response with data and pagination
```

## 🎯 Design Benefits

### Maintainability
- **Clear separation of concerns**
- **Easy to understand and modify**
- **Minimal coupling between layers**

### Testability
- **Independent layer testing**
- **Easy mocking with interfaces**
- **Fast unit tests**

### Scalability
- **Horizontal scaling ready**
- **Database connection pooling**
- **Stateless service design**

### Flexibility
- **Easy to swap implementations**
- **Framework independence**
- **Database independence**

## 🚀 Next Steps

- **Understand each layer in detail**: [API Layer](./api-layer.md), [Business Layer](./business-layer.md), [Data Layer](./data-layer.md)
- **Learn dependency injection patterns**: [Dependency Injection](./dependency-injection.md)
- **See practical implementation examples**: [Examples](../examples/)

---

This architecture provides a solid foundation for building scalable, maintainable Go applications while following industry best practices and clean code principles.
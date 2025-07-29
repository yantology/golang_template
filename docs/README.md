# Go Backend Template Documentation

Welcome to the Go Backend Template documentation! This template provides a clean, well-structured foundation for building scalable Go applications using Clean Architecture principles.

## 📚 Documentation Structure

### 🚀 [Getting Started](./getting-started/)
Essential information to get you up and running quickly.

- **[Setup](./getting-started/setup.md)** - Initial project setup and installation
- **[Configuration](./getting-started/configuration.md)** - Environment configuration and settings  
- **[Development](./getting-started/development.md)** - Development workflow and best practices

### 🏗️ [Architecture](./architecture/)
Deep dive into the clean architecture implementation.

- **[Overview](./architecture/overview.md)** - Clean architecture principles and structure
- **[API Layer](./architecture/api-layer.md)** - HTTP handlers, middleware, and routing
- **[Business Layer](./architecture/business-layer.md)** - Services, validators, and business logic
- **[Data Layer](./architecture/data-layer.md)** - Models, repositories, and database operations
- **[Dependency Injection](./architecture/dependency-injection.md)** - DI patterns and service wiring

### ⚙️ [Configuration](./configuration/)
Comprehensive configuration management guide.

- **[Overview](./configuration/overview.md)** - Configuration management with Viper
- **[Environment Variables](./configuration/environment-variables.md)** - All available env vars and their usage
- **[Adding New Configs](./configuration/adding-new-configs.md)** - How to add new configuration options
- **[Configuration Patterns](./configuration/configuration-patterns.md)** - Best practices for config management

### 🔌 [API Development](./api-development/)
Complete guide to building APIs with this template.

- **[Creating APIs](./api-development/creating-apis.md)** - Step-by-step guide to create new APIs
- **[Handler Patterns](./api-development/handler-patterns.md)** - Handler implementation patterns
- **[Validation](./api-development/validation.md)** - Request validation and error handling
- **[Middleware](./api-development/middleware.md)** - Creating custom middleware
- **[Authentication](./api-development/authentication.md)** - JWT auth implementation
- **[Testing APIs](./api-development/testing-apis.md)** - API testing strategies

### 📡 [API Specification](./api_spesification/)
Detailed API contract documentation for all endpoints.

- **[Overview](./api_spesification/README.md)** - Complete API specification with implementation guide
- **[Authentication](./api_spesification/auth/)** - Auth endpoints (login, register, logout, refresh)
- **[Profiles](./api_spesification/profiles/)** - User profile management endpoints
- **[Articles](./api_spesification/articles/)** - Content creation and management endpoints

### 💾 [Database](./database/)
Database operations and data management.

- **[Models](./database/models.md)** - Creating and managing database models
- **[Repositories](./database/repositories.md)** - Repository pattern implementation
- **[Migrations](./database/migrations.md)** - Database migration management
- **[Queries](./database/queries.md)** - Query building patterns with Squirrel

### 💡 [Examples](./examples/)
Practical examples and real-world implementations.

- **[Complete Feature](./examples/complete-feature.md)** - End-to-end example of adding a complete feature
- **[User Management](./examples/user-management.md)** - User CRUD operations example
- **[Article System](./examples/article-system.md)** - Content management example

## 🏛️ Architecture Overview

This template follows **Clean Architecture** principles with the following layers:

```
┌─────────────────────────────────────────────────────────────┐
│                      🌐 API Layer                           │
│                 (HTTP, Routes, Middleware)                  │
├─────────────────────────────────────────────────────────────┤
│                   💼 Business Layer                         │
│              (Services, Validators, Logic)                  │
├─────────────────────────────────────────────────────────────┤
│                    💾 Data Layer                            │
│            (Models, Repositories, Database)                 │
└─────────────────────────────────────────────────────────────┘
```

### Key Principles

- **Dependency Inversion** - Dependencies point inward toward business logic
- **Interface Segregation** - Small, focused interfaces
- **Single Responsibility** - Each layer has one reason to change
- **Testability** - Easy to mock and test each layer independently

## 🚀 Quick Start

1. **Clone and Setup**
   ```bash
   git clone <repository-url>
   cd golang_template
   cp .env.example .env
   ```

2. **Start Services**
   ```bash
   make db-up          # Start PostgreSQL + Adminer
   make migrate-up     # Run migrations
   ```

3. **Run Application**
   ```bash
   make dev           # Start Go application
   ```

4. **Verify Setup**
   - API: http://localhost:8080/health
   - Adminer: http://localhost:8081

## 🛠️ Development Commands

| Command | Description |
|---------|-------------|
| `make dev` | Start Go application |
| `make db-up` | Start database services |
| `make db-down` | Stop database services |
| `make migrate-up` | Run database migrations |
| `make test` | Run all tests |
| `make lint` | Run code linting |

## 📋 Project Structure

```
golang_template/
├── cmd/api/                    # Application entry point
├── internal/                   # Private application code
│   ├── api/                   # 🌐 API Layer
│   │   ├── handlers/          #   HTTP handlers
│   │   ├── middleware/        #   Custom middleware
│   │   └── routes/            #   Route definitions
│   ├── business/              # 💼 Business Layer
│   │   ├── services/          #   Business services
│   │   └── validators/        #   Business validation
│   ├── data/                  # 💾 Data Layer
│   │   ├── models/            #   Database models
│   │   ├── repositories/      #   Data access
│   │   └── migrations/        #   SQL migrations
│   ├── config/                # ⚙️ Configuration
│   └── pkg/                   # 📦 Internal packages
├── pkg/                       # Public packages
├── tests/                     # Test files
├── docs/                      # 📚 Documentation
└── docker-compose.yml         # Development services
```

## 🤝 Contributing

When contributing to this project:

1. Follow the established architecture patterns
2. Write tests for new functionality
3. Update documentation for any changes
4. Run `make verify` before committing
5. Use conventional commit messages

## 📖 Next Steps

- **New to the project?** Start with [Getting Started](./getting-started/)
- **Want to add an API?** Check [Creating APIs](./api-development/creating-apis.md)
- **Need to configure something?** See [Configuration](./configuration/)
- **Looking for examples?** Browse [Examples](./examples/)

## 🆘 Need Help?

- Check the [troubleshooting section](./getting-started/development.md#troubleshooting) in development docs
- Review the [examples](./examples/) for common patterns
- Examine the existing codebase for reference implementations

---

**Happy Coding!** 🚀
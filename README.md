# Go Backend Template

A clean Go backend template following clean architecture principles with simple PostgreSQL development setup.

## 🏗️ Architecture

This template follows **Clean Architecture** principles with clear separation of concerns:

- **`cmd/`** - Application entry points
- **`internal/`** - Private application code
  - **`api/`** - HTTP handlers, routes, middleware
  - **`business/`** - Business logic and services
  - **`data/`** - Data models, repositories, migrations
  - **`config/`** - Configuration management
  - **`pkg/`** - Internal packages (auth, database, logger, utils)
- **`pkg/`** - Public packages (errors, response, validator)
- **`tests/`** - Integration and E2E tests

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose (for database only)

### 1. Start Database

```bash
# Start PostgreSQL and Adminer with Docker
docker-compose up -d

# Check services are running
docker-compose ps
```

### 2. Setup Application

```bash
# Copy environment configuration
cp .env.example .env

# Install Go dependencies
go mod download

# Run database migrations (requires migrate CLI tool)
make migrate-up

# Start the Go application
go run cmd/api/main.go
```

The application will be available at: http://localhost:8080

## 📊 Services

| Service | URL | Credentials |
|---------|-----|-------------|
| **Go Application** | http://localhost:8080 | - |
| **Adminer** | http://localhost:8081 | Server: `postgres`, User: `postgres`, Password: `dev_password`, DB: `golang_template_dev` |
| **PostgreSQL** | localhost:5432 | Same as above |

## 🛠️ Development

### Environment Configuration

Edit `.env` file for your local development settings:
```bash
# Server
APP_SERVER_PORT=8080
APP_SERVER_ENV=development

# Database (matches docker-compose.yml)
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_USER=postgres
APP_DATABASE_PASSWORD=dev_password
APP_DATABASE_NAME=golang_template_dev

# JWT
APP_JWT_SECRET=your-super-secret-key-change-this-in-production-min-32-characters

# Logging
APP_LOGGER_LEVEL=debug
APP_LOGGER_FORMAT=text
```

### Database Management

#### Using Migrations (Recommended)
Run database migrations to set up your schema:
```bash
# Install migrate CLI tool
make tools

# Run migrations to create tables
make migrate-up

# Check migration status
make migrate-status
```

#### Using Adminer (Alternative)
Use Adminer web interface for database management:
1. Open http://localhost:8081
2. Login with:
   - **Server**: `postgres`
   - **Username**: `postgres`
   - **Password**: `dev_password`
   - **Database**: `golang_template_dev`

### API Endpoints

#### Public Endpoints
- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /api/v1/public/ping` - Ping test
- `GET /api/v1/public/version` - Version info

#### Authentication Endpoints
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh token

#### Protected Endpoints (require authentication)
**Users:**
- `GET /api/v1/users/me` - Get current user
- `PUT /api/v1/users/me` - Update current user
- `GET /api/v1/users/` - List users
- `GET /api/v1/users/:id` - Get user by ID

**Articles:**
- `POST /api/v1/articles/` - Create article
- `GET /api/v1/articles/` - List articles
- `GET /api/v1/articles/:id` - Get article
- `PUT /api/v1/articles/:id` - Update article
- `DELETE /api/v1/articles/:id` - Delete article
- `POST /api/v1/articles/:id/publish` - Publish article
- `POST /api/v1/articles/:id/archive` - Archive article

**Categories:**
- `GET /api/v1/categories/` - List categories
- `POST /api/v1/categories/` - Create category
- `GET /api/v1/categories/:id` - Get category
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

## 🛠️ Development Commands

```bash
# Database services
make db-up          # Start PostgreSQL + Adminer
make db-down        # Stop services  
make db-clean       # Stop and remove data
make db-connect     # Connect via psql

# Go application
make dev            # Run Go app

# Database migrations
make migrate-up     # Apply migrations
make migrate-down   # Rollback migrations
make migrate-create NAME=add_users_table  # Create new migration

# Testing
make test           # Run tests
make test-coverage  # Run with coverage
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration
```

## 🏛️ Architecture Patterns

### Clean Architecture Layers

1. **Entities** (`internal/data/models/`) - Business entities and rules
2. **Use Cases** (`internal/business/services/`) - Application-specific business rules
3. **Interface Adapters** (`internal/api/`) - Controllers, presenters, gateways
4. **Frameworks & Drivers** (`internal/pkg/`) - Database, web framework, external interfaces

### Key Principles

- **Dependency Inversion** - Dependencies point inward toward business logic
- **Interface Segregation** - Small, focused interfaces
- **Single Responsibility** - Each layer has one reason to change
- **Testability** - Easy to mock and test each layer independently

## 🔧 Customization

### Adding New Features

1. **Define Models** - Add to `internal/data/models/`
2. **Create Repository Interface** - Add to `internal/data/repositories/interfaces.go`
3. **Implement Repository** - Create in `internal/data/repositories/`
4. **Create Business Service** - Add to `internal/business/services/`
5. **Add API Handlers** - Create in `internal/api/handlers/`
6. **Define Routes** - Add to `internal/api/routes/`
7. **Create Migration Files** - Add SQL migration files to `internal/data/migrations/`

## 📚 Dependencies

### Core Dependencies
- **Gin** - HTTP web framework
- **PostgreSQL** - Primary database  
- **Viper** - Configuration management
- **Logrus** - Structured logging
- **JWT** - Authentication tokens
- **golang-migrate** - Database migrations

### Development Tools
- **Docker Compose** - PostgreSQL + Adminer for development
- **Adminer** - Database management web UI

## 🆘 Troubleshooting

### Common Issues

**Database Connection Issues**
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Reset database
docker-compose down -v && docker-compose up -d
```

**Migration Issues**
```bash
# Check migration status
make migrate-status

# Reset migrations if needed
make migrate-force VERSION=0
make migrate-up
```

**Application Issues**
```bash
# Check Go application logs
make dev

# Check environment configuration
cat .env
```

---

**Happy Coding!** 🚀
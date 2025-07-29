# Project Setup

This guide will help you set up the Go Backend Template for local development.

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

### Required
- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **Git** - [Install Git](https://git-scm.com/downloads)

### Optional (but recommended)
- **Make** - For using Makefile commands
- **migrate CLI** - For database migrations
- **golangci-lint** - For code linting

## 🚀 Quick Setup

### 1. Clone the Repository

```bash
git clone <your-repository-url>
cd golang_template
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install development tools (optional)
make tools
```

### 3. Environment Configuration

```bash
# Copy the environment template
cp .env.example .env

# Edit the configuration (optional - defaults work for development)
nano .env
```

### 4. Start Database Services

```bash
# Start PostgreSQL and Adminer with Docker
make db-up

# Verify services are running
docker-compose ps
```

### 5. Run Database Migrations

```bash
# Install migrate CLI if not already installed
make tools

# Run migrations to create database schema
make migrate-up

# Check migration status
make migrate-status
```

### 6. Start the Application

```bash
# Start the Go application
make dev

```

### 7. Verify Setup

Open your browser and check:

- **API Health Check**: http://localhost:8080/health
- **Adminer (Database UI)**: http://localhost:8081
- **API Documentation**: http://localhost:8080/api/v1/public/ping

## 🔧 Development Tools Installation

### Install All Tools at Once

```bash
make tools
```

This installs:
- `migrate` - Database migration tool
- `mockery` - Mock generation
- `golangci-lint` - Code linting

### Manual Installation


#### Migrate CLI
```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#### Golangci-lint
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### Mockery
```bash
go install github.com/vektra/mockery/v2@latest
```

## 🗄️ Database Setup Details

### Using Docker Compose (Recommended)

The project includes a `docker-compose.yml` file that sets up:

- **PostgreSQL 15** on port 5432
- **Adminer** (database management UI) on port 8081

```bash
# Start services
make db-up

# Stop services
make db-down

# Reset database (removes all data)
make db-clean
```

### Database Credentials

| Field | Value |
|-------|-------|
| Host | localhost |
| Port | 5432 |
| Database | golang_template_dev |
| Username | postgres |
| Password | dev_password |

### Accessing the Database

#### Via Adminer (Web Interface)
1. Open http://localhost:8081
2. Use the credentials above

#### Via Command Line
```bash
# Using make command
make db-connect

# Or directly with psql
psql "postgres://postgres:dev_password@localhost:5432/golang_template_dev"
```

## 🏗️ Project Structure Overview

After setup, your project structure will look like:

```
golang_template/
├── .env                       # Your environment configuration
├── cmd/api/main.go           # Application entry point
├── internal/                 # Private application code
│   ├── api/                 # HTTP layer
│   ├── business/            # Business logic layer
│   ├── data/                # Data access layer
│   ├── config/              # Configuration management
│   └── pkg/                 # Internal packages
├── pkg/                     # Public packages
├── tests/                   # Test files
├── docs/                    # Documentation
└── docker-compose.yml       # Database services
```

## ✅ Verification Checklist

After setup, verify everything is working:

- [ ] Go version 1.21+ installed (`go version`)
- [ ] Docker services running (`docker-compose ps`)
- [ ] Database accessible via Adminer (http://localhost:8081)
- [ ] Application starts without errors (`make dev`)
- [ ] Health check returns success (http://localhost:8080/health)
- [ ] Migrations applied successfully (`make migrate-status`)

## 🐛 Common Issues

### Port Already in Use

If you get port conflicts:

```bash
# Check what's using the ports
lsof -i :8080  # Go application
lsof -i :5432  # PostgreSQL
lsof -i :8081  # Adminer

# Stop conflicting services or modify ports in docker-compose.yml
```

### Database Connection Issues

```bash
# Check if PostgreSQL container is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Reset database completely
make db-clean && make db-up
```

### Migration Issues

```bash
# Check current migration version
make migrate-status

# Force to a specific version if needed
make migrate-force VERSION=0

# Rerun migrations
make migrate-up
```

### Permission Issues (Linux/Mac)

```bash
# If you get permission denied errors
sudo chown -R $USER:$USER .

# For Docker issues
sudo usermod -aG docker $USER
# Then logout and login again
```

## 🔄 Next Steps

Once setup is complete:

1. **Read the [Configuration Guide](./configuration.md)** to understand environment variables
2. **Follow the [Development Workflow](./development.md)** to learn the development process
3. **Explore the [Architecture Documentation](../architecture/overview.md)** to understand the project structure
4. **Try creating your first API** with the [API Development Guide](../api-development/creating-apis.md)

## 📝 Development Environment

For the best development experience:

1. **Use an IDE with Go support** (VS Code, GoLand, vim-go)
2. **Install Go extensions** for syntax highlighting and debugging
3. **Configure your editor** to run `go fmt` on save
4. **Set up git hooks** for pre-commit checks

### VS Code Extensions

Recommended extensions for VS Code:

- Go (official)
- REST Client (for testing APIs)
- Docker
- PostgreSQL (for database queries)

---

**Need help?** Check the [troubleshooting section](./development.md#troubleshooting) or refer to the main [README](../../README.md).
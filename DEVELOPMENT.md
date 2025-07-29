# Simple Development Setup

This document provides a simplified development setup using just PostgreSQL and Adminer with Docker, while running the Go application manually.

## 🚀 Quick Start

### 1. Start Database Services

```bash
# Start PostgreSQL and Adminer with Docker
docker-compose up -d

# Confirm services are running
docker-compose ps
```

### 2. Setup Go Application

```bash
# Copy environment configuration
cp .env.example .env

# Install Go dependencies
go mod download

# Start the Go application
go run cmd/api/main.go
```

## 📊 Services

| Service | URL | Credentials |
|---------|-----|-------------|
| **Adminer** | http://localhost:8081 | Server: `postgres`, User: `postgres`, Password: `dev_password`, Database: `golang_template_dev` |
| **PostgreSQL** | localhost:5432 | Same as above |
| **Go App** | http://localhost:8080 | - |

## 🛠️ Development Commands

```bash
# Database services
make db-up          # Start PostgreSQL + Adminer
make db-down        # Stop services  
make db-clean       # Stop and remove data
make db-connect     # Connect via psql

# Go application
make dev            # Run Go app

# Testing
make test           # Run tests
make test-coverage  # Run with coverage
```

## 🗄️ Database Management

### Using Adminer (Web UI)
1. Open http://localhost:8081
2. Login with:
   - **Server**: `postgres`
   - **Username**: `postgres`  
   - **Password**: `dev_password`
   - **Database**: `golang_template_dev`

### Using Command Line
```bash
# Connect to database
make db-connect

# Or connect directly
psql "postgres://postgres:dev_password@localhost:5432/golang_template_dev"
```


## 📂 Project Structure

```
golang_template/
├── cmd/api/main.go                 # Application entry point
├── internal/
│   ├── config/                     # Configuration management
│   ├── data/
│   │   ├── models/                 # Database models
│   │   ├── repositories/           # Data access layer
│   │   └── migrations/             # SQL migrations
│   ├── business/services/          # Business logic
│   ├── api/                        # HTTP handlers & routes
│   └── pkg/                        # Internal utilities
├── docker-compose.yml              # Database services only
├── .env.example                    # Environment template
└── Makefile                        # Development commands
```

## 🐛 Troubleshooting

### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker-compose ps

# Check database logs
make logs-db

# Reset database
make db-clean && make db-up
```

### Application Issues
```bash
# Check Go application logs in terminal
# Restart the Go application
make dev
```

### Port Conflicts
If ports 5432 or 8081 are in use:

```bash
# Check what's using the ports
lsof -i :5432
lsof -i :8081

# Stop conflicting services or modify docker-compose.yml
```

## 📝 Development Workflow

1. **Start database**: `make db-up`
2. **Run migrations**: `make migrate-up`  
3. **Start Go app**: `make dev` or `make dev-watch`
4. **Access Adminer**: http://localhost:8081
5. **Test API**: http://localhost:8080/health

This setup gives you a clean separation between infrastructure (Docker) and application development (local Go), making it easy to develop and debug while having a reliable database setup.
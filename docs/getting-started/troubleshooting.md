# Troubleshooting Guide

Solusi cepat untuk masalah-masalah umum saat development. Cari masalah kamu di sini dan ikuti langkah-langkah penyelesaiannya.

## 🚨 Quick Problem Finder

**Kamu mengalami masalah apa?**

| Masalah | Dokumen yang Tepat |
|---------|-------------------|
| 📱 **"Aplikasi tidak bisa jalan"** | → [App Won't Start](#-aplikasi-tidak-bisa-start) |
| 🗄️ **"Database error/connection failed"** | → [Database Issues](#-database-issues) |
| 🔌 **"Port already in use"** | → [Port Problems](#-port-sudah-dipakai) |
| 🧪 **"Tests gagal terus"** | → [Test Failures](#-test-failures) |
| 📦 **"Module/dependency error"** | → [Module Issues](#-module-issues) |
| 🔄 **"Migration stuck/error"** | → [Migration Problems](#-migration-problems) |
| 🐞 **"Aplikasi crash/panic"** | → [Runtime Errors](#-runtime-errors) |
| 🌐 **"API tidak response"** | → [API Issues](#-api-issues) |
| 💻 **"Development environment setup"** | → [Setup Problems](#-setup-problems) |

## 📱 Aplikasi Tidak Bisa Start

### Symptom: `make dev` error atau aplikasi langsung exit

#### Check 1: Port Already Used
```bash
# Cek apa yang pakai port 8080
lsof -i :8080

# Output example:
# COMMAND   PID   USER   FD   TYPE     DEVICE SIZE/OFF NODE NAME
# go      12345  user   8u   IPv6   0x1234567      0t0  TCP *:8080 (LISTEN)

# Kill process yang pakai port
kill -9 12345

# Coba start lagi
make dev
```

#### Check 2: Database Connection
```bash
# Pastikan database jalan
make db-up

# Test connection
make db-connect

# Kalau gagal, lihat logs
make db-logs
```

#### Check 3: Environment Variables
```bash
# Cek env variables
env | grep APP_

# Set environment variables yang diperlukan
export APP_DATABASE_HOST=localhost
export APP_DATABASE_PASSWORD=dev_password

# Coba start lagi
make dev
```

#### Check 4: Go Module Issues
```bash
# Clean dan download ulang
go mod tidy
go mod download
go clean -cache

# Coba start lagi
make dev
```

**💡 Masih gagal?** Lihat error message detail dan lanjut ke section yang sesuai.

## 🗄️ Database Issues

### Symptom: "connection refused", "database not found", "authentication failed"

#### Check 1: Database Container Status
```bash
# Cek status containers
docker-compose ps

# Should show postgres and adminer as "Up"
# Kalau tidak jalan:
make db-down && make db-up
```

#### Check 2: Database Credentials
```bash
# Test connection manual
psql "postgres://postgres:dev_password@localhost:5432/golang_template_dev"

# Kalau gagal login, reset database
make db-clean && make db-up
```

#### Check 3: Database Initialization
```bash
# Pastikan migration jalan
make migrate-status

# Kalau belum ada migration:
make migrate-up

# Kalau migration error:
make db-clean && make db-up && make migrate-up
```

#### Check 4: Port Conflicts
```bash
# Cek port 5432 (PostgreSQL)
lsof -i :5432

# Kalau ada conflict, stop service lain atau ganti port di docker-compose.yml
```

**🔧 Nuclear Option (Development Only):**
```bash
# Reset semua database dan mulai fresh
make db-clean
docker system prune -f
make db-up
make migrate-up
```

## 🔌 Port Sudah Dipakai

### Symptom: "port already in use", "bind: address already in use"

#### Port 8080 (Aplikasi)
```bash
# Find what's using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use killall
killall go

# Restart app
make dev
```

#### Port 5432 (Database)
```bash
# Check port 5432
lsof -i :5432

# Stop other PostgreSQL services
sudo service postgresql stop

# Or kill specific process
kill -9 <PID>

# Restart database
make db-down && make db-up
```

#### Port 8081 (Adminer)
```bash
# Check port 8081
lsof -i :8081

# Usually not a problem, but if needed:
kill -9 <PID>
make db-down && make db-up
```

**🔧 Alternative: Change Ports**
Edit `docker-compose.yml`:
```yaml
services:
  postgres:
    ports:
      - "5433:5432"  # Change from 5432 to 5433
  
  adminer:
    ports:
      - "8082:8080"  # Change from 8081 to 8082
```

## 🧪 Test Failures

### Symptom: `make test` fails, tests timeout, assertion errors

#### Check 1: Database for Integration Tests
```bash
# Setup test database
createdb golang_template_test

# Run migrations on test DB
migrate -path ./internal/data/migrations \
        -database "postgres://postgres:password@localhost/golang_template_test" up

# Run tests again
make test
```

#### Check 2: Clean Test Data
```bash
# Reset test database
dropdb golang_template_test
createdb golang_template_test

# Run migrations
make migrate-up  # Adjust for test DB

# Run tests
make test
```

#### Check 3: Test-Specific Issues
```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v ./internal/business/services -run TestUserService_CreateUser

# Run with race detection
go test -race ./...
```

#### Check 4: Mock Issues
```bash
# Regenerate mocks (if using mockery)
go generate ./...

# Or manually regenerate
mockery --all --output ./mocks
```

**📋 Common Test Issues:**

| Error | Solution |
|-------|----------|
| "database connection failed" | Setup test database |
| "table doesn't exist" | Run migrations on test DB |
| "mock expectations not met" | Check mock setup in test |
| "test timeout" | Increase timeout or check for deadlocks |

## 📦 Module Issues

### Symptom: "module not found", "version conflict", "checksum mismatch"

#### Check 1: Clean Module Cache
```bash
# Clean module cache
go clean -modcache

# Re-download modules
go mod download

# Tidy modules
go mod tidy
```

#### Check 2: Proxy Issues
```bash
# Disable Go proxy temporarily
export GOPROXY=direct

# Download modules
go mod download

# Reset proxy
unset GOPROXY
```

#### Check 3: Version Conflicts
```bash
# Check module dependencies
go mod graph

# Update specific module
go get -u github.com/gin-gonic/gin

# Update all modules
go get -u ./...
```

#### Check 4: Go Version
```bash
# Check Go version
go version

# Should be Go 1.21+
# If not, update Go from https://golang.org/dl/
```

**🔧 Nuclear Option:**
```bash
# Remove go.sum and re-download
rm go.sum
go mod tidy
go mod download
```

## 🔄 Migration Problems

### Symptom: Migration stuck, version mismatch, SQL errors

#### Check 1: Migration Status
```bash
# Check current status
make migrate-status

# Output should show current version
# If stuck, check database logs
make db-logs
```

#### Check 2: Force Migration Version
```bash
# Force to specific version (DANGEROUS!)
make migrate-force VERSION=1

# Then run migrations normally
make migrate-up
```

#### Check 3: SQL Syntax Errors
```bash
# Check migration files manually
cat internal/data/migrations/000001_*.sql

# Test SQL manually
psql "postgres://postgres:dev_password@localhost:5432/golang_template_dev"
\i internal/data/migrations/000001_create_table.up.sql
```

#### Check 4: Migration Lock
```bash
# Check for migration locks in database
psql "postgres://..." -c "SELECT * FROM schema_migrations;"

# If locked, unlock (development only)
psql "postgres://..." -c "UPDATE schema_migrations SET dirty = false;"
```

**🚨 Development Reset (LOSES DATA):**
```bash
# Complete reset
make db-clean
make db-up
make migrate-up
```

## 🐞 Runtime Errors

### Symptom: Panic, crashes, unexpected behavior

#### Check 1: Logs
```bash
# Run with debug logging
APP_LOGGER_LEVEL=debug make dev

# Check for panic stack traces
# Look for specific error messages
```

#### Check 2: Race Conditions
```bash
# Run with race detector
go run -race cmd/api/main.go

# Or in tests
go test -race ./...
```

#### Check 3: Memory Issues
```bash
# Check memory usage
go tool pprof http://localhost:8080/debug/pprof/heap

# Check for goroutine leaks
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

#### Check 4: Configuration
```bash
# Print current config
APP_LOGGER_LEVEL=debug make dev

# Look for config validation errors in startup logs
```

## 🌐 API Issues

### Symptom: 404 errors, 500 errors, no response

#### Check 1: Server Status
```bash
# Check if server is running
curl http://localhost:8080/health

# Should return: {"status":"healthy"}
```

#### Check 2: Routes
```bash
# Check available routes (if you have route listing)
curl http://localhost:8080/debug/routes

# Test specific endpoint
curl -v http://localhost:8080/api/v1/ping
```

#### Check 3: Request Format
```bash
# Test with proper headers
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test"}'
```

#### Check 4: Authentication
```bash
# Test protected endpoint
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/protected
```

## 💻 Setup Problems

### Symptom: Can't install dependencies, Docker issues, environment problems

#### Check 1: Go Installation
```bash
# Check Go version
go version

# Should be 1.21+
# Download from: https://golang.org/dl/
```

#### Check 2: Docker Installation
```bash
# Check Docker
docker --version
docker-compose --version

# Test Docker
docker run hello-world
```

#### Check 3: Required Tools
```bash
# Install missing tools
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Or use make
make tools
```

#### Check 4: Permissions
```bash
# Fix Docker permissions (Linux)
sudo usermod -aG docker $USER
newgrp docker

# Test without sudo
docker ps
```

## 🔧 Emergency Commands

### Development Reset (Nuclear Option)
```bash
# Stop everything
make db-down
killall go

# Clean everything
docker system prune -f
go clean -cache -modcache

# Fresh start
go mod tidy
make db-up
make migrate-up
make dev
```

### Quick Health Check
```bash
# Check all services
echo "🐘 Database:"
docker-compose ps postgres

echo "🖥️  Application:"
curl -s http://localhost:8080/health | jq .

echo "🌐 Adminer:"
curl -s -o /dev/null -w "%{http_code}" http://localhost:8081

echo "🔄 Migration status:"
make migrate-status
```

### Debug Information Gathering
```bash
# Gather debug info
echo "=== SYSTEM INFO ===" > debug.log
go version >> debug.log
docker --version >> debug.log
docker-compose --version >> debug.log

echo "=== PROCESSES ===" >> debug.log
lsof -i :8080 >> debug.log
lsof -i :5432 >> debug.log

echo "=== DOCKER STATUS ===" >> debug.log
docker-compose ps >> debug.log

echo "=== LOGS ===" >> debug.log
make db-logs >> debug.log 2>&1

cat debug.log
```

## 📞 Getting Help

### Before Asking for Help
1. ✅ **Check this troubleshooting guide**
2. ✅ **Try the Emergency Commands**
3. ✅ **Gather debug information**
4. ✅ **Note exact error messages**
5. ✅ **Try to reproduce the issue**

### What to Include When Asking for Help
- **Operating System** (macOS, Linux, Windows)
- **Go version** (`go version`)
- **Docker version** (`docker --version`)
- **Exact error message** (copy-paste, don't paraphrase)
- **Steps to reproduce**
- **What you already tried**

### Useful Log Commands
```bash
# Application logs (if running in background)
journalctl -u your-app -f

# Docker logs
docker-compose logs -f

# System logs (Linux)
tail -f /var/log/syslog

# macOS logs
log stream --predicate 'process CONTAINS "go"'
```

## 🏥 Health Checks

### Quick System Check
```bash
#!/bin/bash
echo "🔍 Running system health check..."

# Check Go
if go version &> /dev/null; then
    echo "✅ Go is installed: $(go version)"
else
    echo "❌ Go is not installed"
fi

# Check Docker
if docker --version &> /dev/null; then
    echo "✅ Docker is installed: $(docker --version)"
else
    echo "❌ Docker is not installed"
fi

# Check database
if docker-compose ps postgres | grep -q "Up"; then
    echo "✅ Database is running"
else
    echo "❌ Database is not running"
fi

# Check application
if curl -s http://localhost:8080/health | grep -q "healthy"; then
    echo "✅ Application is healthy"
else
    echo "❌ Application is not responding"
fi

echo "🎯 Health check completed!"
```

### Monitoring During Development
```bash
# Watch logs in real-time
watch -n 2 'curl -s http://localhost:8080/health | jq .'

# Monitor database connections
watch -n 5 'docker exec postgres_container psql -U postgres -c "SELECT count(*) FROM pg_stat_activity;"'

# Monitor system resources
htop
```

## 🔗 Quick Navigation

**Butuh bantuan untuk topik spesifik?**

| Topic | Document |
|-------|----------|
| 🌅 **Workflow harian** | → [Daily Workflow Guide](daily-workflow.md) |
| 🗄️ **Database local** | → [Database Development Guide](database-development.md) |
| 🏭 **Database production** | → [Database Production Guide](database-production.md) |
| 🧪 **Testing issues** | → [Testing Guide](testing-guide.md) |
| ✨ **Code quality** | → [Code Quality Guide](code-quality.md) |

---

**💡 Pro Tips:**
- Bookmark halaman ini untuk reference cepat
- Save emergency commands di notes
- Setup health check script untuk monitoring
- Document solutions yang kamu temukan untuk future reference
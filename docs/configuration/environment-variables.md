# Environment Variables Reference

Complete reference for all environment variables supported by the Go Backend Template.

## 📋 Environment Variable Convention

- **Prefix**: All environment variables use the `APP_` prefix
- **Naming**: Use UPPERCASE with underscores for nested values
- **Mapping**: `APP_SERVER_PORT` maps to `server.port` in configuration
- **Types**: Automatic type conversion for numbers, booleans, and durations

## 🌍 Environment Priority

Environment variables override all other configuration sources:

1. **Environment Variables** (highest priority)
2. Configuration files (YAML/JSON)
3. Default values (lowest priority)

## 🖥️ Server Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `APP_SERVER_PORT` | string | `"8080"` | HTTP server port |
| `APP_SERVER_HOST` | string | `"0.0.0.0"` | Server bind address |
| `APP_SERVER_ENV` | string | `"development"` | Environment (development/staging/production/test) |
| `APP_SERVER_READ_TIMEOUT` | duration | `"10s"` | HTTP read timeout |
| `APP_SERVER_WRITE_TIMEOUT` | duration | `"10s"` | HTTP write timeout |
| `APP_SERVER_IDLE_TIMEOUT` | duration | `"60s"` | HTTP idle timeout |
| `APP_SERVER_SHUTDOWN_TIMEOUT` | duration | `"30s"` | Graceful shutdown timeout |
| `APP_SERVER_ENABLE_CORS` | bool | `true` | Enable CORS middleware |
| `APP_SERVER_CORS_ORIGINS` | []string | `["*"]` | Allowed CORS origins (comma-separated) |

### Example Server Configuration

```bash
# Development
APP_SERVER_PORT=8080
APP_SERVER_HOST=localhost
APP_SERVER_ENV=development
APP_SERVER_CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# Production
APP_SERVER_PORT=8080
APP_SERVER_HOST=0.0.0.0
APP_SERVER_ENV=production
APP_SERVER_CORS_ORIGINS=https://yourdomain.com,https://api.yourdomain.com
```

## 🗄️ Database Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `APP_DATABASE_TYPE` | string | `"postgres"` | Database type (postgres only) |
| `APP_DATABASE_HOST` | string | `"localhost"` | Database host |
| `APP_DATABASE_PORT` | string | `"5432"` | Database port |
| `APP_DATABASE_USER` | string | `"postgres"` | Database username |
| `APP_DATABASE_PASSWORD` | string | `"postgres"` | Database password |
| `APP_DATABASE_NAME` | string | `"golang_template"` | Database name |
| `APP_DATABASE_SSLMODE` | string | `"disable"` | PostgreSQL SSL mode (disable/require/verify-ca/verify-full) |
| `APP_DATABASE_MAX_OPEN_CONNS` | int | `25` | Maximum open connections |
| `APP_DATABASE_MAX_IDLE_CONNS` | int | `5` | Maximum idle connections |
| `APP_DATABASE_MAX_LIFETIME` | duration | `"300s"` | Connection maximum lifetime |
| `APP_DATABASE_MIGRATION_PATH` | string | `"./internal/data/migrations"` | Migration files path |

### Example Database Configuration

```bash
# Development (Docker)
APP_DATABASE_TYPE=postgres
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_USER=postgres
APP_DATABASE_PASSWORD=dev_password
APP_DATABASE_NAME=golang_template_dev
APP_DATABASE_SSLMODE=disable

# Production
APP_DATABASE_TYPE=postgres
APP_DATABASE_HOST=prod-db.example.com
APP_DATABASE_PORT=5432
APP_DATABASE_USER=app_user
APP_DATABASE_PASSWORD=secure_password_here
APP_DATABASE_NAME=golang_template_prod
APP_DATABASE_SSLMODE=require
APP_DATABASE_MAX_OPEN_CONNS=50
APP_DATABASE_MAX_IDLE_CONNS=10
```

## 🔐 JWT Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `APP_JWT_SECRET` | string | `"your-super-secret-key..."` | JWT signing secret (minimum 32 characters) |
| `APP_JWT_ACCESS_TOKEN_TTL` | duration | `"15m"` | Access token lifetime |
| `APP_JWT_REFRESH_TOKEN_TTL` | duration | `"24h"` | Refresh token lifetime |
| `APP_JWT_ISSUER` | string | `"golang-template"` | JWT issuer claim |
| `APP_JWT_AUDIENCE` | string | `"golang-template-users"` | JWT audience claim |
| `APP_JWT_ALGORITHM` | string | `"HS256"` | JWT signing algorithm (HS256/HS384/HS512) |

### Example JWT Configuration

```bash
# Development
APP_JWT_SECRET=development-secret-key-minimum-32-characters-long
APP_JWT_ACCESS_TOKEN_TTL=15m
APP_JWT_REFRESH_TOKEN_TTL=24h

# Production
APP_JWT_SECRET=production-super-secure-jwt-secret-key-minimum-32-characters-long
APP_JWT_ACCESS_TOKEN_TTL=30m
APP_JWT_REFRESH_TOKEN_TTL=7d
APP_JWT_ISSUER=yourapp
APP_JWT_AUDIENCE=yourapp-users
```

## 📝 Logger Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `APP_LOGGER_LEVEL` | string | `"info"` | Log level (debug/info/warn/error/fatal) |
| `APP_LOGGER_FORMAT` | string | `"json"` | Log format (json/text) |
| `APP_LOGGER_OUTPUT` | string | `"stdout"` | Log output (stdout/stderr/file) |
| `APP_LOGGER_ENABLE_CALLER` | bool | `true` | Include caller information in logs |
| `APP_LOGGER_ENABLE_STACKTRACE` | bool | `false` | Include stack trace in error logs |

### Example Logger Configuration

```bash
# Development
APP_LOGGER_LEVEL=debug
APP_LOGGER_FORMAT=text
APP_LOGGER_OUTPUT=stdout

# Production
APP_LOGGER_LEVEL=info
APP_LOGGER_FORMAT=json
APP_LOGGER_OUTPUT=stdout
APP_LOGGER_ENABLE_STACKTRACE=true
```

## 🔧 Extended Configuration Examples

### Redis Configuration (Optional)

```bash
# Redis Configuration
APP_REDIS_HOST=localhost
APP_REDIS_PORT=6379
APP_REDIS_PASSWORD=
APP_REDIS_DATABASE=0
APP_REDIS_MAX_RETRIES=3
APP_REDIS_DIAL_TIMEOUT=5s
APP_REDIS_READ_TIMEOUT=3s
APP_REDIS_WRITE_TIMEOUT=3s
APP_REDIS_POOL_SIZE=10
APP_REDIS_MIN_IDLE_CONNS=5
APP_REDIS_MAX_CONN_AGE=30m
APP_REDIS_POOL_TIMEOUT=4s
APP_REDIS_IDLE_TIMEOUT=5m
APP_REDIS_IDLE_CHECK_FREQ=1m
```

### Email Configuration (Optional)

```bash
# Email Configuration
APP_EMAIL_SMTP_HOST=smtp.gmail.com
APP_EMAIL_SMTP_PORT=587
APP_EMAIL_USERNAME=your-email@gmail.com
APP_EMAIL_PASSWORD=your-app-password
APP_EMAIL_FROM_ADDRESS=noreply@yourdomain.com
APP_EMAIL_FROM_NAME=Your App Name
APP_EMAIL_USE_TLS=true
```

### Storage Configuration (Optional)

```bash
# Storage Configuration
APP_STORAGE_TYPE=local
APP_STORAGE_LOCAL_PATH=./uploads
APP_STORAGE_MAX_FILE_SIZE=52428800
APP_STORAGE_ALLOWED_EXTENSIONS=.jpg,.jpeg,.png,.gif,.pdf

# S3 Storage (when APP_STORAGE_TYPE=s3)
APP_STORAGE_S3_BUCKET=your-bucket-name
APP_STORAGE_S3_REGION=us-east-1
APP_STORAGE_S3_ACCESS_KEY=your-access-key
APP_STORAGE_S3_SECRET_KEY=your-secret-key
```

## 🌍 Environment-Specific Examples

### Development Environment

```bash
# .env.development
APP_SERVER_ENV=development
APP_SERVER_PORT=8080
APP_SERVER_HOST=localhost
APP_SERVER_CORS_ORIGINS=http://localhost:3000

APP_DATABASE_HOST=localhost
APP_DATABASE_NAME=golang_template_dev
APP_DATABASE_PASSWORD=dev_password
APP_DATABASE_SSLMODE=disable

APP_JWT_SECRET=development-secret-key-minimum-32-characters-long

APP_LOGGER_LEVEL=debug
APP_LOGGER_FORMAT=text
```

### Staging Environment

```bash
# .env.staging
APP_SERVER_ENV=staging
APP_SERVER_PORT=8080
APP_SERVER_HOST=0.0.0.0
APP_SERVER_CORS_ORIGINS=https://staging.yourdomain.com

APP_DATABASE_HOST=staging-db.yourdomain.com
APP_DATABASE_NAME=golang_template_staging
APP_DATABASE_PASSWORD=${STAGING_DB_PASSWORD}
APP_DATABASE_SSLMODE=require

APP_JWT_SECRET=${STAGING_JWT_SECRET}

APP_LOGGER_LEVEL=info
APP_LOGGER_FORMAT=json
```

### Production Environment

```bash
# .env.production
APP_SERVER_ENV=production
APP_SERVER_PORT=8080
APP_SERVER_HOST=0.0.0.0
APP_SERVER_CORS_ORIGINS=https://yourdomain.com

APP_DATABASE_HOST=${PROD_DB_HOST}
APP_DATABASE_PORT=5432
APP_DATABASE_USER=${PROD_DB_USER}
APP_DATABASE_PASSWORD=${PROD_DB_PASSWORD}
APP_DATABASE_NAME=${PROD_DB_NAME}
APP_DATABASE_SSLMODE=require
APP_DATABASE_MAX_OPEN_CONNS=50
APP_DATABASE_MAX_IDLE_CONNS=10

APP_JWT_SECRET=${PROD_JWT_SECRET}
APP_JWT_ACCESS_TOKEN_TTL=30m
APP_JWT_REFRESH_TOKEN_TTL=7d

APP_LOGGER_LEVEL=warn
APP_LOGGER_FORMAT=json
```

## 🧪 Testing Environment

```bash
# .env.test
APP_SERVER_ENV=test
APP_SERVER_PORT=0

APP_DATABASE_HOST=localhost
APP_DATABASE_NAME=golang_template_test
APP_DATABASE_PASSWORD=test_password

APP_JWT_SECRET=test-secret-key-minimum-32-characters-long

APP_LOGGER_LEVEL=error
APP_LOGGER_FORMAT=text
```

## 🔍 Type Conversion Examples

### Duration Values

```bash
# Valid duration formats
APP_SERVER_READ_TIMEOUT=10s        # 10 seconds
APP_SERVER_READ_TIMEOUT=5m         # 5 minutes
APP_SERVER_READ_TIMEOUT=1h         # 1 hour
APP_SERVER_READ_TIMEOUT=2h30m      # 2 hours 30 minutes
APP_JWT_ACCESS_TOKEN_TTL=15m       # 15 minutes
APP_JWT_REFRESH_TOKEN_TTL=24h      # 24 hours
APP_JWT_REFRESH_TOKEN_TTL=7d       # 7 days (not supported by time.Duration, use 168h)
```

### Boolean Values

```bash
# Valid boolean formats (case insensitive)
APP_SERVER_ENABLE_CORS=true
APP_SERVER_ENABLE_CORS=TRUE
APP_SERVER_ENABLE_CORS=false
APP_SERVER_ENABLE_CORS=FALSE
APP_LOGGER_ENABLE_CALLER=1         # 1 = true
APP_LOGGER_ENABLE_CALLER=0         # 0 = false
```

### Array Values

```bash
# Comma-separated values for arrays
APP_SERVER_CORS_ORIGINS=http://localhost:3000,http://localhost:3001,https://yourdomain.com
APP_STORAGE_ALLOWED_EXTENSIONS=.jpg,.jpeg,.png,.gif,.pdf,.doc,.docx
```

### Number Values

```bash
# Integer values
APP_SERVER_PORT=8080
APP_DATABASE_MAX_OPEN_CONNS=25

# Size values (bytes)
APP_STORAGE_MAX_FILE_SIZE=52428800    # 50MB in bytes
APP_SERVER_MAX_REQUEST_SIZE=10485760  # 10MB in bytes
```

## 🔒 Security Considerations

### Sensitive Variables

Mark these variables as sensitive in your deployment system:

```bash
# Secrets - never commit these values
APP_DATABASE_PASSWORD=
APP_JWT_SECRET=
APP_REDIS_PASSWORD=
APP_EMAIL_PASSWORD=
APP_STORAGE_S3_ACCESS_KEY=
APP_STORAGE_S3_SECRET_KEY=
```

### Production Security

```bash
# Production security requirements
APP_JWT_SECRET=              # Minimum 64 characters
APP_DATABASE_SSLMODE=require # Always use SSL in production
APP_SERVER_ENV=production    # Never use development in production
```

## 📚 Usage Examples

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_SERVER_PORT=8080
      - APP_DATABASE_HOST=postgres
      - APP_DATABASE_PASSWORD=dev_password
    depends_on:
      - postgres
  
  postgres:
    image: postgres:15
    environment:
      - POSTGRES_PASSWORD=dev_password
      - POSTGRES_DB=golang_template_dev
```

### Kubernetes ConfigMap

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_SERVER_PORT: "8080"
  APP_SERVER_ENV: "production"
  APP_DATABASE_HOST: "postgres-service"
  APP_LOGGER_LEVEL: "info"
  APP_LOGGER_FORMAT: "json"
```

### Kubernetes Secret

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
data:
  APP_DATABASE_PASSWORD: <base64-encoded-password>
  APP_JWT_SECRET: <base64-encoded-secret>
```

## 🛠️ Environment Variable Validation

The application validates environment variables at startup:

```bash
# This will fail validation
APP_JWT_SECRET=short          # Too short (< 32 characters)
APP_SERVER_PORT=99999         # Invalid port range
APP_DATABASE_MAX_OPEN_CONNS=0 # Must be positive
APP_LOGGER_LEVEL=invalid      # Invalid log level
```

## 🚀 Next Steps

- **Learn configuration patterns**: [Configuration Patterns](./configuration-patterns.md)
- **Understand how to add new configs**: [Adding New Configs](./adding-new-configs.md)
- **See the complete setup**: [Getting Started](../getting-started/setup.md)

---

Use this reference to properly configure your application across different environments while maintaining security and type safety.
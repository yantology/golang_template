# Configuration Guide

This guide explains how to configure the Go Backend Template for different environments and use cases.

## 📋 Configuration Overview

The template uses **Viper** for configuration management, supporting:

- Environment variables (`.env` files)
- YAML configuration files
- Command-line flags
- Default values

### Configuration Priority (highest to lowest)

1. Environment variables with `APP_` prefix
2. YAML configuration files
3. Default values set in code

## 🔧 Environment Variables

### Quick Setup

```bash
# Copy the template
cp .env.example .env

# Edit your configuration
nano .env
```

### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_SERVER_PORT` | `8080` | HTTP server port |
| `APP_SERVER_HOST` | `0.0.0.0` | Server bind address |
| `APP_SERVER_ENV` | `development` | Environment (development/staging/production) |
| `APP_SERVER_READ_TIMEOUT` | `10s` | HTTP read timeout |
| `APP_SERVER_WRITE_TIMEOUT` | `10s` | HTTP write timeout |
| `APP_SERVER_IDLE_TIMEOUT` | `60s` | HTTP idle timeout |
| `APP_SERVER_SHUTDOWN_TIMEOUT` | `30s` | Graceful shutdown timeout |

**Example:**
```bash
APP_SERVER_PORT=3000
APP_SERVER_HOST=127.0.0.1
APP_SERVER_ENV=production
```

### Database Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_DATABASE_TYPE` | `postgres` | Database type (postgres only) |
| `APP_DATABASE_HOST` | `localhost` | Database host |
| `APP_DATABASE_PORT` | `5432` | Database port |
| `APP_DATABASE_USER` | `postgres` | Database username |
| `APP_DATABASE_PASSWORD` | `postgres` | Database password |
| `APP_DATABASE_NAME` | `golang_template` | Database name |
| `APP_DATABASE_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `APP_DATABASE_MAX_OPEN_CONNS` | `25` | Maximum open connections |
| `APP_DATABASE_MAX_IDLE_CONNS` | `5` | Maximum idle connections |
| `APP_DATABASE_MAX_LIFETIME` | `300s` | Connection max lifetime |

**Example:**
```bash
APP_DATABASE_TYPE=postgres
APP_DATABASE_HOST=db.example.com
APP_DATABASE_PORT=5432
APP_DATABASE_USER=app_user
APP_DATABASE_PASSWORD=secure_password
APP_DATABASE_NAME=production_db
APP_DATABASE_SSLMODE=require
```

### JWT Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_JWT_SECRET` | `your-super-secret-key...` | JWT signing secret (min 32 chars) |
| `APP_JWT_ACCESS_TOKEN_TTL` | `15m` | Access token lifetime |
| `APP_JWT_REFRESH_TOKEN_TTL` | `24h` | Refresh token lifetime |
| `APP_JWT_ISSUER` | `golang-template` | JWT issuer |
| `APP_JWT_AUDIENCE` | `golang-template-users` | JWT audience |
| `APP_JWT_ALGORITHM` | `HS256` | JWT signing algorithm |

**Example:**
```bash
APP_JWT_SECRET=my-super-secret-jwt-key-minimum-32-characters-long
APP_JWT_ACCESS_TOKEN_TTL=30m
APP_JWT_REFRESH_TOKEN_TTL=7d
```

### Logger Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_LOGGER_LEVEL` | `info` | Log level (debug/info/warn/error) |
| `APP_LOGGER_FORMAT` | `json` | Log format (json/text) |
| `APP_LOGGER_OUTPUT` | `stdout` | Log output (stdout/stderr/file) |
| `APP_LOGGER_ENABLE_CALLER` | `true` | Include caller info in logs |
| `APP_LOGGER_ENABLE_STACKTRACE` | `false` | Include stack trace in error logs |

**Example:**
```bash
APP_LOGGER_LEVEL=debug
APP_LOGGER_FORMAT=text
APP_LOGGER_OUTPUT=stdout
```

## 📄 YAML Configuration

You can also use YAML files for configuration. Create a `config.yaml` file:

```yaml
# config.yaml
server:
  port: "8080"
  host: "0.0.0.0"
  env: "development"
  read_timeout: "10s"
  write_timeout: "10s"
  idle_timeout: "60s"
  shutdown_timeout: "30s"

database:
  type: "postgres"
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "dev_password"
  name: "golang_template_dev"
  sslmode: "disable"
  max_open_conns: 25
  max_idle_conns: 5
  max_lifetime: "300s"

jwt:
  secret: "your-super-secret-key-change-this-in-production"
  access_token_ttl: "15m"
  refresh_token_ttl: "24h"
  issuer: "golang-template"
  audience: "golang-template-users"
  algorithm: "HS256"

logger:
  level: "info"
  format: "json"
  output: "stdout"
  enable_caller: true
  enable_stacktrace: false
```

## 🌍 Environment-Specific Configuration

### Development (.env.development)

```bash
APP_SERVER_ENV=development
APP_SERVER_PORT=8080
APP_DATABASE_HOST=localhost
APP_DATABASE_NAME=golang_template_dev
APP_LOGGER_LEVEL=debug
APP_LOGGER_FORMAT=text
```

### Staging (.env.staging)

```bash
APP_SERVER_ENV=staging
APP_SERVER_PORT=8080
APP_DATABASE_HOST=staging-db.example.com
APP_DATABASE_NAME=golang_template_staging
APP_DATABASE_SSLMODE=require
APP_LOGGER_LEVEL=info
APP_LOGGER_FORMAT=json
```

### Production (.env.production)

```bash
APP_SERVER_ENV=production
APP_SERVER_PORT=8080
APP_DATABASE_HOST=prod-db.example.com
APP_DATABASE_NAME=golang_template_prod
APP_DATABASE_SSLMODE=require
APP_JWT_SECRET=super-secure-production-secret-key-64-characters-long
APP_LOGGER_LEVEL=warn
APP_LOGGER_FORMAT=json
```

## 🔧 Configuration in Code

### Accessing Configuration

```go
package main

import (
    "github.com/yantology/golang_template/internal/config"
)

func main() {
    // Initialize Viper
    if err := config.InitViper(); err != nil {
        panic(err)
    }
    
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }
    
    // Use configuration
    fmt.Printf("Server will run on port: %s\n", cfg.Server.Port)
}
```

### Configuration Struct

The configuration is mapped to a Go struct:

```go
type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Logger   LoggerConfig   `json:"logger"`
    JWT      JWTConfig      `json:"jwt"`
}

type ServerConfig struct {
    Port            string        `json:"port"`
    Host            string        `json:"host"`
    Env             string        `json:"env"`
    ReadTimeout     time.Duration `json:"read_timeout"`
    WriteTimeout    time.Duration `json:"write_timeout"`
    IdleTimeout     time.Duration `json:"idle_timeout"`
    ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}
```

## 🔒 Security Best Practices

### 1. Environment Variables

- **Never commit `.env` files** to version control
- **Use strong secrets** for JWT and database passwords
- **Rotate secrets regularly** in production
- **Use different secrets** for each environment

### 2. Database Configuration

- **Use SSL/TLS** in production (`APP_DATABASE_SSLMODE=require`)
- **Limit connection pools** to prevent resource exhaustion
- **Use dedicated database users** with minimal privileges

### 3. JWT Configuration

- **Use strong secrets** (minimum 32 characters)
- **Set appropriate token lifetimes** (short for access, longer for refresh)
- **Rotate JWT secrets** periodically
- **Use environment-specific secrets**

## 🔍 Configuration Validation

The template includes configuration validation:

```go
func (c *Config) Validate() error {
    if c.Server.Port == "" {
        return errors.New("server port is required")
    }
    
    if c.JWT.Secret == "" || len(c.JWT.Secret) < 32 {
        return errors.New("JWT secret must be at least 32 characters")
    }
    
    // Additional validation...
    return nil
}
```

## 📋 Environment File Templates

### Minimal Development Setup

```bash
# .env.minimal
APP_SERVER_PORT=8080
APP_DATABASE_PASSWORD=dev_password
APP_JWT_SECRET=development-secret-key-minimum-32-characters-long
```

### Complete Development Setup

```bash
# .env.complete
# Server Configuration
APP_SERVER_PORT=8080
APP_SERVER_HOST=0.0.0.0
APP_SERVER_ENV=development

# Database Configuration
APP_DATABASE_TYPE=postgres
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_USER=postgres
APP_DATABASE_PASSWORD=dev_password
APP_DATABASE_NAME=golang_template_dev
APP_DATABASE_SSLMODE=disable

# JWT Configuration
APP_JWT_SECRET=development-secret-key-minimum-32-characters-long
APP_JWT_ACCESS_TOKEN_TTL=15m
APP_JWT_REFRESH_TOKEN_TTL=24h

# Logger Configuration
APP_LOGGER_LEVEL=debug
APP_LOGGER_FORMAT=text
```

## 🔄 Dynamic Configuration Reload

For production environments, you may want to reload configuration without restarting:

```go
// Watch for configuration changes
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    fmt.Println("Config file changed:", e.Name)
    // Reload configuration logic here
})
```

## 🧪 Testing Configuration

### Test Environment Variables

```bash
# .env.test
APP_SERVER_ENV=test
APP_DATABASE_NAME=golang_template_test
APP_LOGGER_LEVEL=warn
```

### Configuration for Tests

```go
func TestConfig() *Config {
    return &Config{
        Server: ServerConfig{
            Port: "0", // Random available port
            Env:  "test",
        },
        Database: DatabaseConfig{
            Name: "test_db",
        },
        // ... other test configurations
    }
}
```

## 🚀 Next Steps

- **Learn about adding new configuration options**: [Adding New Configs](../configuration/adding-new-configs.md)
- **Understand configuration patterns**: [Configuration Patterns](../configuration/configuration-patterns.md)
- **Set up development workflow**: [Development Guide](./development.md)

---

**Need help?** Check the [troubleshooting section](./development.md#troubleshooting) or refer to the [configuration examples](../examples/).
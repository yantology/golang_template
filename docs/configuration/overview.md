# Configuration Overview

The Go Backend Template uses **Viper** for comprehensive configuration management, supporting multiple configuration sources and environment-specific settings.

## 🎯 Configuration Philosophy

### Design Principles

1. **Environment-Aware**: Different configurations for development, staging, and production
2. **Override Hierarchy**: Environment variables > Configuration files > Defaults
3. **Type Safety**: Strongly-typed configuration structures
4. **Validation**: Configuration validation at startup
5. **Documentation**: Self-documenting configuration with clear examples

### Configuration Sources (Priority Order)

1. **Environment Variables** (highest priority)
2. **Configuration Files** (YAML/JSON)
3. **Default Values** (lowest priority)

## 📂 Configuration Structure

```
internal/config/
├── config.go          # Configuration structures
├── database.go        # Database configuration
├── jwt.go             # JWT configuration  
├── logger.go          # Logger configuration
├── server.go          # Server configuration
└── viper.go           # Viper initialization and defaults
```

## 🔧 Configuration Components

### Main Configuration Structure

```go
// internal/config/config.go
package config

type Config struct {
    Server   ServerConfig   `json:"server" yaml:"server"`
    Database DatabaseConfig `json:"database" yaml:"database"`
    Logger   LoggerConfig   `json:"logger" yaml:"logger"`
    JWT      JWTConfig      `json:"jwt" yaml:"jwt"`
}

func Load() (*Config, error) {
    var cfg Config
    
    // Unmarshal configuration from Viper
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    // Validate configuration
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("configuration validation failed: %w", err)
    }
    
    return &cfg, nil
}

func (c *Config) Validate() error {
    if err := c.Server.Validate(); err != nil {
        return fmt.Errorf("server config: %w", err)
    }
    
    if err := c.Database.Validate(); err != nil {
        return fmt.Errorf("database config: %w", err)
    }
    
    if err := c.JWT.Validate(); err != nil {
        return fmt.Errorf("jwt config: %w", err)
    }
    
    if err := c.Logger.Validate(); err != nil {
        return fmt.Errorf("logger config: %w", err)
    }
    
    return nil
}
```

### Server Configuration

```go
// internal/config/server.go
package config

import (
    "fmt"
    "strconv"
    "time"
)

type ServerConfig struct {
    Port            string        `json:"port" yaml:"port"`
    Host            string        `json:"host" yaml:"host"`
    Env             string        `json:"env" yaml:"env"`
    ReadTimeout     time.Duration `json:"read_timeout" yaml:"read_timeout"`
    WriteTimeout    time.Duration `json:"write_timeout" yaml:"write_timeout"`
    IdleTimeout     time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
    ShutdownTimeout time.Duration `json:"shutdown_timeout" yaml:"shutdown_timeout"`
    EnableCORS      bool          `json:"enable_cors" yaml:"enable_cors"`
    CORSOrigins     []string      `json:"cors_origins" yaml:"cors_origins"`
}

func (c *ServerConfig) Validate() error {
    // Validate port
    if c.Port == "" {
        return fmt.Errorf("port is required")
    }
    
    if port, err := strconv.Atoi(c.Port); err != nil || port < 1 || port > 65535 {
        return fmt.Errorf("invalid port: %s", c.Port)
    }
    
    // Validate host
    if c.Host == "" {
        return fmt.Errorf("host is required")
    }
    
    // Validate environment
    validEnvs := map[string]bool{
        "development": true,
        "staging":     true,
        "production":  true,
        "test":        true,
    }
    if !validEnvs[c.Env] {
        return fmt.Errorf("invalid environment: %s", c.Env)
    }
    
    // Validate timeouts
    if c.ReadTimeout <= 0 {
        return fmt.Errorf("read_timeout must be positive")
    }
    if c.WriteTimeout <= 0 {
        return fmt.Errorf("write_timeout must be positive")
    }
    if c.IdleTimeout <= 0 {
        return fmt.Errorf("idle_timeout must be positive")
    }
    
    return nil
}

func (c *ServerConfig) Address() string {
    return c.Host + ":" + c.Port
}

func (c *ServerConfig) IsDevelopment() bool {
    return c.Env == "development"
}

func (c *ServerConfig) IsProduction() bool {
    return c.Env == "production"
}
```

### Database Configuration

```go
// internal/config/database.go
package config

import (
    "fmt"
    "strconv"
    "time"
)

type DatabaseConfig struct {
    Type         string        `json:"type" yaml:"type"`
    Host         string        `json:"host" yaml:"host"`
    Port         string        `json:"port" yaml:"port"`
    User         string        `json:"user" yaml:"user"`
    Password     string        `json:"password" yaml:"password"`
    Name         string        `json:"name" yaml:"name"`
    SSLMode      string        `json:"sslmode" yaml:"sslmode"`
    MaxOpenConns int           `json:"max_open_conns" yaml:"max_open_conns"`
    MaxIdleConns int           `json:"max_idle_conns" yaml:"max_idle_conns"`
    MaxLifetime  time.Duration `json:"max_lifetime" yaml:"max_lifetime"`
    MigrationPath string       `json:"migration_path" yaml:"migration_path"`
}

func (c *DatabaseConfig) Validate() error {
    // Only PostgreSQL is supported
    if c.Type != "postgres" {
        return fmt.Errorf("unsupported database type: %s (only postgres is supported)", c.Type)
    }
    
    // Validate required fields for PostgreSQL
    if c.Host == "" {
        return fmt.Errorf("host is required")
        }
        if c.Port == "" {
            return fmt.Errorf("port is required for %s", c.Type)
        }
        if c.User == "" {
            return fmt.Errorf("user is required for %s", c.Type)
        }
        if c.Password == "" {
            return fmt.Errorf("password is required for %s", c.Type)
        }
        
        // Validate port number
        if port, err := strconv.Atoi(c.Port); err != nil || port < 1 || port > 65535 {
            return fmt.Errorf("invalid port: %s", c.Port)
        }
    }
    
    if c.Name == "" {
        return fmt.Errorf("database name is required")
    }
    
    // Validate connection pool settings
    if c.MaxOpenConns < 1 {
        return fmt.Errorf("max_open_conns must be at least 1")
    }
    if c.MaxIdleConns < 1 {
        return fmt.Errorf("max_idle_conns must be at least 1")
    }
    if c.MaxIdleConns > c.MaxOpenConns {
        return fmt.Errorf("max_idle_conns cannot be greater than max_open_conns")
    }
    
    return nil
}

func (c *DatabaseConfig) PostgreSQLDSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
    )
}

func (c *DatabaseConfig) DSN() string {
    return c.PostgreSQLDSN()
}
```

### JWT Configuration

```go
// internal/config/jwt.go
package config

import (
    "fmt"
    "time"
)

type JWTConfig struct {
    Secret           string        `json:"secret" yaml:"secret"`
    AccessTokenTTL   time.Duration `json:"access_token_ttl" yaml:"access_token_ttl"`
    RefreshTokenTTL  time.Duration `json:"refresh_token_ttl" yaml:"refresh_token_ttl"`
    Issuer           string        `json:"issuer" yaml:"issuer"`
    Audience         string        `json:"audience" yaml:"audience"`
    Algorithm        string        `json:"algorithm" yaml:"algorithm"`
}

func (c *JWTConfig) Validate() error {
    if c.Secret == "" {
        return fmt.Errorf("jwt secret is required")
    }
    
    if len(c.Secret) < 32 {
        return fmt.Errorf("jwt secret must be at least 32 characters long")
    }
    
    if c.AccessTokenTTL <= 0 {
        return fmt.Errorf("access_token_ttl must be positive")
    }
    
    if c.RefreshTokenTTL <= 0 {
        return fmt.Errorf("refresh_token_ttl must be positive")
    }
    
    if c.RefreshTokenTTL <= c.AccessTokenTTL {
        return fmt.Errorf("refresh_token_ttl must be greater than access_token_ttl")
    }
    
    validAlgorithms := map[string]bool{
        "HS256": true,
        "HS384": true,
        "HS512": true,
    }
    if !validAlgorithms[c.Algorithm] {
        return fmt.Errorf("unsupported jwt algorithm: %s", c.Algorithm)
    }
    
    return nil
}
```

### Logger Configuration

```go
// internal/config/logger.go
package config

import "fmt"

type LoggerConfig struct {
    Level            string `json:"level" yaml:"level"`
    Format           string `json:"format" yaml:"format"`
    Output           string `json:"output" yaml:"output"`
    EnableCaller     bool   `json:"enable_caller" yaml:"enable_caller"`
    EnableStacktrace bool   `json:"enable_stacktrace" yaml:"enable_stacktrace"`
}

func (c *LoggerConfig) Validate() error {
    validLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
        "fatal": true,
    }
    if !validLevels[c.Level] {
        return fmt.Errorf("invalid log level: %s", c.Level)
    }
    
    validFormats := map[string]bool{
        "json": true,
        "text": true,
    }
    if !validFormats[c.Format] {
        return fmt.Errorf("invalid log format: %s", c.Format)
    }
    
    validOutputs := map[string]bool{
        "stdout": true,
        "stderr": true,
        "file":   true,
    }
    if !validOutputs[c.Output] {
        return fmt.Errorf("invalid log output: %s", c.Output)
    }
    
    return nil
}
```

## 🌍 Environment-Specific Configuration

### Development Configuration

```yaml
# config/development.yaml
server:
  port: "8080"
  host: "localhost"
  env: "development"
  enable_cors: true
  cors_origins: ["*"]

database:
  type: "postgres"
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "dev_password"
  name: "golang_template_dev"
  sslmode: "disable"

logger:
  level: "debug"
  format: "text"
  output: "stdout"

jwt:
  secret: "development-secret-key-minimum-32-characters-long"
  access_token_ttl: "15m"
  refresh_token_ttl: "24h"
```

### Production Configuration

```yaml
# config/production.yaml
server:
  port: "8080"
  host: "0.0.0.0"
  env: "production"
  enable_cors: true
  cors_origins: ["https://yourdomain.com"]

database:
  type: "postgres"
  host: "${DB_HOST}"
  port: "5432"
  user: "${DB_USER}"
  password: "${DB_PASSWORD}"
  name: "${DB_NAME}"
  sslmode: "require"
  max_open_conns: 50
  max_idle_conns: 10

logger:
  level: "info"
  format: "json"
  output: "stdout"

jwt:
  secret: "${JWT_SECRET}"
  access_token_ttl: "15m"
  refresh_token_ttl: "7d"
```

## 🔐 Environment Variables

### Complete Environment Variable Reference

| Category | Variable | Description | Default |
|----------|----------|-------------|---------|
| **Server** | `APP_SERVER_PORT` | HTTP server port | `8080` |
| | `APP_SERVER_HOST` | Server bind address | `0.0.0.0` |
| | `APP_SERVER_ENV` | Environment (development/staging/production) | `development` |
| | `APP_SERVER_READ_TIMEOUT` | HTTP read timeout | `10s` |
| | `APP_SERVER_WRITE_TIMEOUT` | HTTP write timeout | `10s` |
| | `APP_SERVER_IDLE_TIMEOUT` | HTTP idle timeout | `60s` |
| | `APP_SERVER_SHUTDOWN_TIMEOUT` | Graceful shutdown timeout | `30s` |
| **Database** | `APP_DATABASE_TYPE` | Database type (postgres only) | `postgres` |
| | `APP_DATABASE_HOST` | Database host | `localhost` |
| | `APP_DATABASE_PORT` | Database port | `5432` |
| | `APP_DATABASE_USER` | Database username | `postgres` |
| | `APP_DATABASE_PASSWORD` | Database password | `postgres` |
| | `APP_DATABASE_NAME` | Database name | `golang_template` |
| | `APP_DATABASE_SSLMODE` | PostgreSQL SSL mode | `disable` |
| | `APP_DATABASE_MAX_OPEN_CONNS` | Maximum open connections | `25` |
| | `APP_DATABASE_MAX_IDLE_CONNS` | Maximum idle connections | `5` |
| | `APP_DATABASE_MAX_LIFETIME` | Connection max lifetime | `300s` |
| **JWT** | `APP_JWT_SECRET` | JWT signing secret (min 32 chars) | `your-super-secret-key...` |
| | `APP_JWT_ACCESS_TOKEN_TTL` | Access token lifetime | `15m` |
| | `APP_JWT_REFRESH_TOKEN_TTL` | Refresh token lifetime | `24h` |
| | `APP_JWT_ISSUER` | JWT issuer | `golang-template` |
| | `APP_JWT_AUDIENCE` | JWT audience | `golang-template-users` |
| | `APP_JWT_ALGORITHM` | JWT signing algorithm | `HS256` |
| **Logger** | `APP_LOGGER_LEVEL` | Log level (debug/info/warn/error) | `info` |
| | `APP_LOGGER_FORMAT` | Log format (json/text) | `json` |
| | `APP_LOGGER_OUTPUT` | Log output (stdout/stderr/file) | `stdout` |
| | `APP_LOGGER_ENABLE_CALLER` | Include caller info in logs | `true` |
| | `APP_LOGGER_ENABLE_STACKTRACE` | Include stack trace in error logs | `false` |

## 🔧 Configuration Loading

### Viper Setup

The current `internal/config/viper.go` provides the foundation:

```go
package config

import (
    "fmt"
    "strings"

    "github.com/spf13/viper"
)

// InitViper initializes Viper configuration management
func InitViper() error {
    // Set config file name and paths
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")
    viper.AddConfigPath("/etc/app")

    // Environment variable configuration
    viper.SetEnvPrefix("APP")
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // Set default values
    setDefaults()

    // Read config file
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // Config file not found; ignore error since we have defaults and env vars
            fmt.Println("No config file found. Using environment variables and defaults.")
        } else {
            // Config file was found but another error was produced
            return fmt.Errorf("error reading config file: %w", err)
        }
    } else {
        fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
    }

    return nil
}

// setDefaults sets default values for all configuration options
func setDefaults() {
    // Server defaults
    viper.SetDefault("server.port", "8080")
    viper.SetDefault("server.host", "0.0.0.0")
    viper.SetDefault("server.env", "development")
    viper.SetDefault("server.read_timeout", "10s")
    viper.SetDefault("server.write_timeout", "10s")
    viper.SetDefault("server.idle_timeout", "60s")
    viper.SetDefault("server.shutdown_timeout", "30s")
    viper.SetDefault("server.enable_cors", true)
    viper.SetDefault("server.cors_origins", []string{"*"})

    // Database defaults
    viper.SetDefault("database.type", "postgres")
    viper.SetDefault("database.host", "localhost")
    viper.SetDefault("database.port", "5432")
    viper.SetDefault("database.user", "postgres")
    viper.SetDefault("database.password", "postgres")
    viper.SetDefault("database.name", "golang_template")
    viper.SetDefault("database.sslmode", "disable")
    viper.SetDefault("database.max_open_conns", 25)
    viper.SetDefault("database.max_idle_conns", 5)
    viper.SetDefault("database.max_lifetime", "300s")
    viper.SetDefault("database.migration_path", "./internal/data/migrations")

    // JWT defaults
    viper.SetDefault("jwt.secret", "your-super-secret-key-change-this-in-production")
    viper.SetDefault("jwt.access_token_ttl", "15m")
    viper.SetDefault("jwt.refresh_token_ttl", "24h")
    viper.SetDefault("jwt.issuer", "golang-template")
    viper.SetDefault("jwt.audience", "golang-template-users")
    viper.SetDefault("jwt.algorithm", "HS256")

    // Logger defaults
    viper.SetDefault("logger.level", "info")
    viper.SetDefault("logger.format", "json")
    viper.SetDefault("logger.output", "stdout")
    viper.SetDefault("logger.enable_caller", true)
    viper.SetDefault("logger.enable_stacktrace", false)
}
```

## 🛠️ Configuration Usage

### In Application Code

```go
// cmd/api/main.go
func main() {
    // Initialize Viper
    if err := config.InitViper(); err != nil {
        log.Fatalf("Failed to initialize configuration: %v", err)
    }
    
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Use configuration
    fmt.Printf("Starting server on %s\n", cfg.Server.Address())
    
    // Pass configuration to components
    db, err := database.Connect(cfg.Database)
    if err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }
}
```

### Configuration Validation Example

```go
func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Configuration error: %v", err)
    }
    
    // Configuration is guaranteed to be valid at this point
    log.Printf("Server running in %s mode", cfg.Server.Env)
}
```

## 🚀 Next Steps

- **Learn to add new configurations**: [Adding New Configs](./adding-new-configs.md)
- **Understand configuration patterns**: [Configuration Patterns](./configuration-patterns.md)
- **See environment variable reference**: [Environment Variables](./environment-variables.md)

---

This configuration system provides flexibility, type safety, and environment-specific settings while maintaining clear defaults and validation.
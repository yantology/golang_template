# Adding New Configuration Options

This guide shows you how to add new configuration options to the Go Backend Template while maintaining type safety, validation, and environment variable support.

## 🎯 Step-by-Step Process

### Step 1: Define Configuration Structure

First, create a new configuration struct or extend an existing one.

#### Adding a New Configuration Section

```go
// internal/config/redis.go
package config

import (
    "fmt"
    "strconv"
    "time"
)

type RedisConfig struct {
    Host            string        `json:"host" yaml:"host"`
    Port            string        `json:"port" yaml:"port"`
    Password        string        `json:"password" yaml:"password"`
    Database        int           `json:"database" yaml:"database"`
    MaxRetries      int           `json:"max_retries" yaml:"max_retries"`
    DialTimeout     time.Duration `json:"dial_timeout" yaml:"dial_timeout"`
    ReadTimeout     time.Duration `json:"read_timeout" yaml:"read_timeout"`
    WriteTimeout    time.Duration `json:"write_timeout" yaml:"write_timeout"`
    PoolSize        int           `json:"pool_size" yaml:"pool_size"`
    MinIdleConns    int           `json:"min_idle_conns" yaml:"min_idle_conns"`
    MaxConnAge      time.Duration `json:"max_conn_age" yaml:"max_conn_age"`
    PoolTimeout     time.Duration `json:"pool_timeout" yaml:"pool_timeout"`
    IdleTimeout     time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
    IdleCheckFreq   time.Duration `json:"idle_check_freq" yaml:"idle_check_freq"`
}

func (c *RedisConfig) Validate() error {
    if c.Host == "" {
        return fmt.Errorf("redis host is required")
    }
    
    if c.Port == "" {
        return fmt.Errorf("redis port is required")
    }
    
    // Validate port number
    if port, err := strconv.Atoi(c.Port); err != nil || port < 1 || port > 65535 {
        return fmt.Errorf("invalid redis port: %s", c.Port)
    }
    
    if c.Database < 0 || c.Database > 15 {
        return fmt.Errorf("redis database must be between 0 and 15")
    }
    
    if c.MaxRetries < 0 {
        return fmt.Errorf("redis max_retries cannot be negative")
    }
    
    if c.PoolSize < 1 {
        return fmt.Errorf("redis pool_size must be at least 1")
    }
    
    if c.MinIdleConns < 0 {
        return fmt.Errorf("redis min_idle_conns cannot be negative")
    }
    
    if c.MinIdleConns > c.PoolSize {
        return fmt.Errorf("redis min_idle_conns cannot be greater than pool_size")
    }
    
    return nil
}

func (c *RedisConfig) Address() string {
    return c.Host + ":" + c.Port
}

func (c *RedisConfig) IsLocal() bool {
    return c.Host == "localhost" || c.Host == "127.0.0.1"
}
```

#### Adding to Existing Configuration Section

```go
// internal/config/server.go - Adding new fields to existing struct
type ServerConfig struct {
    // Existing fields...
    Port            string        `json:"port" yaml:"port"`
    Host            string        `json:"host" yaml:"host"`
    
    // New fields
    MaxRequestSize  int64         `json:"max_request_size" yaml:"max_request_size"`
    EnableProfiling bool          `json:"enable_profiling" yaml:"enable_profiling"`
    ProfilingPath   string        `json:"profiling_path" yaml:"profiling_path"`
    TrustedProxies  []string      `json:"trusted_proxies" yaml:"trusted_proxies"`
    EnableMetrics   bool          `json:"enable_metrics" yaml:"enable_metrics"`
    MetricsPath     string        `json:"metrics_path" yaml:"metrics_path"`
}

// Update validation to include new fields
func (c *ServerConfig) Validate() error {
    // Existing validation...
    
    // New validation
    if c.MaxRequestSize <= 0 {
        return fmt.Errorf("max_request_size must be positive")
    }
    
    if c.EnableProfiling && c.ProfilingPath == "" {
        return fmt.Errorf("profiling_path is required when enable_profiling is true")
    }
    
    if c.EnableMetrics && c.MetricsPath == "" {
        return fmt.Errorf("metrics_path is required when enable_metrics is true")
    }
    
    return nil
}
```

### Step 2: Update Main Configuration

Add the new configuration to the main config struct:

```go
// internal/config/config.go
type Config struct {
    Server   ServerConfig   `json:"server" yaml:"server"`
    Database DatabaseConfig `json:"database" yaml:"database"`
    Logger   LoggerConfig   `json:"logger" yaml:"logger"`
    JWT      JWTConfig      `json:"jwt" yaml:"jwt"`
    Redis    RedisConfig    `json:"redis" yaml:"redis"`      // New configuration
    Email    EmailConfig    `json:"email" yaml:"email"`      // Another new one
    Storage  StorageConfig  `json:"storage" yaml:"storage"`  // And another
}

// Update validation to include new configurations
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
    
    // New validations
    if err := c.Redis.Validate(); err != nil {
        return fmt.Errorf("redis config: %w", err)
    }
    
    if err := c.Email.Validate(); err != nil {
        return fmt.Errorf("email config: %w", err)
    }
    
    if err := c.Storage.Validate(); err != nil {
        return fmt.Errorf("storage config: %w", err)
    }
    
    return nil
}
```

### Step 3: Add Default Values

Update the Viper defaults in `internal/config/viper.go`:

```go
// internal/config/viper.go - Add to setDefaults() function
func setDefaults() {
    // Existing defaults...
    
    // Redis defaults
    viper.SetDefault("redis.host", "localhost")
    viper.SetDefault("redis.port", "6379")
    viper.SetDefault("redis.password", "")
    viper.SetDefault("redis.database", 0)
    viper.SetDefault("redis.max_retries", 3)
    viper.SetDefault("redis.dial_timeout", "5s")
    viper.SetDefault("redis.read_timeout", "3s")
    viper.SetDefault("redis.write_timeout", "3s")
    viper.SetDefault("redis.pool_size", 10)
    viper.SetDefault("redis.min_idle_conns", 5)
    viper.SetDefault("redis.max_conn_age", "30m")
    viper.SetDefault("redis.pool_timeout", "4s")
    viper.SetDefault("redis.idle_timeout", "5m")
    viper.SetDefault("redis.idle_check_freq", "1m")
    
    // Server defaults (new fields)
    viper.SetDefault("server.max_request_size", 10485760) // 10MB
    viper.SetDefault("server.enable_profiling", false)
    viper.SetDefault("server.profiling_path", "/debug/pprof")
    viper.SetDefault("server.trusted_proxies", []string{})
    viper.SetDefault("server.enable_metrics", false)
    viper.SetDefault("server.metrics_path", "/metrics")
    
    // Email defaults
    viper.SetDefault("email.smtp_host", "localhost")
    viper.SetDefault("email.smtp_port", "587")
    viper.SetDefault("email.username", "")
    viper.SetDefault("email.password", "")
    viper.SetDefault("email.from_address", "noreply@example.com")
    viper.SetDefault("email.from_name", "Go Template")
    viper.SetDefault("email.use_tls", true)
    
    // Storage defaults
    viper.SetDefault("storage.type", "local")
    viper.SetDefault("storage.local_path", "./uploads")
    viper.SetDefault("storage.max_file_size", 52428800) // 50MB
    viper.SetDefault("storage.allowed_extensions", []string{".jpg", ".jpeg", ".png", ".gif", ".pdf"})
}
```

### Step 4: Update Environment Variables

Add the new environment variables to `.env.example`:

```bash
# .env.example - Add new sections

# Redis Configuration
APP_REDIS_HOST=localhost
APP_REDIS_PORT=6379
APP_REDIS_PASSWORD=
APP_REDIS_DATABASE=0
APP_REDIS_MAX_RETRIES=3
APP_REDIS_POOL_SIZE=10

# Server Configuration (new options)
APP_SERVER_MAX_REQUEST_SIZE=10485760
APP_SERVER_ENABLE_PROFILING=false
APP_SERVER_PROFILING_PATH=/debug/pprof
APP_SERVER_ENABLE_METRICS=false
APP_SERVER_METRICS_PATH=/metrics

# Email Configuration
APP_EMAIL_SMTP_HOST=smtp.gmail.com
APP_EMAIL_SMTP_PORT=587
APP_EMAIL_USERNAME=your-email@gmail.com
APP_EMAIL_PASSWORD=your-app-password
APP_EMAIL_FROM_ADDRESS=noreply@yourdomain.com
APP_EMAIL_FROM_NAME=Your App Name
APP_EMAIL_USE_TLS=true

# Storage Configuration
APP_STORAGE_TYPE=local
APP_STORAGE_LOCAL_PATH=./uploads
APP_STORAGE_MAX_FILE_SIZE=52428800
```

### Step 5: Create Additional Configuration Files

Add more specific configuration structures as needed:

```go
// internal/config/email.go
package config

import "fmt"

type EmailConfig struct {
    SMTPHost     string `json:"smtp_host" yaml:"smtp_host"`
    SMTPPort     string `json:"smtp_port" yaml:"smtp_port"`
    Username     string `json:"username" yaml:"username"`
    Password     string `json:"password" yaml:"password"`
    FromAddress  string `json:"from_address" yaml:"from_address"`
    FromName     string `json:"from_name" yaml:"from_name"`
    UseTLS       bool   `json:"use_tls" yaml:"use_tls"`
}

func (c *EmailConfig) Validate() error {
    if c.SMTPHost == "" {
        return fmt.Errorf("smtp_host is required")
    }
    
    if c.SMTPPort == "" {
        return fmt.Errorf("smtp_port is required")
    }
    
    if c.FromAddress == "" {
        return fmt.Errorf("from_address is required")
    }
    
    if c.FromName == "" {
        return fmt.Errorf("from_name is required")
    }
    
    return nil
}

func (c *EmailConfig) SMTPAddress() string {
    return c.SMTPHost + ":" + c.SMTPPort
}
```

```go
// internal/config/storage.go
package config

import (
    "fmt"
    "strings"
)

type StorageConfig struct {
    Type              string   `json:"type" yaml:"type"`
    LocalPath         string   `json:"local_path" yaml:"local_path"`
    MaxFileSize       int64    `json:"max_file_size" yaml:"max_file_size"`
    AllowedExtensions []string `json:"allowed_extensions" yaml:"allowed_extensions"`
    
    // S3 configuration (optional)
    S3Bucket    string `json:"s3_bucket" yaml:"s3_bucket"`
    S3Region    string `json:"s3_region" yaml:"s3_region"`
    S3AccessKey string `json:"s3_access_key" yaml:"s3_access_key"`
    S3SecretKey string `json:"s3_secret_key" yaml:"s3_secret_key"`
}

func (c *StorageConfig) Validate() error {
    validTypes := map[string]bool{
        "local": true,
        "s3":    true,
    }
    if !validTypes[c.Type] {
        return fmt.Errorf("unsupported storage type: %s", c.Type)
    }
    
    if c.Type == "local" && c.LocalPath == "" {
        return fmt.Errorf("local_path is required for local storage")
    }
    
    if c.Type == "s3" {
        if c.S3Bucket == "" {
            return fmt.Errorf("s3_bucket is required for S3 storage")
        }
        if c.S3Region == "" {
            return fmt.Errorf("s3_region is required for S3 storage")
        }
        if c.S3AccessKey == "" {
            return fmt.Errorf("s3_access_key is required for S3 storage")
        }
        if c.S3SecretKey == "" {
            return fmt.Errorf("s3_secret_key is required for S3 storage")
        }
    }
    
    if c.MaxFileSize <= 0 {
        return fmt.Errorf("max_file_size must be positive")
    }
    
    return nil
}

func (c *StorageConfig) IsLocal() bool {
    return c.Type == "local"
}

func (c *StorageConfig) IsS3() bool {
    return c.Type == "s3"
}

func (c *StorageConfig) IsExtensionAllowed(filename string) bool {
    if len(c.AllowedExtensions) == 0 {
        return true // No restrictions
    }
    
    // Extract extension
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return false // No extension
    }
    
    ext := "." + strings.ToLower(parts[len(parts)-1])
    
    for _, allowed := range c.AllowedExtensions {
        if strings.ToLower(allowed) == ext {
            return true
        }
    }
    
    return false
}
```

### Step 6: Update Documentation

Add the new configuration options to the documentation:

```go
// docs/configuration/environment-variables.md
// Add new sections for Redis, Email, Storage configurations
```

### Step 7: Use Configuration in Code

Create services that use the new configuration:

```go
// internal/pkg/redis/redis.go
package redis

import (
    "context"
    "time"
    
    "github.com/go-redis/redis/v8"
    "github.com/yantology/golang_template/internal/config"
)

func NewClient(cfg config.RedisConfig) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:         cfg.Address(),
        Password:     cfg.Password,
        DB:           cfg.Database,
        MaxRetries:   cfg.MaxRetries,
        DialTimeout:  cfg.DialTimeout,
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
        PoolSize:     cfg.PoolSize,
        MinIdleConns: cfg.MinIdleConns,
        MaxConnAge:   cfg.MaxConnAge,
        PoolTimeout:  cfg.PoolTimeout,
        IdleTimeout:  cfg.IdleTimeout,
        IdleCheckFrequency: cfg.IdleCheckFreq,
    })
}

func Connect(cfg config.RedisConfig) (*redis.Client, error) {
    client := NewClient(cfg)
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }
    
    return client, nil
}
```

```go
// internal/pkg/email/email.go
package email

import (
    "crypto/tls"
    "fmt"
    "net/smtp"
    
    "github.com/yantology/golang_template/internal/config"
)

type Service struct {
    config config.EmailConfig
    auth   smtp.Auth
}

func NewService(cfg config.EmailConfig) *Service {
    var auth smtp.Auth
    if cfg.Username != "" && cfg.Password != "" {
        auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)
    }
    
    return &Service{
        config: cfg,
        auth:   auth,
    }
}

func (s *Service) SendEmail(to []string, subject, body string) error {
    // Create message
    message := fmt.Sprintf("From: %s <%s>\r\n"+
        "To: %s\r\n"+
        "Subject: %s\r\n"+
        "\r\n"+
        "%s\r\n", s.config.FromName, s.config.FromAddress, 
        strings.Join(to, ","), subject, body)
    
    // Send email
    if s.config.UseTLS {
        return s.sendWithTLS(to, []byte(message))
    }
    
    return smtp.SendMail(s.config.SMTPAddress(), s.auth, s.config.FromAddress, to, []byte(message))
}

func (s *Service) sendWithTLS(to []string, message []byte) error {
    // Implementation for TLS email sending
    // ... implementation details
}
```

## 🔧 Configuration Templates

### Complete Example: Adding Cache Configuration

Here's a complete example of adding cache configuration:

#### 1. Configuration Structure

```go
// internal/config/cache.go
package config

import (
    "fmt"
    "time"
)

type CacheConfig struct {
    Type            string        `json:"type" yaml:"type"`
    DefaultTTL      time.Duration `json:"default_ttl" yaml:"default_ttl"`
    CleanupInterval time.Duration `json:"cleanup_interval" yaml:"cleanup_interval"`
    
    // Redis cache (when type = "redis")
    Redis RedisConfig `json:"redis" yaml:"redis"`
    
    // Memory cache (when type = "memory")
    MaxSize     int           `json:"max_size" yaml:"max_size"`
    MaxItemSize int           `json:"max_item_size" yaml:"max_item_size"`
}

func (c *CacheConfig) Validate() error {
    validTypes := map[string]bool{
        "memory": true,
        "redis":  true,
        "none":   true,
    }
    if !validTypes[c.Type] {
        return fmt.Errorf("unsupported cache type: %s", c.Type)
    }
    
    if c.Type != "none" && c.DefaultTTL <= 0 {
        return fmt.Errorf("default_ttl must be positive")
    }
    
    if c.Type == "memory" {
        if c.MaxSize <= 0 {
            return fmt.Errorf("max_size must be positive for memory cache")
        }
        if c.MaxItemSize <= 0 {
            return fmt.Errorf("max_item_size must be positive for memory cache")
        }
        if c.CleanupInterval <= 0 {
            return fmt.Errorf("cleanup_interval must be positive for memory cache")
        }
    }
    
    if c.Type == "redis" {
        if err := c.Redis.Validate(); err != nil {
            return fmt.Errorf("redis cache config: %w", err)
        }
    }
    
    return nil
}

func (c *CacheConfig) IsEnabled() bool {
    return c.Type != "none"
}

func (c *CacheConfig) IsMemory() bool {
    return c.Type == "memory"
}

func (c *CacheConfig) IsRedis() bool {
    return c.Type == "redis"
}
```

#### 2. Add to Main Config

```go
// internal/config/config.go
type Config struct {
    // ... existing fields
    Cache CacheConfig `json:"cache" yaml:"cache"`
}

func (c *Config) Validate() error {
    // ... existing validations
    if err := c.Cache.Validate(); err != nil {
        return fmt.Errorf("cache config: %w", err)
    }
    return nil
}
```

#### 3. Add Defaults

```go
// internal/config/viper.go
func setDefaults() {
    // ... existing defaults
    
    // Cache defaults
    viper.SetDefault("cache.type", "memory")
    viper.SetDefault("cache.default_ttl", "1h")
    viper.SetDefault("cache.cleanup_interval", "10m")
    viper.SetDefault("cache.max_size", 1000)
    viper.SetDefault("cache.max_item_size", 1048576) // 1MB
    
    // Cache Redis defaults (nested)
    viper.SetDefault("cache.redis.host", "localhost")
    viper.SetDefault("cache.redis.port", "6379")
    viper.SetDefault("cache.redis.database", 1) // Different from main Redis
}
```

#### 4. Environment Variables

```bash
# Cache Configuration
APP_CACHE_TYPE=memory
APP_CACHE_DEFAULT_TTL=1h
APP_CACHE_CLEANUP_INTERVAL=10m
APP_CACHE_MAX_SIZE=1000
APP_CACHE_MAX_ITEM_SIZE=1048576

# Cache Redis Configuration (if using Redis cache)
APP_CACHE_REDIS_HOST=localhost
APP_CACHE_REDIS_PORT=6379
APP_CACHE_REDIS_DATABASE=1
```

#### 5. Use in Application

```go
// internal/pkg/cache/cache.go
package cache

import (
    "github.com/yantology/golang_template/internal/config"
)

func NewCache(cfg config.CacheConfig) (Cache, error) {
    if !cfg.IsEnabled() {
        return &NoOpCache{}, nil
    }
    
    switch cfg.Type {
    case "memory":
        return NewMemoryCache(cfg)
    case "redis":
        return NewRedisCache(cfg.Redis)
    default:
        return nil, fmt.Errorf("unsupported cache type: %s", cfg.Type)
    }
}
```

## 🧪 Testing Configuration

### Unit Tests

```go
// internal/config/config_test.go
package config

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
)

func TestRedisConfig_Validate(t *testing.T) {
    tests := []struct {
        name    string
        config  RedisConfig
        wantErr bool
    }{
        {
            name: "valid config",
            config: RedisConfig{
                Host:     "localhost",
                Port:     "6379",
                Database: 0,
                PoolSize: 10,
            },
            wantErr: false,
        },
        {
            name: "missing host",
            config: RedisConfig{
                Port:     "6379",
                Database: 0,
                PoolSize: 10,
            },
            wantErr: true,
        },
        {
            name: "invalid port",
            config: RedisConfig{
                Host:     "localhost",
                Port:     "invalid",
                Database: 0,
                PoolSize: 10,
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Integration Tests

```go
// tests/integration/config_test.go
package integration

import (
    "os"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/yantology/golang_template/internal/config"
)

func TestConfigurationFromEnv(t *testing.T) {
    // Set environment variables
    os.Setenv("APP_REDIS_HOST", "redis-test")
    os.Setenv("APP_REDIS_PORT", "6380")
    defer func() {
        os.Unsetenv("APP_REDIS_HOST")
        os.Unsetenv("APP_REDIS_PORT")
    }()
    
    // Initialize and load config
    err := config.InitViper()
    assert.NoError(t, err)
    
    cfg, err := config.Load()
    assert.NoError(t, err)
    
    // Verify environment variables are loaded
    assert.Equal(t, "redis-test", cfg.Redis.Host)
    assert.Equal(t, "6380", cfg.Redis.Port)
}
```

## 🎯 Best Practices

### 1. Configuration Validation
- Always validate configuration at startup
- Provide clear error messages
- Validate relationships between fields

### 2. Environment Variables
- Use consistent naming (APP_ prefix)
- Document all environment variables
- Provide sensible defaults

### 3. Type Safety
- Use appropriate types (time.Duration, not strings)
- Use enums for limited value sets
- Provide helper methods for common operations

### 4. Documentation
- Update .env.example
- Add to configuration documentation
- Include validation rules in comments

### 5. Testing
- Test configuration validation
- Test environment variable loading
- Test default values

## 🚀 Next Steps

- **Learn configuration patterns**: [Configuration Patterns](./configuration-patterns.md)
- **See environment variable reference**: [Environment Variables](./environment-variables.md)
- **Understand the complete setup**: [Getting Started](../getting-started/setup.md)

---

Adding new configuration options is straightforward when following these patterns, ensuring your application remains configurable and maintainable.
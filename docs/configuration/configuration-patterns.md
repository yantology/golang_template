# Configuration Patterns

This guide covers best practices and patterns for configuration management in the Go Backend Template.

## 🎯 Configuration Design Patterns

### 1. Hierarchical Configuration

Organize configuration in a logical hierarchy that reflects your application structure:

```go
type Config struct {
    Server    ServerConfig    `json:"server"`
    Database  DatabaseConfig  `json:"database"`
    External  ExternalConfig  `json:"external"`
    Features  FeatureConfig   `json:"features"`
}

type ExternalConfig struct {
    Redis    RedisConfig    `json:"redis"`
    Email    EmailConfig    `json:"email"`
    Storage  StorageConfig  `json:"storage"`
    Payment  PaymentConfig  `json:"payment"`
}

type FeatureConfig struct {
    EnableMetrics   bool `json:"enable_metrics"`
    EnableProfiling bool `json:"enable_profiling"`
    EnableCache     bool `json:"enable_cache"`
}
```

### 2. Environment-Aware Configuration

Create configuration that adapts to different environments:

```go
type Config struct {
    Environment string `json:"environment"`
    // ... other fields
}

func (c *Config) IsDevelopment() bool {
    return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
    return c.Environment == "production"
}

func (c *Config) IsTest() bool {
    return c.Environment == "test"
}

// Environment-specific behavior
func (c *Config) GetLogLevel() string {
    if c.IsDevelopment() {
        return "debug"
    }
    return "info"
}

func (c *Config) GetCORSOrigins() []string {
    if c.IsDevelopment() {
        return []string{"*"}
    }
    return c.Server.CORSOrigins
}
```

### 3. Configuration Builder Pattern

For complex configuration scenarios:

```go
type ConfigBuilder struct {
    config Config
}

func NewConfigBuilder() *ConfigBuilder {
    return &ConfigBuilder{
        config: Config{
            Server:   defaultServerConfig(),
            Database: defaultDatabaseConfig(),
            Logger:   defaultLoggerConfig(),
        },
    }
}

func (b *ConfigBuilder) WithEnvironment(env string) *ConfigBuilder {
    b.config.Environment = env
    
    // Apply environment-specific defaults
    switch env {
    case "development":
        b.config.Logger.Level = "debug"
        b.config.Logger.Format = "text"
        b.config.Server.EnableCORS = true
    case "production":
        b.config.Logger.Level = "info"
        b.config.Logger.Format = "json"
        b.config.Server.EnableCORS = false
    }
    
    return b
}

func (b *ConfigBuilder) WithDatabase(host, name string) *ConfigBuilder {
    b.config.Database.Host = host
    b.config.Database.Name = name
    return b
}

func (b *ConfigBuilder) WithFeatures(features FeatureConfig) *ConfigBuilder {
    b.config.Features = features
    return b
}

func (b *ConfigBuilder) Build() (*Config, error) {
    if err := b.config.Validate(); err != nil {
        return nil, err
    }
    return &b.config, nil
}

// Usage
config, err := NewConfigBuilder().
    WithEnvironment("production").
    WithDatabase("prod-db.example.com", "myapp_prod").
    WithFeatures(FeatureConfig{
        EnableMetrics:   true,
        EnableCache:     true,
    }).
    Build()
```

### 4. Configuration Profiles

Support multiple configuration profiles:

```go
type ConfigProfile struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Config      Config `json:"config"`
}

var profiles = map[string]ConfigProfile{
    "development": {
        Name:        "development",
        Description: "Local development environment",
        Config: Config{
            Server: ServerConfig{
                Port: "8080",
                Host: "localhost",
                Env:  "development",
            },
            Database: DatabaseConfig{
                Host:    "localhost",
                Name:    "myapp_dev",
                SSLMode: "disable",
            },
            Logger: LoggerConfig{
                Level:  "debug",
                Format: "text",
            },
        },
    },
    "docker": {
        Name:        "docker",
        Description: "Docker development environment",
        Config: Config{
            Server: ServerConfig{
                Port: "8080",
                Host: "0.0.0.0",
                Env:  "development",
            },
            Database: DatabaseConfig{
                Host:    "postgres",
                Name:    "myapp_dev",
                SSLMode: "disable",
            },
        },
    },
    "production": {
        Name:        "production",
        Description: "Production environment",
        Config: Config{
            Server: ServerConfig{
                Port: "8080",
                Host: "0.0.0.0",
                Env:  "production",
            },
            Database: DatabaseConfig{
                SSLMode:      "require",
                MaxOpenConns: 50,
                MaxIdleConns: 10,
            },
            Logger: LoggerConfig{
                Level:  "info",
                Format: "json",
            },
        },
    },
}

func LoadProfile(profileName string) (*Config, error) {
    profile, exists := profiles[profileName]
    if !exists {
        return nil, fmt.Errorf("profile %s not found", profileName)
    }
    
    // Start with profile config
    config := profile.Config
    
    // Override with environment variables
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## 🔧 Advanced Configuration Patterns

### 1. Dynamic Configuration

Configuration that can change at runtime:

```go
type DynamicConfig struct {
    mu     sync.RWMutex
    config Config
    
    // Callbacks for configuration changes
    changeCallbacks []func(old, new Config)
}

func NewDynamicConfig(initial Config) *DynamicConfig {
    return &DynamicConfig{
        config: initial,
    }
}

func (dc *DynamicConfig) Get() Config {
    dc.mu.RLock()
    defer dc.mu.RUnlock()
    return dc.config
}

func (dc *DynamicConfig) Update(newConfig Config) error {
    if err := newConfig.Validate(); err != nil {
        return err
    }
    
    dc.mu.Lock()
    oldConfig := dc.config
    dc.config = newConfig
    dc.mu.Unlock()
    
    // Notify callbacks
    for _, callback := range dc.changeCallbacks {
        go callback(oldConfig, newConfig)
    }
    
    return nil
}

func (dc *DynamicConfig) OnChange(callback func(old, new Config)) {
    dc.changeCallbacks = append(dc.changeCallbacks, callback)
}

// Usage
dynamicConfig := NewDynamicConfig(config)

// React to configuration changes
dynamicConfig.OnChange(func(old, new Config) {
    if old.Logger.Level != new.Logger.Level {
        logger.SetLevel(new.Logger.Level)
    }
})
```

### 2. Configuration Validation with Dependencies

Validate configuration based on relationships between fields:

```go
func (c *Config) Validate() error {
    // Basic field validation
    if err := c.Server.Validate(); err != nil {
        return fmt.Errorf("server: %w", err)
    }
    
    // Cross-field validation
    if c.Features.EnableCache && c.Cache.Type == "" {
        return fmt.Errorf("cache type must be specified when cache is enabled")
    }
    
    if c.Features.EnableMetrics && c.Server.MetricsPath == "" {
        return fmt.Errorf("metrics path must be specified when metrics are enabled")
    }
    
    // Environment-specific validation
    if c.IsProduction() {
        if err := c.validateProduction(); err != nil {
            return fmt.Errorf("production validation: %w", err)
        }
    }
    
    return nil
}

func (c *Config) validateProduction() error {
    if c.Database.SSLMode == "disable" {
        return fmt.Errorf("SSL must be enabled in production")
    }
    
    if len(c.JWT.Secret) < 64 {
        return fmt.Errorf("JWT secret must be at least 64 characters in production")
    }
    
    if c.Logger.Level == "debug" {
        return fmt.Errorf("debug logging should not be used in production")
    }
    
    return nil
}
```

### 3. Configuration Merging

Merge configuration from multiple sources:

```go
type ConfigSource interface {
    Load() (Config, error)
    Priority() int
}

type FileConfigSource struct {
    path string
}

func (f *FileConfigSource) Load() (Config, error) {
    // Load from file
}

func (f *FileConfigSource) Priority() int {
    return 1 // Lower priority
}

type EnvironmentConfigSource struct{}

func (e *EnvironmentConfigSource) Load() (Config, error) {
    // Load from environment variables
}

func (e *EnvironmentConfigSource) Priority() int {
    return 2 // Higher priority
}

type ConfigMerger struct {
    sources []ConfigSource
}

func (cm *ConfigMerger) AddSource(source ConfigSource) {
    cm.sources = append(cm.sources, source)
    
    // Sort by priority
    sort.Slice(cm.sources, func(i, j int) bool {
        return cm.sources[i].Priority() < cm.sources[j].Priority()
    })
}

func (cm *ConfigMerger) Merge() (Config, error) {
    var final Config
    
    for _, source := range cm.sources {
        config, err := source.Load()
        if err != nil {
            continue // Skip failed sources
        }
        
        // Merge configurations (higher priority overwrites)
        final = mergeConfigs(final, config)
    }
    
    return final, nil
}

func mergeConfigs(base, override Config) Config {
    // Use reflection or manual merging to combine configs
    // Override non-zero values from override into base
    return merged
}
```

### 4. Feature Flags Integration

Integrate feature flags with configuration:

```go
type FeatureFlag struct {
    Name        string `json:"name"`
    Enabled     bool   `json:"enabled"`
    Description string `json:"description"`
    Conditions  map[string]interface{} `json:"conditions,omitempty"`
}

type FeatureFlags struct {
    flags map[string]FeatureFlag
}

func NewFeatureFlags() *FeatureFlags {
    return &FeatureFlags{
        flags: make(map[string]FeatureFlag),
    }
}

func (ff *FeatureFlags) IsEnabled(flagName string) bool {
    flag, exists := ff.flags[flagName]
    if !exists {
        return false
    }
    
    return flag.Enabled
}

func (ff *FeatureFlags) EnableFeature(flagName string) {
    flag := ff.flags[flagName]
    flag.Enabled = true
    ff.flags[flagName] = flag
}

// Integration with configuration
type ConfigWithFeatures struct {
    Config
    Features *FeatureFlags `json:"features"`
}

func (c *ConfigWithFeatures) ShouldEnableCache() bool {
    return c.Features.IsEnabled("cache")
}

func (c *ConfigWithFeatures) ShouldEnableMetrics() bool {
    return c.Features.IsEnabled("metrics")
}
```

### 5. Configuration Templates

Use templates for dynamic configuration generation:

```go
type ConfigTemplate struct {
    template *template.Template
    data     map[string]interface{}
}

func NewConfigTemplate(templateStr string) (*ConfigTemplate, error) {
    tmpl, err := template.New("config").Parse(templateStr)
    if err != nil {
        return nil, err
    }
    
    return &ConfigTemplate{
        template: tmpl,
        data:     make(map[string]interface{}),
    }, nil
}

func (ct *ConfigTemplate) SetData(key string, value interface{}) {
    ct.data[key] = value
}

func (ct *ConfigTemplate) Generate() (string, error) {
    var buf bytes.Buffer
    if err := ct.template.Execute(&buf, ct.data); err != nil {
        return "", err
    }
    return buf.String(), nil
}

// Usage
configTemplate := `
server:
  port: "{{.Port}}"
  host: "{{.Host}}"
  env: "{{.Environment}}"

database:
  host: "{{.DatabaseHost}}"
  name: "{{.AppName}}_{{.Environment}}"
  
{{if eq .Environment "production"}}
  sslmode: "require"
  max_open_conns: 50
{{else}}
  sslmode: "disable"
  max_open_conns: 10
{{end}}
`

tmpl, _ := NewConfigTemplate(configTemplate)
tmpl.SetData("Port", "8080")
tmpl.SetData("Host", "0.0.0.0")
tmpl.SetData("Environment", "production")
tmpl.SetData("DatabaseHost", "prod-db.example.com")
tmpl.SetData("AppName", "myapp")

configYAML, _ := tmpl.Generate()
```

## 🛡️ Security Patterns

### 1. Secret Management

```go
type SecretManager interface {
    GetSecret(key string) (string, error)
}

type VaultSecretManager struct {
    client *vault.Client
}

func (v *VaultSecretManager) GetSecret(key string) (string, error) {
    // Fetch from HashiCorp Vault
}

type EnvSecretManager struct{}

func (e *EnvSecretManager) GetSecret(key string) (string, error) {
    value := os.Getenv(key)
    if value == "" {
        return "", fmt.Errorf("secret %s not found", key)
    }
    return value, nil
}

type ConfigWithSecrets struct {
    Config
    secretManager SecretManager
}

func (c *ConfigWithSecrets) GetDatabasePassword() (string, error) {
    return c.secretManager.GetSecret("DATABASE_PASSWORD")
}

func (c *ConfigWithSecrets) GetJWTSecret() (string, error) {
    return c.secretManager.GetSecret("JWT_SECRET")
}
```

### 2. Configuration Encryption

```go
type EncryptedConfig struct {
    data []byte
    key  []byte
}

func NewEncryptedConfig(plaintext string, key []byte) (*EncryptedConfig, error) {
    encrypted, err := encrypt([]byte(plaintext), key)
    if err != nil {
        return nil, err
    }
    
    return &EncryptedConfig{
        data: encrypted,
        key:  key,
    }, nil
}

func (ec *EncryptedConfig) Decrypt() (string, error) {
    decrypted, err := decrypt(ec.data, ec.key)
    if err != nil {
        return "", err
    }
    return string(decrypted), nil
}

func (ec *EncryptedConfig) Load() (Config, error) {
    configStr, err := ec.Decrypt()
    if err != nil {
        return Config{}, err
    }
    
    var config Config
    if err := yaml.Unmarshal([]byte(configStr), &config); err != nil {
        return Config{}, err
    }
    
    return config, nil
}
```

## 🧪 Testing Patterns

### 1. Configuration Mocking

```go
type MockConfig struct {
    Config
    overrides map[string]interface{}
}

func NewMockConfig(base Config) *MockConfig {
    return &MockConfig{
        Config:    base,
        overrides: make(map[string]interface{}),
    }
}

func (m *MockConfig) SetDatabaseHost(host string) {
    m.Config.Database.Host = host
}

func (m *MockConfig) SetLogLevel(level string) {
    m.Config.Logger.Level = level
}

func (m *MockConfig) SetTestMode() {
    m.Config.Server.Env = "test"
    m.Config.Logger.Level = "error"
    m.Config.Database.Name = "test_db"
}

// Test helper
func TestConfig() *MockConfig {
    base := Config{
        Server: ServerConfig{
            Port: "0", // Random port for tests
            Host: "localhost",
            Env:  "test",
        },
        Database: DatabaseConfig{
            Host: "localhost",
            Name: "test_db",
        },
        Logger: LoggerConfig{
            Level:  "error",
            Format: "text",
        },
    }
    
    return NewMockConfig(base)
}
```

### 2. Configuration Validation Testing

```go
func TestConfigValidation(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid config",
            config: Config{
                Server:   validServerConfig(),
                Database: validDatabaseConfig(),
                JWT:      validJWTConfig(),
            },
            wantErr: false,
        },
        {
            name: "invalid JWT secret",
            config: Config{
                Server:   validServerConfig(),
                Database: validDatabaseConfig(),
                JWT:      JWTConfig{Secret: "short"},
            },
            wantErr: true,
            errMsg:  "jwt secret must be at least 32 characters",
        },
        {
            name: "production without SSL",
            config: Config{
                Server: ServerConfig{Env: "production"},
                Database: DatabaseConfig{
                    Host:    "localhost",
                    SSLMode: "disable",
                },
                JWT: validJWTConfig(),
            },
            wantErr: true,
            errMsg:  "SSL must be enabled in production",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errMsg != "" {
                    assert.Contains(t, err.Error(), tt.errMsg)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 🎯 Best Practices Summary

### 1. Structure
- Use hierarchical configuration structures
- Group related settings together
- Use meaningful names and types

### 2. Validation
- Validate configuration at startup
- Provide clear error messages
- Include cross-field validation

### 3. Environment Awareness
- Support multiple environments
- Use environment-specific defaults
- Validate environment-specific requirements

### 4. Security
- Never commit secrets to version control
- Use secret management systems
- Encrypt sensitive configuration

### 5. Testing
- Create test-specific configurations
- Mock configuration for unit tests
- Test configuration validation thoroughly

### 6. Documentation
- Document all configuration options
- Provide examples for different environments
- Include validation rules and constraints

## 🚀 Next Steps

- **Learn about specific configurations**: [Environment Variables](./environment-variables.md)
- **See how to add new configs**: [Adding New Configs](./adding-new-configs.md)
- **Understand the complete setup**: [Configuration Overview](./overview.md)

---

These patterns provide a solid foundation for managing configuration in complex applications while maintaining security, testability, and maintainability.
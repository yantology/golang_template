package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type DatabaseType string

const (
	PostgreSQL DatabaseType = "postgres"
)

type DatabaseConfig struct {
	Type          DatabaseType  `json:"type"`
	Host          string        `json:"host"`
	Port          string        `json:"port"`
	User          string        `json:"user"`
	Password      string        `json:"password"`
	Name          string        `json:"name"`
	SSLMode       string        `json:"sslmode"`
	MaxOpenConns  int           `json:"max_open_conns"`
	MaxIdleConns  int           `json:"max_idle_conns"`
	MaxLifetime   time.Duration `json:"max_lifetime"`
	MigrationPath string        `json:"migration_path"`
}

// LoadDatabaseConfig loads database configuration from Viper
func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Type:          DatabaseType(viper.GetString("database.type")),
		Host:          viper.GetString("database.host"),
		Port:          viper.GetString("database.port"),
		User:          viper.GetString("database.user"),
		Password:      viper.GetString("database.password"),
		Name:          viper.GetString("database.name"),
		SSLMode:       viper.GetString("database.sslmode"),
		MaxOpenConns:  viper.GetInt("database.max_open_conns"),
		MaxIdleConns:  viper.GetInt("database.max_idle_conns"),
		MaxLifetime:   viper.GetDuration("database.max_lifetime"),
		MigrationPath: viper.GetString("database.migration_path"),
	}
}

// Validate validates database configuration
func (c DatabaseConfig) Validate(isProduction bool) error {
	if c.Type != PostgreSQL {
		return fmt.Errorf("invalid database type: %s (only postgres is supported)", c.Type)
	}

	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Port == "" {
		return fmt.Errorf("database port is required")
	}

	if c.User == "" {
		return fmt.Errorf("database user is required")
	}

	if c.Name == "" {
		return fmt.Errorf("database name is required")
	}

	// In production, password should be provided
	if isProduction && c.Password == "" {
		return fmt.Errorf("database password is required in production environment")
	}

	if c.MaxOpenConns <= 0 {
		return fmt.Errorf("max open connections must be positive")
	}

	if c.MaxIdleConns <= 0 {
		return fmt.Errorf("max idle connections must be positive")
	}

	if c.MaxIdleConns > c.MaxOpenConns {
		return fmt.Errorf("max idle connections cannot be greater than max open connections")
	}

	if c.MaxLifetime <= 0 {
		return fmt.Errorf("connection max lifetime must be positive")
	}

	if c.MigrationPath == "" {
		return fmt.Errorf("migration path is required")
	}

	return nil
}

// GetDSN returns the database connection string
func (c DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

// GetDriverName returns the database driver name
func (c DatabaseConfig) GetDriverName() string {
	return "postgres"
}
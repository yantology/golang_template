package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port            string        `json:"port"`
	Host            string        `json:"host"`
	Env             string        `json:"env"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
	EnableCORS      bool          `json:"enable_cors"`
	CORSOrigins     []string      `json:"cors_origins"`
}

// LoadServerConfig loads server configuration from Viper
func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Port:            viper.GetString("server.port"),
		Host:            viper.GetString("server.host"),
		Env:             viper.GetString("server.env"),
		ReadTimeout:     viper.GetDuration("server.read_timeout"),
		WriteTimeout:    viper.GetDuration("server.write_timeout"),
		IdleTimeout:     viper.GetDuration("server.idle_timeout"),
		ShutdownTimeout: viper.GetDuration("server.shutdown_timeout"),
		EnableCORS:      viper.GetBool("server.enable_cors"),
		CORSOrigins:     viper.GetStringSlice("server.cors_origins"),
	}
}

// Validate validates server configuration
func (c ServerConfig) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Host == "" {
		return fmt.Errorf("server host is required")
	}

	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
		"test":        true,
	}

	if !validEnvs[c.Env] {
		return fmt.Errorf("invalid environment: %s (must be one of: development, staging, production, test)", c.Env)
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("server read timeout must be positive")
	}

	if c.WriteTimeout <= 0 {
		return fmt.Errorf("server write timeout must be positive")
	}

	if c.IdleTimeout <= 0 {
		return fmt.Errorf("server idle timeout must be positive")
	}

	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("server shutdown timeout must be positive")
	}


	return nil
}

// IsProduction returns true if environment is production
func (c ServerConfig) IsProduction() bool {
	return c.Env == "production"
}

// IsDevelopment returns true if environment is development
func (c ServerConfig) IsDevelopment() bool {
	return c.Env == "development"
}

// IsStaging returns true if environment is staging
func (c ServerConfig) IsStaging() bool {
	return c.Env == "staging"
}

// IsTest returns true if environment is test
func (c ServerConfig) IsTest() bool {
	return c.Env == "test"
}

// GetAddress returns the server address in host:port format
func (c ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
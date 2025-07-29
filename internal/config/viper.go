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
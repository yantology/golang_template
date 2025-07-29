package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type LoggerConfig struct {
	Level            string `json:"level"`
	Format           string `json:"format"`
	Output           string `json:"output"`
	EnableCaller     bool   `json:"enable_caller"`
	EnableStacktrace bool   `json:"enable_stacktrace"`
}

// LoadLoggerConfig loads logger configuration from Viper
func LoadLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:            viper.GetString("logger.level"),
		Format:           viper.GetString("logger.format"),
		Output:           viper.GetString("logger.output"),
		EnableCaller:     viper.GetBool("logger.enable_caller"),
		EnableStacktrace: viper.GetBool("logger.enable_stacktrace"),
	}
}

// Validate validates logger configuration
func (c LoggerConfig) Validate() error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
		"panic": true,
	}

	if !validLevels[c.Level] {
		return fmt.Errorf("invalid log level: %s (valid levels: debug, info, warn, error, fatal, panic)", c.Level)
	}

	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}

	if !validFormats[c.Format] {
		return fmt.Errorf("invalid log format: %s (valid formats: json, text)", c.Format)
	}

	validOutputs := map[string]bool{
		"stdout": true,
		"stderr": true,
		"file":   true,
	}

	if !validOutputs[c.Output] {
		return fmt.Errorf("invalid log output: %s (valid outputs: stdout, stderr, file)", c.Output)
	}

	return nil
}
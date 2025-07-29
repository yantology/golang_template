package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type JWTConfig struct {
	Secret           string        `json:"secret"`
	AccessTokenTTL   time.Duration `json:"access_token_ttl"`
	RefreshTokenTTL  time.Duration `json:"refresh_token_ttl"`
	Issuer           string        `json:"issuer"`
	Audience         string        `json:"audience"`
	Algorithm        string        `json:"algorithm"`
}

// LoadJWTConfig loads JWT configuration from Viper
func LoadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:          viper.GetString("jwt.secret"),
		AccessTokenTTL:  viper.GetDuration("jwt.access_token_ttl"),
		RefreshTokenTTL: viper.GetDuration("jwt.refresh_token_ttl"),
		Issuer:          viper.GetString("jwt.issuer"),
		Audience:        viper.GetString("jwt.audience"),
		Algorithm:       viper.GetString("jwt.algorithm"),
	}
}

// Validate validates JWT configuration
func (c JWTConfig) Validate(isProduction bool) error {
	if c.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	// In production, ensure the secret is strong enough
	if isProduction {
		if len(c.Secret) < 32 {
			return fmt.Errorf("JWT secret must be at least 32 characters in production")
		}

		// Check if using default secret
		if c.Secret == "your-super-secret-key-change-this-in-production" {
			return fmt.Errorf("please change the default JWT secret in production")
		}
	}

	if c.AccessTokenTTL <= 0 {
		return fmt.Errorf("access token TTL must be positive")
	}

	if c.RefreshTokenTTL <= 0 {
		return fmt.Errorf("refresh token TTL must be positive")
	}

	if c.AccessTokenTTL >= c.RefreshTokenTTL {
		return fmt.Errorf("refresh token TTL must be greater than access token TTL")
	}

	if c.Issuer == "" {
		return fmt.Errorf("JWT issuer is required")
	}

	if c.Audience == "" {
		return fmt.Errorf("JWT audience is required")
	}

	validAlgorithms := map[string]bool{
		"HS256": true,
		"HS384": true,
		"HS512": true,
		"RS256": true,
		"RS384": true,
		"RS512": true,
		"ES256": true,
		"ES384": true,
		"ES512": true,
	}

	if !validAlgorithms[c.Algorithm] {
		return fmt.Errorf("invalid JWT algorithm: %s", c.Algorithm)
	}

	return nil
}
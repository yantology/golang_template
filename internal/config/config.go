package config

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logger   LoggerConfig   `json:"logger"`
	JWT      JWTConfig      `json:"jwt"`
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "localhost",
			Env:  "development",
		},
	}, nil
}
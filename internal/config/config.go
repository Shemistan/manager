package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// DatabaseConfig represents database configuration.
type DatabaseConfig struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	User    string `toml:"user"`
	Name    string `toml:"name"`
	SSLMode string `toml:"sslmode"`
}

// Config represents the application configuration.
type Config struct {
	ServiceName string         `toml:"service_name"`
	ServiceEnv  string         `toml:"service_env"`
	HTTPPort    int            `toml:"http_port"`
	Database    DatabaseConfig `toml:"database"`
}

// Load reads the configuration from a TOML file.
func Load(filePath string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	return &cfg, nil
}

// GetDatabasePassword retrieves the database password from environment variables.
func GetDatabasePassword() string {
	return os.Getenv("DB_PASSWORD")
}

// GetDatabaseUser retrieves the database user from environment variables (or uses config value).
func GetDatabaseUser(configUser string) string {
	envUser := os.Getenv("DB_USER")
	if envUser != "" {
		return envUser
	}
	return configUser
}

// BuildDSN constructs a PostgreSQL DSN string.
func BuildDSN(cfg *Config) string {
	user := GetDatabaseUser(cfg.Database.User)
	password := GetDatabasePassword()
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		user,
		password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
}

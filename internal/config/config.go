package config

import (
	"fmt"
	"os"
	"strings"
)

// DatabaseConfig represents database configuration.
type DatabaseConfig struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	User    string `toml:"user"`
	Name    string `toml:"name"`
	SSLMode string `toml:"sslmode"`
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enabled  bool   `toml:"enabled"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
	CAFile   string `toml:"ca_file"`
}

// Config represents the application configuration.
type Config struct {
	ServiceName string         `toml:"service_name"`
	ServiceEnv  string         `toml:"service_env"`
	HTTPPort    int            `toml:"http_port"`
	Database    DatabaseConfig `toml:"database"`
	TLS         TLSConfig      `toml:"tls"`
}

// Load reads the configuration from a TOML file and environment variables.
func Load() (*Config, error) {
	var cfg Config

	// Override with environment variables
	// Database configuration
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Database.Port) // nolint:errcheck
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		cfg.Database.Name = name
	}
	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		cfg.Database.SSLMode = sslmode
	}

	// App configuration
	if port := os.Getenv("SERVICE_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.HTTPPort) // nolint:errcheck
	}

	// TLS configuration
	if tlsEnabled := os.Getenv("TLS_ENABLED"); tlsEnabled != "" {
		cfg.TLS.Enabled = strings.ToLower(tlsEnabled) == "true"
	}
	if certFile := os.Getenv("TLS_CERT_FILE"); certFile != "" {
		cfg.TLS.CertFile = certFile
	}
	if keyFile := os.Getenv("TLS_KEY_FILE"); keyFile != "" {
		cfg.TLS.KeyFile = keyFile
	}
	if caFile := os.Getenv("TLS_CA_FILE"); caFile != "" {
		cfg.TLS.CAFile = caFile
	}

	// Set defaults
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}

	return &cfg, nil
}

// GetDatabasePassword retrieves the database password from environment variables.
func GetDatabasePassword() string {
	return os.Getenv("DB_PASSWORD")
}

// GetDatabaseUser retrieves the database user from environment variables (or uses config value).
func GetDatabaseUser() string {
	return os.Getenv("DB_USER")
}

// BuildDSN constructs a PostgreSQL DSN string.
func BuildDSN(cfg *Config) string {
	user := GetDatabaseUser()
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

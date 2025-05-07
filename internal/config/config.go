// Package config provides functionality for loading application configuration from a YAML file.
package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config holds the application's configuration settings loaded from a YAML file.
//
// Fields:
//   - Port: the port number the server should listen on
//   - Debug: whether debug mode is enabled
//   - Database: the type of database used (e.g., "sqlite", "postgres")
//   - DatabaseDSN: the database connection string (DSN)
//   - AuthProvider: the authentication provider ("local", "ldap", etc.)
//   - JWTSecret: the secret key used for signing JWT tokens
type Config struct {
	Port         int    `yaml:"port"`
	Debug        bool   `yaml:"debug"`
	Database     string `yaml:"database"`      // "sqlite", "postgres", "mysql", etc.
	DatabaseDSN  string `yaml:"database_dsn"`  // connection string for the database
	AuthProvider string `yaml:"auth_provider"` // "local", "ldap", etc.
	JWTSecret    string `yaml:"jwt_secret"`    // secret for JWT token signing
}

// LoadConfig loads a YAML configuration file from the specified path.
// Environment variables take precedence over values from the config file.
//
// It reads the file contents, unmarshals the YAML into a Config struct,
// and overrides applicable fields using environment variables if set.
//
// Parameters:
//   - path: the file path to the YAML config file
//
// Returns:
//   - *Config: pointer to the loaded Config struct
//   - error: non-nil if reading or parsing fails
func LoadConfig(path string) (*Config, error) {
	if envPath := os.Getenv("IAM_CONFIG_PATH"); envPath != "" {
		// Allow overriding the config path using IAM_CONFIG_PATH environment variable
		path = envPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if portStr := os.Getenv("IAM_PORT"); portStr != "" {
		// Override YAML port with environment variable IAM_PORT
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = port
		}
	}

	if db := os.Getenv("IAM_DATABASE"); db != "" {
		// Override database engine from environment IAM_DATABASE
		cfg.Database = db
	}

	if dsn := os.Getenv("IAM_DATABASE_DSN"); dsn != "" {
		// Override database DSN from environment IAM_DATABASE_DSN
		cfg.DatabaseDSN = dsn
	}

	if provider := os.Getenv("IAM_AUTH_PROVIDER"); provider != "" {
		// Override auth provider from environment IAM_AUTH_PROVIDER
		cfg.AuthProvider = provider
	}

	return &cfg, nil
}

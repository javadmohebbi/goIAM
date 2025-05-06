// Package config provides functionality for loading application configuration from a YAML file.
package config

import (
	"os"

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
//
// It reads the file contents, unmarshals the YAML into a Config struct,
// and returns a pointer to the Config instance or an error if loading fails.
//
// Parameters:
//   - path: the file path to the YAML config file
//
// Returns:
//   - *Config: pointer to the loaded Config struct
//   - error: non-nil if reading or parsing fails
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

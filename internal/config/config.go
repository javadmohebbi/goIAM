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
//   - AuthProviders: a list of authentication providers ("local", "ldap", etc.)
//   - JWTSecret: the secret key used for signing JWT tokens
type Config struct {
	Port          int                  `yaml:"port"`
	Debug         bool                 `yaml:"debug"`
	Database      string               `yaml:"database"`       // "sqlite", "postgres", "mysql", etc.
	DatabaseDSN   string               `yaml:"database_dsn"`   // connection string for the database
	AuthProviders []AuthProviderConfig `yaml:"auth_providers"` // list of typed auth providers
	JWTSecret     string               `yaml:"jwt_secret"`     // secret for JWT token signing
	Validation    ValidationConfig     `yaml:"validation"`
	AppName       string               `yaml:"appName"`    // Application name used in CLI and logs
	ServerName    string               `yaml:"serverName"` // Server name for headers or UI
}

type ValidationConfig struct {
	EmailRegex        string `yaml:"email_regex"`
	PhoneRegex        string `yaml:"phone_regex"`
	PasswordRegex     string `yaml:"password_regex"`
	WebsiteRegex      string `yaml:"website_regex"`
	PasswordMinLength int    `yaml:"password_min_length"`
}

// AuthProviderConfig holds a named provider and its configuration block (as raw map).
type AuthProviderConfig struct {
	Name   string                 `yaml:"name"`   // e.g., "local", "ldap", "firebase"
	Config map[string]interface{} `yaml:"config"` // raw provider-specific configuration
}

// LDAPConfig defines configuration fields for an LDAP provider.
type LDAPConfig struct {
	Server string `yaml:"server"`
	BaseDN string `yaml:"baseDn"`
}

// Auth0Config defines configuration fields for an Auth0 provider.
type Auth0Config struct {
	Domain       string `yaml:"domain"`
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	Audience     string `yaml:"audience"`
}

// EntraIDConfig defines configuration fields for Microsoft Entra ID.
type EntraIDConfig struct {
	TenantID      string   `yaml:"tenantId"`
	ClientID      string   `yaml:"clientId"`
	ClientSecret  string   `yaml:"clientSecret"`
	AuthorityHost string   `yaml:"authorityHost"`
	Scopes        []string `yaml:"scopes"`
}

// FirebaseConfig defines configuration fields for Firebase Authentication.
type FirebaseConfig struct {
	ProjectID       string `yaml:"projectId"`
	CredentialsFile string `yaml:"credentialsFile"`
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

	// Apply default validation config if not set
	if cfg.Validation.EmailRegex == "" {
		cfg.Validation.EmailRegex = `^[^@\s]+@[^@\s]+\.[^@\s]+$`
	}
	if cfg.Validation.PhoneRegex == "" {
		cfg.Validation.PhoneRegex = `^\+?[0-9]{7,15}$`
	}
	if cfg.Validation.PasswordRegex == "" {
		cfg.Validation.PasswordRegex = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&]{6,}$`
	}
	if cfg.Validation.WebsiteRegex == "" {
		cfg.Validation.WebsiteRegex = `^https?://[\w\-\.]+\.\w+`
	}
	if cfg.Validation.PasswordMinLength == 0 {
		cfg.Validation.PasswordMinLength = 6
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

	// Deprecated: single string override is not supported for multi-provider mode
	// if provider := os.Getenv("IAM_AUTH_PROVIDER"); provider != "" {
	// 	cfg.AuthProvider = provider
	// }

	if appName := os.Getenv("IAM_APP_NAME"); appName != "" {
		cfg.AppName = appName
	}
	if serverName := os.Getenv("IAM_SERVER_NAME"); serverName != "" {
		cfg.ServerName = serverName
	}

	return &cfg, nil
}

// As unmarshals the raw Config map into the provided target struct.
//
// Example usage:
//
//	var cfg LDAPConfig
//	if err := provider.As(&cfg); err != nil { ... }
func (ap *AuthProviderConfig) As(target any) error {
	data, err := yaml.Marshal(ap.Config)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, target)
}

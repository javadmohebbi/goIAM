package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port         int    `yaml:"port"`
	Debug        bool   `yaml:"debug"`
	Database     string `yaml:"database"`      // "sqlite", "postgres", "mysql", "sqlserver", "clickhouse"
	DatabaseDSN  string `yaml:"database_dsn"`  // connection string
	AuthProvider string `yaml:"auth_provider"` // "local", "ldap", etc.
	JWTSecret    string `yaml:"jwt_secret"`
}

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

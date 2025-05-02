package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port         int    `yaml:"port"`
	Debug        bool   `yaml:"debug"`
	DatabaseDSN  string `yaml:"database_dsn"`
	AuthProvider string `yaml:"auth_provider"` // e.g., local, ldap, auth0
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

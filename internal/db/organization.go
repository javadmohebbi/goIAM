package db

import "gorm.io/gorm"

// Organization represents a tenant in the multi-tenant IAM system.
// It provides logical isolation for users, groups, roles, and policies.
//
// Slug is a short, URL-safe identifier used for subdomain-based routing
// and stable human-readable identifiers in the API.
type Organization struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"` // Unique organization name
	Slug        string `gorm:"uniqueIndex;not null"` // Short, URL-safe identifier for subdomain and API routing
	Description string // Optional description of the organization
	Users       []User // Users in the organization
}

// Package db defines the IAM (Identity and Access Management) data models,
// including Group, Role, and Policy relationships used by the authentication and authorization system.
package db

import "gorm.io/gorm"

// Organization represents a tenant in the multi-tenant IAM system.
// It provides logical isolation for users, groups, roles, and policies.
type Organization struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex;not null"` // Unique organization name
	Users []User // Users in the organization
}

type Group struct {
	gorm.Model
	Name           string `gorm:"not null;uniqueIndex:idx_org_group_name"` // Unique within org
	OrganizationID uint   // Tenant scoping
	Organization   Organization
	Users          []User   `gorm:"many2many:user_groups;"`
	Policies       []Policy `gorm:"many2many:group_policies;"`
}

// Role represents a job-based permission set within an organization.
//
// Fields:
//   - Name: a name unique within the organization
//   - OrganizationID: foreign key to the organization this role belongs to
//   - Organization: the organization entity this role belongs to
//   - Users: users assigned this role (many-to-many)
//   - Policies: policies assigned to this role (many-to-many)
type Role struct {
	gorm.Model
	Name           string `gorm:"not null;uniqueIndex:idx_org_role_name"` // Unique within org
	OrganizationID uint
	Organization   Organization
	Users          []User   `gorm:"many2many:user_roles;"`
	Policies       []Policy `gorm:"many2many:role_policies;"`
}

// Policy defines access control rules scoped to an organization.
//
// Fields:
//   - Name: a name unique within the organization
//   - OrganizationID: foreign key to the organization this policy belongs to
//   - Organization: the organization entity this policy belongs to
//   - Description: optional human-readable description of the policy
type Policy struct {
	gorm.Model
	Name           string `gorm:"not null;uniqueIndex:idx_org_policy_name"` // Unique within org
	OrganizationID uint
	Organization   Organization
	Description    string
}

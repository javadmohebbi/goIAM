package db

import "gorm.io/gorm"

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
	Slug           string `gorm:"uniqueIndex:idx_org_role_slug"`          // Unique slug within organization
	Description    string // Optional role description
	OrganizationID uint
	Organization   Organization
	Users          []User   `gorm:"many2many:user_roles;"`
	Policies       []Policy `gorm:"many2many:role_policies;"`
}

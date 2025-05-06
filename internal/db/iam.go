// Package db defines the IAM (Identity and Access Management) data models,
// including Group, Role, and Policy relationships used by the authentication and authorization system.
package db

import "gorm.io/gorm"

// Group represents a collection of users that can share common policies.
//
// Fields:
//   - Name: a unique name for the group
//   - Users: the users belonging to this group (many-to-many)
//   - Policies: policies applied to this group (many-to-many)
type Group struct {
	gorm.Model
	Name     string   `gorm:"uniqueIndex;not null"`      // Unique group name
	Users    []User   `gorm:"many2many:user_groups;"`    // Users in the group
	Policies []Policy `gorm:"many2many:group_policies;"` // Group-wide policies
}

// Role represents a set of permissions assigned to users, usually by job function.
//
// Fields:
//   - Name: a unique name for the role
//   - Users: users assigned this role (many-to-many)
//   - Policies: policies associated with this role (many-to-many)
type Role struct {
	gorm.Model
	Name     string   `gorm:"uniqueIndex;not null"`     // Unique role name
	Users    []User   `gorm:"many2many:user_roles;"`    // Users with this role
	Policies []Policy `gorm:"many2many:role_policies;"` // Role-wide policies
}

// Policy defines a named access control rule that can be linked to groups or roles.
//
// Fields:
//   - Name: unique identifier for the policy
//   - Description: human-readable explanation of the policy's intent
type Policy struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"` // Unique policy name
	Description string // Optional description
}

package db

import "gorm.io/gorm"

// Group represents a collection of users within an organization,
// typically used to assign shared access policies.
//
// Each group belongs to a single organization and may have many users and policies.
type Group struct {
	gorm.Model

	Name           string       `gorm:"not null;uniqueIndex:idx_org_group_name"` // Unique group name within the organization
	Slug           string       `gorm:"uniqueIndex:idx_org_group_slug"`          // URL-safe identifier for routing or CLI use
	Description    string       // Optional description for the group
	OrganizationID uint         // Foreign key reference to the owning organization
	Organization   Organization // GORM association to the organization

	Users    []User   `gorm:"many2many:user_groups;"`    // Many-to-many relationship with users
	Policies []Policy `gorm:"many2many:group_policies;"` // Many-to-many relationship with policies
}

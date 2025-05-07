package db

import (
	"errors"

	"gorm.io/gorm"
)

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

// CreateRole creates a new Role in the database.
//
// Parameters:
//   - role: pointer to a Role struct with required fields populated.
//
// Returns an error if the insertion fails.
func CreateRole(role *Role) error {
	return DB.Create(role).Error
}

// GetRoleByID retrieves a Role by its ID.
//
// Parameters:
//   - id: the primary key of the role.
//
// Returns the role and an error if not found or query fails.
func GetRoleByID(id uint) (*Role, error) {
	var role Role
	if err := DB.Preload("Users").Preload("Policies").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// UpdateRole updates the provided Role record in the database.
//
// Parameters:
//   - role: pointer to a Role struct with updates applied.
//
// Returns an error if the update fails.
func UpdateRole(role *Role) error {
	return DB.Save(role).Error
}

// DeleteRole deletes a Role from the database by its ID.
//
// Parameters:
//   - id: the primary key of the role to delete.
//
// Returns an error if deletion fails.
func DeleteRole(id uint) error {
	if id == 0 {
		return errors.New("invalid role ID")
	}
	return DB.Delete(&Role{}, id).Error
}

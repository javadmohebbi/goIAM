package db

import (
	"errors"

	"gorm.io/gorm"
)

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

// CreateGroup inserts a new Group into the database.
//
// Parameters:
//   - group: pointer to a Group struct with required fields populated.
//
// Returns an error if the operation fails.
func CreateGroup(group *Group) error {
	return DB.Create(group).Error
}

// GetGroupByID retrieves a Group by its primary key.
//
// Parameters:
//   - id: the group's ID.
//
// Returns the group and an error if not found.
func GetGroupByID(id uint) (*Group, error) {
	var group Group
	if err := DB.Preload("Users").Preload("Policies").First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// UpdateGroup saves changes to an existing Group.
//
// Parameters:
//   - group: pointer to the modified Group.
//
// Returns an error if the update fails.
func UpdateGroup(group *Group) error {
	return DB.Save(group).Error
}

// DeleteGroup removes a Group from the database by its ID.
//
// Parameters:
//   - id: the group's ID.
//
// Returns an error if the delete operation fails.
func DeleteGroup(id uint) error {
	if id == 0 {
		return errors.New("invalid group ID")
	}
	return DB.Delete(&Group{}, id).Error
}

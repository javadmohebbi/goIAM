// Package db provides data models and database access methods for policy resources,
// including CRUD operations and JSON import/export utilities.
package db

import (
	"encoding/json"
	"os"

	"gorm.io/gorm"
)

// PolicyAction defines an allowed or denied action in a policy statement.
type PolicyAction struct {
	gorm.Model
	PolicyStatementID uint
	Action            string
}

// PolicyResource defines a resource the statement applies to.
type PolicyResource struct {
	gorm.Model
	PolicyStatementID uint
	Resource          string
	// OrganizationID references the organization this resource belongs to.
	OrganizationID uint
}

// CreateResource creates a new PolicyResource associated with a given organization and policy statement.
func CreateResource(orgID, statementID uint, name string) error {
	return DB.Create(&PolicyResource{
		OrganizationID:    orgID,
		PolicyStatementID: statementID,
		Resource:          name,
	}).Error
}

// GetResourcesByOrg retrieves all PolicyResource records for the specified organization.
func GetResourcesByOrg(orgID uint) ([]PolicyResource, error) {
	var resources []PolicyResource
	err := DB.Where("organization_id = ?", orgID).Find(&resources).Error
	return resources, err
}

// UpdateResourceName updates the Resource field of a PolicyResource by its ID.
func UpdateResourceName(id uint, newName string) error {
	return DB.Model(&PolicyResource{}).Where("id = ?", id).Update("resource", newName).Error
}

// DeleteResource removes a PolicyResource from the database by ID.
func DeleteResource(id uint) error {
	return DB.Delete(&PolicyResource{}, id).Error
}

// Each record in the file must conform to the PolicyResource structure. The function expects
// a JSON array of objects, each including the fields: PolicyStatementID, Resource, and OrganizationID.
//
// Example JSON file:
// [
//
//	{
//	  "PolicyStatementID": 100,
//	  "Resource": "org:42:users:create",
//	  "OrganizationID": 42
//	},
//	{
//	  "PolicyStatementID": 101,
//	  "Resource": "org:42:groups:read",
//	  "OrganizationID": 42
//	}
//
// ]
//
// Parameters:
//   - path: the file path to the JSON file.
//
// Returns an error if the file cannot be read, parsed, or if database insertion fails.
func ImportResourcesFromJSON(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var resources []PolicyResource
	if err := json.Unmarshal(file, &resources); err != nil {
		return err
	}
	return DB.Create(&resources).Error
}

// Each record in the file must conform to the PolicyResource structure. The function expects
// a JSON array of objects, each including the fields: PolicyStatementID, Resource, and OrganizationID.
//
// Example JSON file:
// [
//
//	{
//	  "PolicyStatementID": 100,
//	  "Resource": "org:42:users:create",
//	  "OrganizationID": 42
//	},
//	{
//	  "PolicyStatementID": 101,
//	  "Resource": "org:42:groups:read",
//	  "OrganizationID": 42
//	}
//
// ]
//
// Parameters:
//   - path: the file path to the JSON file.
//
// Returns an error if the file cannot be read, parsed, or if database insertion fails.
func ExportResourcesToJSON(path string, orgID uint) error {
	resources, err := GetResourcesByOrg(orgID)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

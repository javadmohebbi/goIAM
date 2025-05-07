package db

import (
	"errors"

	"gorm.io/gorm"
)

// Policy defines access control rules scoped to an organization.
// Fields:
//   - Name: a name unique within the organization
//   - OrganizationID: foreign key to the organization this policy belongs to
//   - Organization: the organization entity this policy belongs to
//   - Description: optional human-readable description of the policy
type Policy struct {
	gorm.Model
	Name           string `gorm:"not null;uniqueIndex:idx_org_policy_name"` // Unique within org
	Slug           string `gorm:"uniqueIndex:idx_org_policy_slug"`          // Unique slug within organization
	OrganizationID uint
	Organization   Organization
	Description    string
	Statements     []PolicyStatement `gorm:"foreignKey:PolicyID"` // List of statements under this policy
}

// PolicyStatement represents a single rule inside a policy,
// describing an effect ("Allow" or "Deny") and its related actions and resources.
type PolicyStatement struct {
	gorm.Model
	PolicyID uint
	Policy   Policy
	Effect   string // "Allow" or "Deny"

	Actions   []PolicyAction   `gorm:"foreignKey:PolicyStatementID"`
	Resources []PolicyResource `gorm:"foreignKey:PolicyStatementID"`
}

// CreatePolicy inserts a new Policy into the database.
//
// Parameters:
//   - policy: pointer to a Policy struct with required fields.
//
// Returns an error if insertion fails.
func CreatePolicy(policy *Policy) error {
	return DB.Create(policy).Error
}

// GetPolicyByID retrieves a Policy and its statements by ID.
//
// Parameters:
//   - id: the Policy's ID.
//
// Returns the policy and an error if not found.
func GetPolicyByID(id uint) (*Policy, error) {
	var policy Policy
	err := DB.Preload("Statements.Actions").
		Preload("Statements.Resources").
		First(&policy, id).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// UpdatePolicy updates the provided Policy record.
//
// Parameters:
//   - policy: pointer to the modified Policy.
//
// Returns an error if update fails.
func UpdatePolicy(policy *Policy) error {
	return DB.Save(policy).Error
}

// DeletePolicy removes a Policy from the database by its ID.
//
// Parameters:
//   - id: the policy's ID.
//
// Returns an error if the delete operation fails.
func DeletePolicy(id uint) error {
	if id == 0 {
		return errors.New("invalid policy ID")
	}
	return DB.Delete(&Policy{}, id).Error
}

package db

import "gorm.io/gorm"

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

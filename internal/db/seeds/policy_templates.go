package seeds

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/javadmohebbi/goIAM/internal/db"
	"gorm.io/gorm"
)

// SeedDefaultPoliciesForOrg seeds standard policies for a new organization.
//
// It ensures that each organization has:
//   - FullAccess: wildcard access to everything
//   - ReadOnly: read-only actions on all resources
//   - SelfManage: user can manage their own profile
func SeedDefaultPoliciesForOrg(orgID uint, dbConn *gorm.DB) error {
	// Define default policy templates with their statements.
	// Each template includes actions and resources under an effect.
	policies := []struct {
		Name        string
		Slug        string
		Description string
		Statements  []db.PolicyStatement
	}{
		{
			Name:        "FullAccess",
			Slug:        "full-access",
			Description: "Grants full access to all actions and resources.",
			Statements: []db.PolicyStatement{{
				Effect: "Allow",
				Actions: []db.PolicyAction{{
					Action: "*",
				}},
				Resources: []db.PolicyResource{{
					Resource:       "*",
					OrganizationID: orgID,
				}},
			}},
		},
		{
			Name:        "ReadOnly",
			Slug:        "read-only",
			Description: "Grants read-only access to users, roles, groups, and policies.",
			Statements: []db.PolicyStatement{{
				Effect: "Allow",
				Actions: []db.PolicyAction{
					{Action: "user:read"},
					{Action: "group:read"},
					{Action: "role:read"},
					{Action: "policy:read"},
				},
				Resources: []db.PolicyResource{{
					Resource:       "*",
					OrganizationID: orgID,
				}},
			}},
		},
		{
			Name:        "SelfManage",
			Slug:        "self-manage",
			Description: "Allows users to manage their own profile.",
			Statements: []db.PolicyStatement{{
				Effect: "Allow",
				Actions: []db.PolicyAction{
					{Action: "user:read"},
					{Action: "user:update"},
				},
				Resources: []db.PolicyResource{{
					Resource:       "org:{org_id}:user:{user_id}",
					OrganizationID: orgID,
				}},
			}},
		},
	}

	// Loop through each policy template and attempt to create it if it doesn't exist.
	for _, tpl := range policies {
		// Check if a policy with the same slug already exists for this organization.
		var existing db.Policy
		err := dbConn.Where("slug = ? AND organization_id = ?", tpl.Slug, orgID).First(&existing).Error
		if err == nil {
			continue // Already exists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("db error: %w", err)
		}

		// If policy doesn't exist, construct a new one with a unique slug and attach statements.
		newPolicy := db.Policy{
			Name:           tpl.Name,
			Slug:           tpl.Slug + "-" + uuid.New().String()[:6],
			Description:    tpl.Description,
			OrganizationID: orgID,
			Statements:     tpl.Statements,
		}

		// Insert the new policy and its related statements/actions/resources.
		if err := dbConn.Create(&newPolicy).Error; err != nil {
			return fmt.Errorf("failed to seed policy %s: %w", tpl.Name, err)
		}
	}

	return nil
}

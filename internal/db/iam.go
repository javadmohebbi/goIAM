// Package db defines the IAM (Identity and Access Management) data models,
// including Group, Role, and Policy relationships used by the authentication and authorization system.
package db

// EvaluatePolicy determines whether a user is allowed to perform the specified action on the given resource.
//
// It aggregates all policies assigned to the user directly, via groups, and via roles,
// then evaluates their policy statements by checking matching actions, resources, and effects.
//
// Returns true if an "Allow" policy applies and is not overridden by a matching "Deny".
func EvaluatePolicy(user User, action string, resource string) bool {
	// Gather all relevant policy IDs
	policyIDs := map[uint]struct{}{}

	// User's direct policies
	for _, p := range user.Policies {
		policyIDs[p.ID] = struct{}{}
	}
	// Group policies
	for _, g := range user.Groups {
		for _, p := range g.Policies {
			policyIDs[p.ID] = struct{}{}
		}
	}
	// Role policies
	for _, r := range user.Roles {
		for _, p := range r.Policies {
			policyIDs[p.ID] = struct{}{}
		}
	}

	// Track effective decision
	allowed := false

	for pid := range policyIDs {
		var policy Policy
		// Preload Actions and Resources for each statement
		if err := DB.Preload("Statements.Actions").Preload("Statements.Resources").First(&policy, pid).Error; err != nil {
			continue
		}

		for _, stmt := range policy.Statements {
			actionMatch := false
			for _, a := range stmt.Actions {
				if a.Action == action || a.Action == "*" {
					actionMatch = true
					break
				}
			}

			resourceMatch := false
			for _, r := range stmt.Resources {
				if (r.Resource == resource || r.Resource == "*") && r.OrganizationID == user.OrganizationID {
					resourceMatch = true
					break
				}
			}

			if actionMatch && resourceMatch {
				if stmt.Effect == "Deny" {
					return false // Deny overrides everything
				}
				if stmt.Effect == "Allow" {
					allowed = true
				}
			}
		}
	}

	return allowed
}

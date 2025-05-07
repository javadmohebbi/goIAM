package db

import "gorm.io/gorm"

// PolicyAction defines an allowed or denied action in a statement.
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
}

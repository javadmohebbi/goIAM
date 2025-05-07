// Package validation provides reusable input validation utilities that are
// configurable via the main application config.
package validation

import "github.com/javadmohebbi/goIAM/internal/config"

// Validation encapsulates input validation logic using configuration rules
// such as regex patterns and length constraints defined in the application's config.
type Validation struct {
	cfg *config.Config
}

// New creates a new Validation instance using the provided configuration.
func New(c *config.Config) *Validation {
	return &Validation{
		cfg: c,
	}
}

package api

import (
	"time"

	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/validation"
)

// API provides shared dependencies to API route handlers.
//
// It holds the application configuration and a centralized validation utility.
type API struct {
	cfg        *config.Config
	validation *validation.Validation

	startTime time.Time
}

// New returns a new instance of the API struct,
// initialized with configuration and validation logic.
func New(c *config.Config) *API {
	return &API{
		cfg:        c,
		validation: validation.New(c),
	}
}

// Package api contains HTTP server initialization and route registration.
package api

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

// StartServer initializes and starts the Fiber HTTP server.
//
// It registers basic health check routes and, depending on the configured
// authentication provider, it sets up appropriate route groups (e.g., local auth).
//
// Parameters:
//   - cfg: pointer to the application's Config struct containing port, auth provider, and other settings.
//
// Returns:
//   - error if the server fails to start.
func (a *API) StartServer() error {
	// Initialize Fiber app instance
	app := fiber.New()

	// Simple health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("healthy")
	})

	// Register routes depending on auth provider
	if a.cfg.AuthProvider == "local" {
		a.RegisterLocalRoutes(app)
	}

	// Start the server on the specified port
	return app.Listen(fmt.Sprintf(":%d", a.cfg.Port))
}

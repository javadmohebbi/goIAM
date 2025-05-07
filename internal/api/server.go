// Package api contains HTTP server initialization and route registration.
package api

import (
	"fmt"
	"log"
	"runtime"
	"time"

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
	app := fiber.New(
		fiber.Config{
			ServerHeader: a.cfg.ServerName,
			AppName:      a.cfg.AppName,
		},
	)

	// store fiber app in API
	a._app = app

	// server start time
	a.startTime = time.Now()

	// Simple health check endpoint
	app.Get("/echo", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":        "ok",
			"uptime":        time.Since(a.startTime).String(),
			"go_version":    runtime.Version(),
			"auth_provider": a.cfg.AuthProviders,
			"port":          a.cfg.Port,
			"app":           a._app.Config(),
		})
	})

	// Log all incoming requests
	app.Use(func(c fiber.Ctx) error {
		log.Printf("Incoming %s request to %s\n", c.Method(), c.Path())
		// fmt.Printf("---BEGIN--")
		// log.Println(c.Request().String())
		// log.Println(c.Request().Body())
		// fmt.Printf("---END--")
		return c.Next()
	})

	// Register routes depending on auth provider
	a.registerRoutes(app)

	// Start the server on the specified port
	return app.Listen(
		fmt.Sprintf(":%d", a.cfg.Port), // listen on port and all interfaces
		fiber.ListenConfig{
			DisableStartupMessage: false, // if true, it won't show start up message of fiber
			EnablePrintRoutes:     false, // it's false by default, if true, it will return all routes and groups
		},
	)
}

// stop and close the API server
func (a *API) StopAndClose() {
	// for future use
	log.Println("Stopping the server...!")

	// shut down the app
	a._app.RebuildTree().Shutdown()

	log.Println("Stopped!")
}

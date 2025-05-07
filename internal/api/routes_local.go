// Package api provides HTTP route handlers for goIAM.
// This file contains local authentication and 2FA setup routes.
package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/middleware"
)

// RegisterLocalRoutes sets up all routes used for local authentication and 2FA.
//
// It includes endpoints for:
//   - user registration
//   - login with optional 2FA
//   - 2FA setup and validation
//   - regenerating backup codes
//   - disabling 2FA
func (a *API) RegisterLocalRoutes(app *fiber.App) {
	app.Post("/auth/register", a.handleRegister)
	app.Post("/auth/login", a.handleLogin())

	secure := app.Group("/secure", middleware.RequireAuth(a.cfg))

	secure.Post("/auth/2fa/setup", a.handle2FASetup())
	secure.Post("/auth/2fa/verify", a.handle2FAVerify())
	secure.Post("/auth/2fa/disable", a.handle2FADisable())
	secure.Post("/auth/backup-codes/regenerate", a.handleBackupCodes())
}

package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/middleware"
)

// registerAuthRoutes defines all secure authentication-related routes.
//
// These routes are grouped under the /secure prefix and protected by RequireAuth.
// They also apply fine-grained policy checks using RequireAccess middleware.
// Includes routes for 2FA management, user profile updates, and backup codes.
func (a *API) registerAuthRoutes(secure fiber.Router) {
	secure.Post("/auth/2fa/setup", a.handle2FASetup())
	secure.Post("/auth/2fa/verify", a.handle2FAVerify())
	secure.Post("/auth/2fa/disable", a.handle2FADisable())
	secure.Post("/auth/backup-codes/regenerate", a.handleBackupCodes())

	secure.Get("/auth/profile",
		middleware.RequireAccess("read", "org:{org_id}:user:{user_id}", a.cfg),
		a.handleGetProfile)

	secure.Patch("/auth/profile",
		middleware.RequireAccess("update", "org:{org_id}:user:{user_id}", a.cfg),
		a.handleUpdateProfile)

	secure.Post("/auth/profile/password",
		middleware.RequireAccess("update", "org:{org_id}:user:{user_id}:password", a.cfg),
		a.handleChangePassword)

	secure.Post("/auth/profile/2fa/enable",
		middleware.RequireAccess("update", "org:{org_id}:user:{user_id}:2fa", a.cfg),
		a.handleEnable2FA)

	secure.Post("/auth/profile/2fa/disable",
		middleware.RequireAccess("update", "org:{org_id}:user:{user_id}:2fa", a.cfg),
		a.handleDisable2FA)
}

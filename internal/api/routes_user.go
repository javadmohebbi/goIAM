package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/middleware"
)

// registerUserRoutes defines routes for managing users within the authenticated user's organization.
func (a *API) registerUserRoutes(secure fiber.Router) {
	// Create a new user within the caller's organization
	secure.Post("/create",
		middleware.RequireAccess("create", "org:{org_id}:user", a.cfg),
		a.handleCreateUser)

	// // Update an existing user by ID
	// secure.Patch("/:username",
	// 	middleware.RequireAccess("update", "org:{org_id}:user:{user_id}", a.cfg, a.iamDB),
	// 	a.handleUpdateUser)

	// // Delete a user by ID
	// secure.Delete("/:username",
	// 	middleware.RequireAccess("delete", "org:{org_id}:user:{user_id}", a.cfg, a.iamDB),
	// 	a.handleDeleteUser)

	// // Get a specific user by ID
	// secure.Get("/:username",
	// 	middleware.RequireAccess("read", "org:{org_id}:user:{user_id}", a.cfg, a.iamDB),
	// 	a.handleGetUser)

	// // List all users in the caller's organization
	// secure.Get("/users",
	// 	middleware.RequireAccess("read", "org:{org_id}:user", a.cfg, a.iamDB),
	// 	a.handleListUsers)
}

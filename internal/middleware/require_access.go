package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
	"gorm.io/gorm"
)

// RequireAccess returns a Fiber middleware that enforces access control by evaluating the user's policy.
//
// Parameters:
//   - action: the action being performed (e.g., "read", "write", "delete").
//   - resourceTemplate: a string representing the resource with optional placeholders like {user_id}, {org_id}.
//   - cfg: application configuration reference.
//   - dbConn: database connection used for policy evaluation.
func RequireAccess(action string, resourceTemplate string, cfg *config.Config, dbConn *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		user, ok := c.Locals("user").(db.User)
		if !ok {
			return fiber.ErrUnauthorized
		}

		// Replace placeholders in resource template
		resource := strings.ReplaceAll(resourceTemplate, "{user_id}", fmt.Sprint(user.ID))
		resource = strings.ReplaceAll(resource, "{org_id}", fmt.Sprint(user.OrganizationID))

		// Evaluate access
		if !db.EvaluatePolicy(user, action, resource) {
			return fiber.NewError(fiber.StatusForbidden, "Access denied")
		}

		return c.Next()
	}
}

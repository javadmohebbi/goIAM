package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
	"gorm.io/gorm"
)

// RequireAuth is a Fiber middleware that verifies the Authorization Bearer JWT token,
// checks if the user exists, and enforces 2FA if required.
// On success, it stores the `db.User` in c.Locals("user") for route handlers.
func RequireAuth(cfg *config.Config, iamDB *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Extract bearer token
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and verify JWT
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid claims")
		}

		// Extract user ID from token
		userID, ok := claims["sub"].(float64)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user ID")
		}

		// Load user from DB
		var user db.User
		if err := iamDB.First(&user, uint(userID)).Error; err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "user not found")
		}

		// Check if 2FA is required but not verified
		// Skip 2FA check only for /2fa/verify and /2fa/setup
		path := c.Path()
		// verified := claims["2fa"] == true
		// A JWT with "2fa": true means user already passed 2FA
		verified, _ := claims["2fa"].(bool)

		if user.Requires2FA && !verified &&
			path != "/s/auth/2fa/verify" && path != "/s/auth/2fa/setup" {
			return fiber.NewError(fiber.StatusForbidden, "2FA required")
		}

		// Store user object in Fiber context
		c.Locals("user", user)
		return c.Next()
	}
}

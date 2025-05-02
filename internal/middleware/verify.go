package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
)

func RequireAuth(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

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

		userID, ok := claims["sub"].(float64)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user ID")
		}

		var user db.User
		if err := db.DB.First(&user, uint(userID)).Error; err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "user not found")
		}

		if user.Requires2FA && !user.TwoFAVerified {
			return fiber.NewError(fiber.StatusForbidden, "2FA required")
		}

		c.Locals("user", user)
		return c.Next()
	}
}

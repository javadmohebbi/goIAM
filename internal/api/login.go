package api

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/javadmohebbi/goIAM/internal/auth"
	"github.com/javadmohebbi/goIAM/internal/db"
)

// handleLoginInput represents the expected JSON structure for login.
type handleLoginInput struct {
	Username   string `json:"username"`    // required
	Password   string `json:"password"`    // required
	BackupCode string `json:"backup_code"` // optional
}

// handleLogin returns a Fiber handler that performs user login,
// validates credentials, and returns either a 2FA challenge or a JWT.
func (a *API) handleLogin() fiber.Handler {
	return func(c fiber.Ctx) error {
		var org db.Organization
		if err := db.DB.First(&org).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "default organization not found")
		}

		var body handleLoginInput
		if err := c.Bind().Body(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid input")
		}

		var user db.User
		if err := db.DB.Preload("BackupCodes").Where("username = ? AND organization_id = ?", body.Username, org.ID).First(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "user not found")
		}

		if !auth.CheckPasswordHash(body.Password, user.PasswordHash) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}

		if user.Requires2FA && body.BackupCode == "" {
			totpToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub":  user.ID,
				"name": user.Username,
				"exp":  time.Now().Add(5 * time.Minute).Unix(),
			})
			signed, err := totpToken.SignedString([]byte(a.cfg.JWTSecret))
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "token creation failed")
			}
			return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
				"message": "2FA required",
				"token":   signed,
			})
		}

		if user.Requires2FA && body.BackupCode != "" {
			valid := false
			for _, bc := range user.BackupCodes {
				if !bc.Used && auth.CheckBackupCode(body.BackupCode, bc.CodeHash) {
					valid = true
					bc.Used = true
					db.DB.Save(&bc)
					break
				}
			}
			if !valid {
				return fiber.NewError(fiber.StatusForbidden, "invalid backup code")
			}
		}

		finalToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  user.ID,
			"name": user.Username,
			"exp":  time.Now().Add(24 * time.Hour).Unix(),
		})
		signed, err := finalToken.SignedString([]byte(a.cfg.JWTSecret))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "token creation failed")
		}

		return c.JSON(fiber.Map{"token": signed})
	}
}

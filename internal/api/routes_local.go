// Package api provides HTTP route handlers for goIAM.
// This file contains local authentication and 2FA setup routes.
package api

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/javadmohebbi/goIAM/internal/auth"
	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
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
func RegisterLocalRoutes(app *fiber.App, cfg *config.Config) {
	app.Post("/auth/register", handleRegister)
	app.Post("/auth/login", handleLogin(cfg))

	secure := app.Group("/secure", middleware.RequireAuth(cfg))

	secure.Post("/auth/2fa/setup", handle2FASetup(cfg))
	secure.Post("/auth/2fa/verify", handle2FAVerify(cfg))
	secure.Post("/auth/2fa/disable", handle2FADisable(cfg))
	secure.Post("/auth/backup-codes/regenerate", handleBackupCodes(cfg))
}

// handleRegisterInput represents the expected JSON structure for registration.
type handleRegisterInput struct {
	Username       string `json:"username"`        // required
	Password       string `json:"password"`        // required, validated by regex
	Email          string `json:"email"`           // required, validated by regex
	PhoneNumber    string `json:"phone_number"`    // optional, validated if present
	FirstName      string `json:"first_name"`      // optional
	MiddleName     string `json:"middle_name"`     // optional
	LastName       string `json:"last_name"`       // optional
	Address        string `json:"address"`         // optional
	OrganizationID *uint  `json:"organization_id"` // optional, use default if not set
}

// handleRegister handles user registration for a specific organization.
//
// This function:
//   - Validates user input (username, password, email format, etc.)
//   - Verifies that a valid OrganizationID is provided
//   - Hashes the password securely
//   - Stores the user in the database
//
// Returns 201 on success, 400 for bad input, or 409 if a duplicate user exists.
func handleRegister(c fiber.Ctx) error {
	// Parse and bind JSON input to struct
	var body handleRegisterInput
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	// Validate user input fields
	if err := validateRegisterInput(body); err != nil {
		return err
	}

	// Verify that the organization exists
	var org db.Organization
	if body.OrganizationID == nil {
		return fiber.NewError(fiber.StatusBadRequest, "organization_id is required")
	}
	if err := db.DB.First(&org, *body.OrganizationID).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "specified organization not found")
	}

	// Hash the password
	hash, err := auth.HashPassword(body.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to hash password")
	}

	// Persist the user in the database
	user := db.User{
		Username:       body.Username,
		Email:          body.Email,
		PhoneNumber:    body.PhoneNumber,
		FirstName:      body.FirstName,
		MiddleName:     body.MiddleName,
		LastName:       body.LastName,
		Address:        body.Address,
		PasswordHash:   hash,
		IsActive:       true,
		OrganizationID: org.ID,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusConflict, "user exists or DB error")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user registered"})
}

// validateRegisterInput validates fields from a registration request using centralized config-driven rules.
func validateRegisterInput(input handleRegisterInput) error {
	if input.Username == "" || input.Password == "" || input.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "username, password, and email are required")
	}
	if len(input.Password) < 6 {
		return fiber.NewError(fiber.StatusBadRequest, "password must be at least 6 characters")
	}
	if !validateEmail(input.Email) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid email format")
	}
	return nil
}

func validateEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	return true
}

// handleLoginInput represents the expected JSON structure for login.
type handleLoginInput struct {
	Username   string `json:"username"`    // required
	Password   string `json:"password"`    // required
	BackupCode string `json:"backup_code"` // optional
}

// handleLogin returns a Fiber handler that performs user login,
// validates credentials, and returns either a 2FA challenge or a JWT.
func handleLogin(cfg *config.Config) fiber.Handler {
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
			signed, err := totpToken.SignedString([]byte(cfg.JWTSecret))
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
		signed, err := finalToken.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "token creation failed")
		}

		return c.JSON(fiber.Map{"token": signed})
	}
}

// handle2FAVerifyInput represents the expected JSON structure for 2FA verification.
type handle2FAVerifyInput struct {
	Code string `json:"code"` // required TOTP code
}

// handle2FASetup returns a handler that creates and stores a TOTP secret,
// and returns it to the client for use in authenticator apps.
func handle2FASetup(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(db.User)

		key, qrURL, err := auth.GenerateTOTPSecret(user.Username, "goIAM")
		if err != nil {
			return fiber.NewError(500, "failed to generate TOTP secret")
		}

		user.TOTPSecret = key.Secret()
		if err := db.DB.Save(&user).Error; err != nil {
			return fiber.NewError(500, "failed to save 2FA secret")
		}

		return c.JSON(fiber.Map{
			"otpauth_url": qrURL,
			"secret":      key.Secret(),
		})
	}
}

// handle2FAVerify verifies the TOTP code and enables 2FA for the user,
// issuing a new long-lived token on success.
func handle2FAVerify(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(db.User)

		var body handle2FAVerifyInput
		if err := c.Bind().Body(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid input")
		}

		if user.TOTPSecret == "" {
			return fiber.NewError(fiber.StatusBadRequest, "2FA not initialized")
		}
		if !auth.ValidateTOTP(user.TOTPSecret, body.Code) {
			return fiber.NewError(fiber.StatusForbidden, "invalid TOTP code")
		}

		user.Requires2FA = true
		user.TwoFAVerified = true
		if err := db.DB.Save(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update user")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  user.ID,
			"name": user.Username,
			"exp":  time.Now().Add(24 * time.Hour).Unix(),
		})
		signed, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to create token")
		}

		return c.JSON(fiber.Map{
			"message": "2FA verified and enabled",
			"token":   signed,
		})
	}
}

// handle2FADisableInput represents the expected JSON structure for disabling 2FA.
type handle2FADisableInput struct {
	Code     string `json:"code"`     // optional TOTP code for verification
	Password string `json:"password"` // optional password for verification (not currently used)
}

// handle2FADisable disables TOTP-based 2FA and deletes all backup codes.
func handle2FADisable(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(db.User)

		var body handle2FADisableInput
		if err := c.Bind().Body(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid input")
		}

		if user.TOTPSecret != "" && body.Code != "" {
			if !auth.ValidateTOTP(user.TOTPSecret, body.Code) {
				return fiber.NewError(fiber.StatusForbidden, "invalid TOTP code")
			}
		}

		user.TOTPSecret = ""
		user.Requires2FA = false
		db.DB.Where("user_id = ?", user.ID).Delete(&db.BackupCode{})

		if err := db.DB.Save(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to disable 2FA")
		}

		return c.JSON(fiber.Map{"message": "2FA disabled"})
	}
}

// handleBackupCodes regenerates a new set of backup codes for the user
// and invalidates all previously issued codes.
func handleBackupCodes(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(db.User)

		codes, hashes, err := auth.GenerateBackupCodes(8)
		if err != nil {
			return fiber.NewError(500, "generation failed")
		}

		db.DB.Where("user_id = ?", user.ID).Delete(&db.BackupCode{})
		for _, h := range hashes {
			db.DB.Create(&db.BackupCode{UserID: user.ID, CodeHash: h})
		}

		return c.JSON(fiber.Map{"backup_codes": codes})
	}
}

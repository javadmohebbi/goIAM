package api

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/javadmohebbi/goIAM/internal/auth"
	"github.com/javadmohebbi/goIAM/internal/db"
)

// handle2FAVerifyInput represents the expected JSON structure for 2FA verification.
type handle2FAVerifyInput struct {
	Code string `json:"code"` // required TOTP code
}

// handle2FASetup returns a handler that creates and stores a TOTP secret,
// and returns it to the client for use in authenticator apps.
func (a *API) handle2FASetup() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(db.User)

		key, qrURL, err := auth.GenerateTOTPSecret(user.Username, "goIAM")
		if err != nil {
			return fiber.NewError(500, "failed to generate TOTP secret")
		}

		user.TOTPSecret = key.Secret()
		if err := a.iamDB.Save(&user).Error; err != nil {
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
func (a *API) handle2FAVerify() fiber.Handler {
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
		if err := a.iamDB.Save(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update user")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  user.ID,
			"name": user.Username,
			"exp":  time.Now().Add(24 * time.Hour).Unix(),
		})
		signed, err := token.SignedString([]byte(a.cfg.JWTSecret))
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
func (a *API) handle2FADisable() fiber.Handler {
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
		a.iamDB.Where("user_id = ?", user.ID).Delete(&db.BackupCode{})

		if err := a.iamDB.Save(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to disable 2FA")
		}

		return c.JSON(fiber.Map{"message": "2FA disabled"})
	}
}

// handleBackupCodes regenerates a new set of backup codes for the user
// and invalidates all previously issued codes.
func (a *API) handleBackupCodes() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(db.User)

		codes, hashes, err := auth.GenerateBackupCodes(8)
		if err != nil {
			return fiber.NewError(500, "generation failed")
		}

		a.iamDB.Where("user_id = ?", user.ID).Delete(&db.BackupCode{})
		for _, h := range hashes {
			a.iamDB.Create(&db.BackupCode{UserID: user.ID, CodeHash: h})
		}

		return c.JSON(fiber.Map{"backup_codes": codes})
	}
}

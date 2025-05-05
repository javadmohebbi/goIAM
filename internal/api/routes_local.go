package api

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/javadmohebbi/goIAM/internal/auth"
	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
	"github.com/javadmohebbi/goIAM/internal/middleware"
)

func RegisterLocalRoutes(app *fiber.App, cfg *config.Config) {

	app.Post("/auth/register", func(c fiber.Ctx) error {
		var body struct {
			Username    string `json:"username"`
			Password    string `json:"password"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
			FirstName   string `json:"first_name"`
			MiddleName  string `json:"middle_name"`
			LastName    string `json:"last_name"`
			Address     string `json:"address"`
		}
		if err := c.Bind().Body(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid input")
		}

		hash, err := auth.HashPassword(body.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to hash password")
		}

		user := db.User{
			Username:     body.Username,
			Email:        body.Email,
			PhoneNumber:  body.PhoneNumber,
			FirstName:    body.FirstName,
			MiddleName:   body.MiddleName,
			LastName:     body.LastName,
			Address:      body.Address,
			PasswordHash: hash,
			IsActive:     true,
		}

		if err := db.DB.Create(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusConflict, "user exists or DB error")
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user registered"})
	})

	app.Post("/auth/login", func(c fiber.Ctx) error {
		var body struct {
			Username   string `json:"username"`
			Password   string `json:"password"`
			BackupCode string `json:"backup_code"` // optional
		}
		if err := c.Bind().Body(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid input")
		}

		var user db.User
		if err := db.DB.Preload("BackupCodes").Where("username = ?", body.Username).First(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "user not found")
		}

		// Always check password first
		if !auth.CheckPasswordHash(body.Password, user.PasswordHash) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}

		// ✅ 2FA required but no backup code → return token + status 202
		if user.Requires2FA && body.BackupCode == "" {
			// Return short-lived token to allow /secure/auth/2fa/verify
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

		// ✅ 2FA + backup code handling
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

		// ✅ Issue regular token (24h)
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
	})

	secure := app.Group("/secure", middleware.RequireAuth(cfg))

	secure.Post("/auth/backup-codes/regenerate", func(c fiber.Ctx) error {
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
	})

	secure.Post("/auth/2fa/setup", func(c fiber.Ctx) error {
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
	})

	secure.Post("/auth/2fa/verify", func(c fiber.Ctx) error {
		user := c.Locals("user").(db.User)

		var body struct {
			Code string `json:"code"`
		}
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
		user.TwoFAVerified = true // runtime only
		if err := db.DB.Save(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update user")
		}

		// Generate final JWT token
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

	})

	secure.Post("/auth/2fa/setup", func(c fiber.Ctx) error {
		user := c.Locals("user").(db.User)

		key, qrURL, err := auth.GenerateTOTPSecret(user.Username, "goIAM")
		if err != nil {
			return fiber.NewError(500, "failed to generate TOTP secret")
		}

		// Save to DB
		user.TOTPSecret = key.Secret()
		if err := db.DB.Save(&user).Error; err != nil {
			return fiber.NewError(500, "failed to save 2FA secret")
		}

		return c.JSON(fiber.Map{
			"otpauth_url": qrURL,
			"secret":      key.Secret(), // show only once
		})
	})

	secure.Post("/auth/2fa/disable", func(c fiber.Ctx) error {
		user := c.Locals("user").(db.User)

		var body struct {
			Code     string `json:"code"`     // optional
			Password string `json:"password"` // optional
		}
		if err := c.Bind().Body(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid input")
		}

		// OPTIONAL: require recent password or valid code
		if user.TOTPSecret != "" && body.Code != "" {
			if !auth.ValidateTOTP(user.TOTPSecret, body.Code) {
				return fiber.NewError(fiber.StatusForbidden, "invalid TOTP code")
			}
		}

		// Disable 2FA
		user.TOTPSecret = ""
		user.Requires2FA = false
		db.DB.Where("user_id = ?", user.ID).Delete(&db.BackupCode{})

		if err := db.DB.Save(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to disable 2FA")
		}

		return c.JSON(fiber.Map{"message": "2FA disabled"})
	})

}

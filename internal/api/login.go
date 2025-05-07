package api

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/javadmohebbi/goIAM/internal/auth"
	"github.com/javadmohebbi/goIAM/internal/db"
	"github.com/mssola/user_agent"
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
		if err := a.iamDB.First(&org).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "default organization not found")
		}

		var body handleLoginInput
		if err := c.Bind().Body(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid input")
		}

		var user db.User
		if err := a.iamDB.Preload("BackupCodes").Where("username = ? AND organization_id = ?", body.Username, org.ID).First(&user).Error; err != nil {
			a.storeLoginActivity(c, db.User{Username: body.Username}, "user_not_found")
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credential") // user not found
		}

		if !auth.CheckPasswordHash(body.Password, user.PasswordHash) {
			a.storeLoginActivity(c, user, "invalid_password")
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
					a.iamDB.Save(&bc)
					break
				}
			}
			if !valid {
				a.storeLoginActivity(c, user, "invalid_backup_code")
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

		a.storeLoginActivity(c, user, "success")

		return c.JSON(fiber.Map{"token": signed})
	}
}

// storeLoginActivity creates an audit log record for a login attempt,
// capturing metadata such as user agent, browser, OS, and IP address.
// It logs both successful and failed login attempts, with a provided status label.
func (a *API) storeLoginActivity(c fiber.Ctx, user db.User, status string) {
	ua := user_agent.New(string(c.Request().Header.UserAgent()))
	browser, _ := ua.Browser()

	audit := db.LoginActivity{
		UserID:    user.ID,                                // ID of the user attempting login
		Username:  user.Username,                          // Username of the user attempting login
		IP:        c.IP(),                                 // IP address from which the login was attempted
		UserAgent: string(c.Request().Header.UserAgent()), // Raw User-Agent string
		OS:        ua.OS(),                                // Operating system extracted from User-Agent
		Browser:   browser,                                // Browser name extracted from User-Agent
		Device:    ua.Platform(),                          // Device platform extracted from User-Agent
		Status:    status,                                 // Status of the login attempt (e.g. "success", "invalid_password")
		Success:   status == "success",                    // true if login was successful
	}

	go a.iamDB.Create(&audit)
}

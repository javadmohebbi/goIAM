package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
	"github.com/javadmohebbi/goIAM/internal/smtpclient"
)

// handleResetPasswordRequestLocal send the instruction to user to reset his password
func (a *API) handleResetPasswordRequestLocal(c fiber.Ctx) error {
	type req struct {
		Username string `json:"username"` // required
		Email    string `json:"email"`    // required
	}

	var body req
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	var user db.User

	// if username is provided
	if body.Username != "" {
		if err := a.iamDB.Model(&db.User{}).Where("username = ?", body.Username).First(&user); err != nil {
			// log error
		}

		return sendResetPasswordEmail(user, a.cfg)
	}

	// if email provided
	if body.Email != "" {
		if err := a.iamDB.Model(&db.User{}).Where("email = ?", body.Email).First(&user); err != nil {
			// log error
		}
		return sendResetPasswordEmail(user, a.cfg)
	}

	// anyway, we will show a message to user
	// that if the provided info is correct,
	// you will recieve an email to reset your password
	return nil
}

// sendResetPasswordEmail generates a password reset token, populates an email template
// with the provided user information and configuration, and sends the reset email.
//
// Parameters:
//   - u: the user who requested password reset
//   - cfg: application configuration including SMTP settings
//
// Behavior:
//   - Uses the user's FirstName if available, otherwise falls back to Username
//   - Generates a UUID token (should be stored for later validation — TODO)
//   - Replaces placeholders in the HTML template and sends the email
func sendResetPasswordEmail(u db.User, cfg *config.Config) error {
	// Choose the name to personalize the email
	_name := u.Username
	if u.FirstName != "" {
		_name = u.FirstName
	}

	// Generate a unique reset token (should be stored with expiration — TODO)
	token := uuid.New().String()

	// Prepare template placeholders
	placeholders := map[string]string{
		"Name":    _name,
		"AppName": cfg.AppName,
		"Year":    fmt.Sprintf("%d", time.Now().Year()),
		"Token":   token,
	}

	// Send the reset password email using HTML template
	return smtpclient.SendEmailFromHTMLTemplate(cfg, "Reset Your Password",
		[]string{u.Email}, "templates/reset-password.html", placeholders)
}

package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
	"github.com/javadmohebbi/goIAM/internal/smtpclient"
	"github.com/javadmohebbi/goIAM/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// handleCreateUser allows an authenticated user to create another user within their organization.
// The organization ID is extracted from the authenticated user's context.
func (a *API) handleCreateUser(c fiber.Ctx) error {
	authUser, ok := c.Locals("user").(db.User)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var body struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"-"` // Radnomly generated password for a newly created user
		FirstName   string `json:"first_name"`
		MiddleName  string `json:"middle_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if body.Password == "" {
		body.Password, _ = utils.GenerateRandomString(16)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to hash password")
	}

	user := db.User{
		Username:       body.Username,
		Email:          body.Email,
		PasswordHash:   string(hashedPassword),
		FirstName:      body.FirstName,
		MiddleName:     body.MiddleName,
		LastName:       body.LastName,
		PhoneNumber:    body.PhoneNumber,
		OrganizationID: authUser.OrganizationID,

		// all new users are inactive unless they activate themselves
		// with the provided link in the email or other way of communicateion
		IsActive:      false,
		EmailVerified: false,
		PhoneVerified: false,
	}

	if err := a.iamDB.Create(&user).Error; err != nil {
		var errMsg string
		errMsg = err.Error()
		if strings.Contains(strings.ToUpper(err.Error()), "UNIQUE") {
			if strings.Contains(err.Error(), "username") {
				errMsg = "username already exists"
			}
			if strings.Contains(err.Error(), "email") {
				errMsg = "email already exists"
			}
		}
		return fiber.NewError(
			fiber.StatusInternalServerError,
			fmt.Sprintf(
				"%s. %v",
				"failed to create user",
				errMsg,
			),
		)
	}

	// send email
	_ = sendResetPasswordEmail(user, a.cfg)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_name": user.Username,
		"message":   "User created",
		"created":   true,
	})
}

// sendUserCreationEmail generates an activation token, populates an email template
// with the provided user information and configuration, and sends the account activation email.
//
// Parameters:
//   - u: the user who has just been created and needs to activate their account
//   - cfg: application configuration including SMTP settings
//
// Behavior:
//   - Uses the user's FirstName if available, otherwise falls back to Username
//   - Generates a UUID token (should be stored for later validation — TODO)
//   - Replaces placeholders in the HTML template and sends the email
func sendUserCreationEmail(u db.User, cfg *config.Config) error {
	// Choose the name to personalize the email
	_name := u.Username
	if u.FirstName != "" {
		_name = u.FirstName
	}
	// Generate a unique activation token (should be stored with expiration — TODO)
	token := uuid.New().String()
	// Prepare template placeholders
	placeholders := map[string]string{
		"Name":    _name,
		"AppName": cfg.AppName,
		"Year":    fmt.Sprintf("%d", time.Now().Year()),
		"Token":   token,
	}

	// Send the activation email using HTML template
	return smtpclient.SendEmailFromHTMLTemplate(cfg, "Activate Your Account",
		[]string{u.Email}, "templates/reset-password.html", placeholders)
}

package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/db"
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_name": user.Username,
		"message":   "User created",
		"created":   true,
	})
}

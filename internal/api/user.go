package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/db"
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
		Password    string `json:"password"`
		FirstName   string `json:"first_name"`
		MiddleName  string `json:"middle_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
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
	}

	if err := a.iamDB.Create(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create user")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_id": user.ID,
	})
}

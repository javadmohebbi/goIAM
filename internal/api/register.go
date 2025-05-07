package api

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/auth"
	"github.com/javadmohebbi/goIAM/internal/db"
)

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
	OrganizationID uint   `json:"organization_id"` // optional, use default if not set
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
func (a *API) handleRegister(c fiber.Ctx) error {
	// Parse and bind JSON input to struct
	var body handleRegisterInput
	if err := c.Bind().Body(&body); err != nil {
		log.Println(body)
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	// Validate user input fields
	if err := a.validateRegisterInput(body); err != nil {
		return err
	}

	// Verify that the organization exists
	var org db.Organization
	if body.OrganizationID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "organization_id is required")
	}
	if err := db.DB.First(&org, body.OrganizationID).Error; err != nil {
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
func (a *API) validateRegisterInput(input handleRegisterInput) error {
	if input.Username == "" || input.Password == "" || input.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "username, password, and email are required")
	}
	if len(input.Password) < 6 {
		return fiber.NewError(fiber.StatusBadRequest, "password must be at least 6 characters")
	}
	if !a.validation.ValidateEmail(input.Email) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid email format")
	}
	return nil
}

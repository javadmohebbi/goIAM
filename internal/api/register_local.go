package api

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/javadmohebbi/goIAM/internal/auth"
	"github.com/javadmohebbi/goIAM/internal/db"
	"github.com/javadmohebbi/goIAM/internal/db/seeds"
)

// handleRegisterInput represents the expected JSON structure for registration.
//
// A user may either:
//   - provide an existing organization_id
//   - or leave it blank and supply an optional organization_name and/or slug,
//     in which case a new organization will be created.
type handleRegisterInput struct {
	Username         string `json:"username"`          // required
	Password         string `json:"password"`          // required, validated by regex
	Email            string `json:"email"`             // required, validated by regex
	PhoneNumber      string `json:"phone_number"`      // optional, validated if present
	FirstName        string `json:"first_name"`        // optional
	MiddleName       string `json:"middle_name"`       // optional
	LastName         string `json:"last_name"`         // optional
	Address          string `json:"address"`           // optional
	OrganizationID   uint   `json:"organization_id"`   // optional, use default if not set
	OrganizationName string `json:"organization_name"` // optional: name of new org if org_id not provided
	OrganizationSlug string `json:"organization_slug"` // optional: custom slug (generated from name if not given)
}

// handleRegister handles user registration for a specific organization.
//
// This function:
//   - Validates user input (username, password, email format, etc.)
//   - Verifies that a valid OrganizationID is provided or creates a new organization if not
//   - Hashes the password securely
//   - Stores the user in the database
//
// Returns 201 on success, 400 for bad input, or 409 if a duplicate user exists.
func (a *API) handleRegisterLocal(c fiber.Ctx) error {
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

	var org db.Organization

	// Create new organization always
	orgName := body.OrganizationName
	if orgName == "" {
		suffix := uuid.New().String()[:8]
		orgName = "goIAM Organization " + suffix
	}
	orgSlug := strings.ToLower(strings.ReplaceAll(orgName, " ", "-"))

	// Ensure slug is unique
	var existing db.Organization
	if err := a.iamDB.Where("slug = ?", orgSlug).First(&existing).Error; err == nil {
		orgSlug = orgSlug + "-" + uuid.New().String()[:4]
	}

	org = db.Organization{
		Name:        orgName,
		Slug:        orgSlug,
		Description: "Created automatically during registration",
	}
	if err := a.iamDB.Create(&org).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: organizations.name") {
			return fiber.NewError(fiber.StatusConflict, "organization name already exists")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create organization")
	}

	// Seed default policies into the new organization
	if err := seeds.SeedDefaultPoliciesForOrg(org.ID, a.iamDB); err != nil {
		log.Printf("failed to seed default policies: %v", err)
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

	if err := a.iamDB.Create(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusConflict, "user exists or DB error")
	}

	// Check how many users exist in this organization
	var userCount int64
	if err := a.iamDB.Model(&db.User{}).Where("organization_id = ?", org.ID).Count(&userCount).Error; err == nil {
		if userCount == 1 {
			// First user gets FullAccess
			var fullAccess db.Policy
			if err := a.iamDB.
				Where("slug LIKE ? AND organization_id = ?", "full-access%", org.ID).
				First(&fullAccess).Error; err == nil {
				a.iamDB.Model(&user).Association("Policies").Append(&fullAccess)
			}
		} else {
			// Other users get SelfManage
			var selfManage db.Policy
			if err := a.iamDB.
				Where("slug LIKE ? AND organization_id = ?", "self-manage%", org.ID).
				First(&selfManage).Error; err == nil {
				a.iamDB.Model(&user).Association("Policies").Append(&selfManage)
			}
		}
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

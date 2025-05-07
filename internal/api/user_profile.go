package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/db"
)

// handleGetProfile returns the authenticated user's profile as JSON.
//
// Requires the user to be set in the Fiber context via prior authentication middleware.
// Returns 401 if the user is not present in context.
func (a *API) handleGetProfile(c fiber.Ctx) error {
	user, ok := c.Locals("user").(db.User)
	if !ok {
		return fiber.ErrUnauthorized
	}
	return c.JSON(user)
}

// handleUpdateProfile updates the authenticated user's editable profile fields.
//
// Expects a JSON body with updatable fields. It will not allow changing sensitive fields like
// organization ID, password hash, or 2FA secrets. Returns 400 on bad input or 500 on DB error.
func (a *API) handleUpdateProfile(c fiber.Ctx) error {
	user, ok := c.Locals("user").(db.User)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var updates map[string]interface{}
	if err := c.Bind().Body(&updates); err != nil {
		return fiber.ErrBadRequest
	}

	if err := user.UpdateProfile(a.iamDB, updates); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "update failed")
	}
	return c.SendStatus(fiber.StatusOK)
}

// handleChangePassword sets a new password hash for the authenticated user.
//
// Expects a JSON payload with `new_password_hash` already securely hashed.
// Returns 400 on invalid input or 500 on database error.
func (a *API) handleChangePassword(c fiber.Ctx) error {
	user, ok := c.Locals("user").(db.User)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var payload struct {
		NewHash string `json:"new_password_hash"`
	}
	if err := c.Bind().Body(&payload); err != nil {
		return fiber.ErrBadRequest
	}

	if err := user.UpdatePasswordHash(a.iamDB, payload.NewHash); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "password update failed")
	}
	return c.SendStatus(fiber.StatusOK)
}

// handleEnable2FA enables two-factor authentication for the authenticated user.
//
// Expects a JSON payload with `totp_hash` representing the pre-validated TOTP secret.
// Updates the user record to require 2FA.
func (a *API) handleEnable2FA(c fiber.Ctx) error {
	user, ok := c.Locals("user").(db.User)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var body struct {
		TOTPHash string `json:"totp_hash"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := user.Enable2FA(a.iamDB, body.TOTPHash); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "2FA enable failed")
	}
	return c.SendStatus(fiber.StatusOK)
}

// handleDisable2FA disables two-factor authentication for the authenticated user.
//
// Clears the stored TOTP secret and disables 2FA requirement for the user.
func (a *API) handleDisable2FA(c fiber.Ctx) error {
	user, ok := c.Locals("user").(db.User)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if err := user.Disable2FA(a.iamDB); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "2FA disable failed")
	}
	return c.SendStatus(fiber.StatusOK)
}

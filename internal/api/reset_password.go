package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/db"
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
		return nil
	}

	// if email provided
	if body.Email != "" {
		if err := a.iamDB.Model(&db.User{}).Where("username = ?", body.Username).First(&user); err != nil {
			// log error
		}
		return nil
	}

	// anyway, we will show a message to user
	// that if the provided info is correct,
	// you will recieve an email to reset your password
	return nil
}

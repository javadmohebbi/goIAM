package api

import "github.com/gofiber/fiber/v3"

// handleResetPasswordRequestLocal send the instruction to user to reset his password
func (a *API) handleResetPasswordRequestLocal(c fiber.Ctx) error {
	type req struct {
		Username string `json:"username"` // required
	}

	var body req
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	// return nil
}

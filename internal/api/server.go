package api

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/config"
)

func StartServer(cfg *config.Config) error {
	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("goIAM is healthy ✅")
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	return app.Listen(addr)
}

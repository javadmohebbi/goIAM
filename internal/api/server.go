package api

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/config"
)

func StartServer(cfg *config.Config) error {
	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("healthy")
	})

	if cfg.AuthProvider == "local" {
		RegisterLocalRoutes(app, cfg)
	}

	return app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

package web

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/lord-server/panorama/internal/config"
	"github.com/lord-server/panorama/internal/web/handlers"
)

func Serve(config *config.Config) {
	app := fiber.New()

	app.Static("/", "./static")
	app.Static("/tiles", config.System.TilesPath, fiber.Static{
		MaxAge: 5,
	})

	app.Route("/api/v1", func(router fiber.Router) {
		app.Get("/metadata", handlers.Metadata(config))
	})

	err := app.Listen(config.Web.ListenAddress)
	if err != nil {
		slog.Error("failed to start web server", "err", err)
	}
}

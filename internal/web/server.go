package web

import (
	"log"

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

	app.Get("/metadata.json", handlers.Metadata(config))

	log.Fatal(app.Listen(config.Web.ListenAddress))
}

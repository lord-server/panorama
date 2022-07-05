package web

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/weqqr/panorama/pkg/config"
	"github.com/weqqr/panorama/pkg/web/handlers"
)

func Serve(addr string, config *config.Config) {
	app := fiber.New()

	app.Static("/", "./static")
	app.Static("/tiles", config.TilesPath, fiber.Static{
		MaxAge: 5,
	})

	app.Get("/metadata.json", handlers.Metadata(config))

	log.Fatal(app.Listen(config.ListenAddress))
}

package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/lord-server/panorama/internal/config"
)

func Metadata(config *config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"title":      config.Web.Title,
			"zoomLevels": config.Renderer.ZoomLevels,
		})
	}
}

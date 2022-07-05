package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/weqqr/panorama/pkg/config"
)

func Metadata(config *config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"title":      config.Title,
			"zoomLevels": config.ZoomLevels,
		})
	}
}

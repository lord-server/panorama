package render

import (
	"github.com/lord-server/panorama/internal/game"
	"github.com/lord-server/panorama/internal/raster"
	"github.com/lord-server/panorama/internal/spatial"
	"github.com/lord-server/panorama/internal/world"
)

type TilePosition struct {
	X, Y int
}

type Renderer interface {
	RenderTile(pos TilePosition, w *world.World, game *game.Game) *raster.RenderBuffer
	ProjectRegion(region spatial.Region) spatial.ProjectedRegion
	// ListTilesWithBlock(x, y, z int) []TilePosition
	// ListTilesInsideRegion(region config.Region) []TilePosition
}

package render

import (
	"github.com/lord-server/panorama/pkg/game"
	"github.com/lord-server/panorama/pkg/raster"
	"github.com/lord-server/panorama/pkg/spatial"
	"github.com/lord-server/panorama/pkg/world"
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

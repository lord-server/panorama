package render

import (
	"github.com/weqqr/panorama/pkg/game"
	"github.com/weqqr/panorama/pkg/raster"
	"github.com/weqqr/panorama/pkg/spatial"
	"github.com/weqqr/panorama/pkg/world"
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

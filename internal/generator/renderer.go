package generator

import (
	"github.com/lord-server/panorama/internal/game"
	"github.com/lord-server/panorama/internal/generator/rasterizer"
	"github.com/lord-server/panorama/internal/world"
	"github.com/lord-server/panorama/pkg/geom"
)

type TilePosition struct {
	X, Y int
}

type Renderer interface {
	RenderTile(pos TilePosition, w *world.World, game *game.Game) *rasterizer.RenderBuffer
	ProjectRegion(region geom.Region) geom.ProjectedRegion
	// ListTilesWithBlock(x, y, z int) []TilePosition
	// ListTilesInsideRegion(region config.Region) []TilePosition
}

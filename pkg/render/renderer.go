package render

import (
	"image"

	"github.com/weqqr/panorama/pkg/game"
	"github.com/weqqr/panorama/pkg/world"
)

type TilePosition struct {
	X, Y int
}

type Renderer interface {
	RenderTile(pos TilePosition, w *world.World, game *game.Game) *image.NRGBA
	ListTilesWithBlock(x, y, z int) []TilePosition
}

package render

import (
	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/world"
	"image"
)

type TilePosition struct {
	X, Y int
}

type Renderer interface {
	RenderTile(pos TilePosition, w *world.World, game *game.Game) *image.NRGBA
	ListTilesWithBlock(x, y, z int) []TilePosition
}

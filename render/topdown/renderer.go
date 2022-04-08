package topdown

import (
	"image"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/lm"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/world"
)

const BaseResolution = 16
const TileSize = BaseResolution * world.MapBlockSize

type Renderer struct {
	nr render.NodeRasterizer

	lowerLimit int
	upperLimit int
}

func NewRenderer(lowerLimit, upperLimit int) Renderer {
	return Renderer{
		nr: render.NewNodeRasterizer(BaseResolution, BaseResolution, BaseResolution, lm.TopDownProjection()),

		lowerLimit: lowerLimit,
		upperLimit: upperLimit,
	}
}

func (r *Renderer) RenderTile(tilePos render.TilePosition, w *world.World, game *game.Game) *image.NRGBA {
	rect := image.Rect(0, 0, TileSize, TileSize)
	target := raster.NewRenderBuffer(rect)

	blockX := tilePos.X
	blockZ := tilePos.Y

	for y := r.lowerLimit; y < r.upperLimit; y++ {
		block, err := w.GetBlock(blockX, y, blockZ)
		if err != nil || block == nil {
			continue
		}

		depthOffset := float32(-y * world.MapBlockSize)
		renderBlock(target, &r.nr, block, game, depthOffset)
	}

	return target.Color
}

func (r *Renderer) ListTilesWithBlock(x, y, z int) []render.TilePosition {
	return []render.TilePosition{}
}

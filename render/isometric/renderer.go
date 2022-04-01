package isometric

import (
	"image"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/world"
)

type Renderer struct {
	nr NodeRasterizer

	lowerLimit int
	upperLimit int
}

func NewRenderer(lowerLimit, upperLimit int) Renderer {
	return Renderer{
		nr: NewNodeRasterizer(),

		lowerLimit: lowerLimit,
		upperLimit: upperLimit,
	}
}

func (r *Renderer) RenderTile(tilePos render.TilePosition, w *world.World, game *game.Game) *image.NRGBA {
	tilePos.Y *= 2

	rect := image.Rect(0, 0, TileBlockWidth, TileBlockWidth)
	target := raster.NewRenderBuffer(rect)

	centerX := tilePos.Y - tilePos.X
	centerY := 0
	centerZ := tilePos.Y + tilePos.X

	for i := r.lowerLimit; i < r.upperLimit; i++ {
		for z := -2; z <= 2; z++ {
			for x := -2; x <= 2; x++ {
				blockX := centerX + x + i
				blockY := centerY + i
				blockZ := centerZ + z + i

				block, err := w.GetBlock(blockX, blockY, blockZ)
				if err != nil {
					panic(err)
				}

				if block == nil {
					continue
				}

				tileOffsetX := world.MapBlockSize * BaseResolution / 2 * (z - x)
				tileOffsetY := world.MapBlockSize * BaseResolution / 4 * (z + x)

				depthOffset := -float32(z+x)/math.Sqrt2*world.MapBlockSize - 2*float32(i)*math.Sqrt2*world.MapBlockSize
				renderBlock(target, &r.nr, block, game, tileOffsetX, tileOffsetY, depthOffset)
			}
		}
	}

	return target.Color
}

func (r *Renderer) ListTilesWithBlock(x, y, z int) []render.TilePosition {
	return []render.TilePosition{}
}

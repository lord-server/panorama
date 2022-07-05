package isometric

import (
	"image"
	"math"

	"github.com/weqqr/panorama/pkg/game"
	"github.com/weqqr/panorama/pkg/raster"
	"github.com/weqqr/panorama/pkg/render"
	"github.com/weqqr/panorama/pkg/world"
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
		for z := -3; z <= 3; z++ {
			for x := -3; x <= 3; x++ {
				blockX := centerX + x + i
				blockY := centerY + i
				blockZ := centerZ + z + i

				neighborhood := BlockNeighborhood{}

				neighborhood.FetchBlock(1, 1, 1, blockX, blockY, blockZ, w)
				// neighborhood.FetchBlock(0, 1, 1, blockX-1, blockY, blockZ, w)
				neighborhood.FetchBlock(2, 1, 1, blockX+1, blockY, blockZ, w)
				// neighborhood.FetchBlock(1, 0, 1, blockX, blockY-1, blockZ, w)
				neighborhood.FetchBlock(1, 2, 1, blockX, blockY+1, blockZ, w)
				// neighborhood.FetchBlock(1, 1, 0, blockX, blockY, blockZ-1, w)
				neighborhood.FetchBlock(1, 1, 2, blockX, blockY, blockZ+1, w)

				tileOffsetX := BaseResolution / 2 * (z - x) * world.MapBlockSize
				tileOffsetY := (BaseResolution/4*(z+x+2*i) - i*YOffsetCoef) * world.MapBlockSize

				depthOffset := (-float32(z+x+2*i)/math.Sqrt2 - 0.5*float32(i)) * world.MapBlockSize
				renderBlock(target, &r.nr, &neighborhood, game, tileOffsetX, tileOffsetY, depthOffset)
			}
		}
	}

	return target.Color
}

func (r *Renderer) ListTilesWithBlock(x, y, z int) []render.TilePosition {
	return []render.TilePosition{}
}

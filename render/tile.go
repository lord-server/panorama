package render

import (
	"image"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/world"
)

func RenderTile(world *world.World, game *game.Game) *image.NRGBA {
	rect := image.Rect(0, 0, TileBlockWidth, TileBlockHeight)
	tile := image.NewNRGBA(rect)
	depth := NewDepthBuffer(rect)

	centerX := 5
	centerY := 0
	centerZ := -6

	originX, originY := 0, 0

	nr := NewNodeRasterizer()

	for z := -1; z <= 1; z++ {
		for y := -1; y <= 1; y++ {
			for x := -1; x <= 1; x++ {
				block, err := world.GetBlock(centerX+x, centerY+y, centerZ+z)
				if err != nil {
					continue
				}

				blockColor, blockDepth := RenderBlock(&nr, block, game)

				tileOffsetX := originX + 16*BaseResolution*(z-x)/2
				tileOffsetY := originY + 16*BaseResolution/4*(z+x) - 16*YOffsetCoef*y

				depthOffset := (-float32(z+x)/math.Sqrt2 - 0.5*(float32(y))) * 16
				overlayWithDepth(tile, depth, blockColor, blockDepth, image.Pt(tileOffsetX, tileOffsetY), depthOffset)
			}
		}
	}

	return tile
}

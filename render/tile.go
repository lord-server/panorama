package render

import (
	"image"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/imaging"
	"github.com/weqqr/panorama/world"
)

func RenderTile(tileX, tileY int, nr *NodeRasterizer, w *world.World, game *game.Game) *image.NRGBA {
	rect := image.Rect(0, 0, TileBlockWidth, TileBlockWidth)
	target := imaging.NewRenderBuffer(rect)

	centerX := tileY - tileX
	centerY := 0
	centerZ := tileY + tileX

	upperLimit := 5
	lowerLimit := -5

	for i := lowerLimit; i < upperLimit; i++ {
		for z := -2; z <= 2; z++ {
			for y := 0; y <= 0; y++ {
				for x := -2; x <= 2; x++ {
					blockX := centerX + x + i
					blockY := centerY + y + i
					blockZ := centerZ + z + i

					block, err := w.GetBlock(blockX, blockY, blockZ)
					if err != nil {
						panic(err)
					}

					if block == nil {
						continue
					}

					tileOffsetX := world.MapBlockSize * BaseResolution * (z - x) / 2
					tileOffsetY := world.MapBlockSize*BaseResolution/4*(z+x) - world.MapBlockSize*YOffsetCoef*y

					depthOffset := -float32(z+x)/math.Sqrt2*world.MapBlockSize - 0.5*float32(y)*world.MapBlockSize - 2*float32(i)*math.Sqrt2*world.MapBlockSize
					RenderBlock(target, nr, block, game, tileOffsetX, tileOffsetY, depthOffset)
				}
			}
		}
	}

	return target.Color
}

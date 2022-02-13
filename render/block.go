package render

import (
	"fmt"
	"image"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/world"
)

func overlayWithDepth(target *image.NRGBA, targetDepth *DepthBuffer, source *image.NRGBA, sourceDepth *DepthBuffer, origin image.Point, depthOffset float32) {
	if source == nil {
		return
	}

	width := source.Rect.Dx()
	height := source.Rect.Dy()

	for y := origin.Y; y < origin.Y+height; y++ {
		for x := origin.X; x < origin.X+width; x++ {
			if x > 1000 || y > 1000 {
				fmt.Printf("x=%v y=%v origin=%v\n w=%v, h=%v", x, y, origin, width, height)
			}
			targetZ := targetDepth.At(x, y)
			sourceZ := sourceDepth.At(x-origin.X, y-origin.Y) + depthOffset

			if sourceZ > targetZ {
				continue
			}

			targetDepth.Set(x, y, sourceZ)

			c := source.NRGBAAt(x-origin.X, y-origin.Y)
			if c.A == 0 {
				// TODO: support opacity
				continue
			}
			target.SetNRGBA(x, y, c)
		}
	}
}

func RenderBlock(target *image.NRGBA, targetDepth *DepthBuffer, nr *NodeRasterizer, block *world.MapBlock, game *game.Game, offsetX, offsetY int, depth float32) {
	rect := image.Rect(0, 0, TileBlockWidth, TileBlockHeight)

	// FIXME: nodes must define their origin points
	originX, originY := rect.Dx()/2-BaseResolution/2, rect.Dy()/2+BaseResolution/4+2

	for z := world.MapBlockSize - 1; z >= 0; z-- {
		for y := world.MapBlockSize - 1; y >= 0; y-- {
			for x := world.MapBlockSize - 1; x >= 0; x-- {
				node := block.GetNode(x, y, z)
				nodeName := block.ResolveName(node.ID)
				gameNode := game.Node(nodeName)

				nodeColor, nodeDepth := nr.Render(nodeName, &gameNode)

				tileOffsetX := originX + BaseResolution*(z-x)/2 + offsetX
				tileOffsetY := originY + BaseResolution/4*(z+x) - YOffsetCoef*y + offsetY

				depthOffset := -float32(z+x)/math.Sqrt2 - 0.5*(float32(y)) + depth
				overlayWithDepth(target, targetDepth, nodeColor, nodeDepth, image.Pt(tileOffsetX, tileOffsetY), depthOffset)
			}
		}
	}
}

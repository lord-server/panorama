package render

import (
	"image"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/world"
)

func Block(target *raster.RenderBuffer, nr *NodeRasterizer, block *world.MapBlock, game *game.Game, offsetX, offsetY int, depth float32) {
	rect := image.Rect(0, 0, TileBlockWidth, TileBlockHeight)

	// FIXME: nodes must define their origin points
	originX, originY := rect.Dx()/2-BaseResolution/2, rect.Dy()/2+BaseResolution/4+2

	for z := world.MapBlockSize - 1; z >= 0; z-- {
		for y := world.MapBlockSize - 1; y >= 0; y-- {
			for x := world.MapBlockSize - 1; x >= 0; x-- {
				node := block.GetNode(x, y, z)
				nodeName := block.ResolveName(node.ID)
				gameNode := game.Node(nodeName)

				renderedNode := nr.Render(nodeName, &gameNode)

				tileOffsetX := originX + BaseResolution*(z-x)/2 + offsetX
				tileOffsetY := originY + BaseResolution/4*(z+x) - YOffsetCoef*y + offsetY

				depthOffset := -float32(z+x)/math.Sqrt2 - 0.5*(float32(y)) + depth
				target.OverlayDepthAware(renderedNode, image.Pt(tileOffsetX, tileOffsetY), depthOffset)
			}
		}
	}
}

package topdown

import (
	"image"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/world"
)

func renderBlock(target *raster.RenderBuffer, nr *render.NodeRasterizer, block *world.MapBlock, g *game.Game, depth float32) {
	// FIXME: nodes must define their origin points
	for z := world.MapBlockSize - 1; z >= 0; z-- {
		for y := world.MapBlockSize - 1; y >= 0; y-- {
			for x := world.MapBlockSize - 1; x >= 0; x-- {
				tileOffsetX := x * BaseResolution
				tileOffsetY := z * BaseResolution

				node := block.GetNode(x, y, z)
				nodeName := block.ResolveName(node.ID)

				nodeDef := g.NodeDef(nodeName)

				renderableNode := render.RenderableNode{
					Name:   nodeName,
					Light:  1.0,
					Param2: node.Param2,
				}
				renderedNode := nr.Render(renderableNode, &nodeDef)

				depthOffset := -float32(y)
				target.OverlayDepthAware(renderedNode, image.Pt(tileOffsetX, tileOffsetY), depthOffset)
			}
		}
	}
}

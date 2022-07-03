package isometric

import (
	"image"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/world"
)

func renderBlock(target *raster.RenderBuffer, nr *render.NodeRasterizer, neighborhood *render.BlockNeighborhood, g *game.Game, offsetX, offsetY int, depth float32) {
	rect := image.Rect(0, 0, render.BaseTileSize, render.TileBlockHeight)

	// FIXME: nodes must define their origin points
	originX, originY := rect.Dx()/2-render.BaseResolution/2, rect.Dy()/2+render.BaseResolution/4+2

	for z := world.MapBlockSize - 1; z >= 0; z-- {
		for y := world.MapBlockSize - 1; y >= 0; y-- {
			for x := world.MapBlockSize - 1; x >= 0; x-- {
				tileOffsetX := originX + render.BaseResolution*(z-x)/2 + offsetX
				tileOffsetY := originY + render.BaseResolution/4*(z+x) - render.YOffset*y + offsetY

				// Fast path: Don't bother with nodes outside viewport
				nodeTileTooLow := tileOffsetX <= target.Color.Rect.Min.X-render.BaseResolution || tileOffsetY <= target.Color.Rect.Min.Y-render.BaseResolution-render.BaseResolution/8
				nodeTileTooHigh := tileOffsetX >= target.Color.Rect.Max.X || tileOffsetY >= target.Color.Rect.Max.Y

				if nodeTileTooLow || nodeTileTooHigh {
					continue
				}

				name, param1, param2 := neighborhood.GetNode(x, y, z)

				// Fast path: checking for air immediately is faster than fetching NodeDefinition
				if name == "air" {
					continue
				}

				nodeDef := g.NodeDef(name)

				light := render.DecodeLight(param1)
				if l := render.DecodeLight(neighborhood.GetParam1(x+1, y, z)); l > light {
					light = l
				}
				if l := render.DecodeLight(neighborhood.GetParam1(x, y+1, z)); l > light {
					light = l
				}
				if l := render.DecodeLight(neighborhood.GetParam1(x, y, z+1)); l > light {
					light = l
				}

				renderableNode := render.RenderableNode{
					Name:   name,
					Light:  light,
					Param2: param2,
				}
				renderedNode := nr.Render(renderableNode, &nodeDef)

				depthOffset := -float32(z+x)/math.Sqrt2 - 0.5*(float32(y)) + depth
				target.OverlayDepthAware(renderedNode, image.Pt(tileOffsetX, tileOffsetY), depthOffset)
			}
		}
	}
}

package isometric

import (
	"image"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/lm"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/world"
)

type Renderer struct {
	nr render.NodeRasterizer

	lowerLimit int
	upperLimit int
}

func NewRenderer(lowerLimit, upperLimit int) Renderer {
	return Renderer{
		nr: render.NewNodeRasterizer(lm.DimetricProjection()),

		lowerLimit: lowerLimit,
		upperLimit: upperLimit,
	}
}

func (r *Renderer) renderBlock(target *raster.RenderBuffer, neighborhood *render.BlockNeighborhood, game *game.Game, offsetX, offsetY int, depth float32) {
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

				nodeDef := game.NodeDef(name)

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
				renderedNode := r.nr.Render(renderableNode, &nodeDef)

				depthOffset := -float32(z+x)/math.Sqrt2 - 0.5*(float32(y)) + depth
				target.OverlayDepthAware(renderedNode, image.Pt(tileOffsetX, tileOffsetY), depthOffset)
			}
		}
	}
}

func (r *Renderer) RenderTile(tilePos render.TilePosition, w *world.World, game *game.Game) *image.NRGBA {
	tilePos.Y *= 2

	rect := image.Rect(0, 0, render.BaseTileSize, render.BaseTileSize)
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

				neighborhood := render.BlockNeighborhood{}
				neighborhood.FetchBlock(1, 1, 1, blockX, blockY, blockZ, w)
				neighborhood.FetchBlock(2, 1, 1, blockX+1, blockY, blockZ, w)
				neighborhood.FetchBlock(1, 2, 1, blockX, blockY+1, blockZ, w)
				neighborhood.FetchBlock(1, 1, 2, blockX, blockY, blockZ+1, w)

				tileOffsetX := render.BaseResolution / 2 * (z - x) * world.MapBlockSize
				tileOffsetY := (render.BaseResolution/4*(z+x+2*i) - i*render.YOffset) * world.MapBlockSize

				depthOffset := (-float32(z+x+2*i)/math.Sqrt2 - 0.5*float32(i)) * world.MapBlockSize
				r.renderBlock(target, &neighborhood, game, tileOffsetX, tileOffsetY, depthOffset)
			}
		}
	}

	return target.Color
}

func (r *Renderer) ListTilesWithBlock(x, y, z int) []render.TilePosition {
	return []render.TilePosition{}
}

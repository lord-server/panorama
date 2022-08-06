package isometric

import (
	"image"
	"math"

	"github.com/weqqr/panorama/pkg/game"
	"github.com/weqqr/panorama/pkg/raster"
	"github.com/weqqr/panorama/pkg/region"
	"github.com/weqqr/panorama/pkg/render"
	"github.com/weqqr/panorama/pkg/world"
)

type Renderer struct {
	nr NodeRasterizer

	region region.Region
	game   *game.Game
}

func NewRenderer(region region.Region, game *game.Game) *Renderer {
	return &Renderer{
		nr:     NewNodeRasterizer(),
		region: region,
		game:   game,
	}
}

func (r *Renderer) renderBlock(target *raster.RenderBuffer, neighborhood *BlockNeighborhood, offsetX, offsetY int, depth float32) {
	rect := image.Rect(0, 0, TileBlockWidth, TileBlockHeight)

	// FIXME: nodes must define their origin points
	originX, originY := rect.Dx()/2-BaseResolution/2, rect.Dy()/2+BaseResolution/4+2

	for z := world.MapBlockSize - 1; z >= 0; z-- {
		for y := world.MapBlockSize - 1; y >= 0; y-- {
			for x := world.MapBlockSize - 1; x >= 0; x-- {
				tileOffsetX := originX + BaseResolution*(z-x)/2 + offsetX
				tileOffsetY := originY + BaseResolution/4*(z+x) - YOffsetCoef*y + offsetY

				// Fast path: Don't bother with nodes outside viewport
				nodeTileTooLow := tileOffsetX <= target.Color.Rect.Min.X-BaseResolution || tileOffsetY <= target.Color.Rect.Min.Y-BaseResolution-BaseResolution/8
				nodeTileTooHigh := tileOffsetX >= target.Color.Rect.Max.X || tileOffsetY >= target.Color.Rect.Max.Y

				if nodeTileTooLow || nodeTileTooHigh {
					continue
				}

				name, param1, param2 := neighborhood.GetNode(x, y, z)

				// Fast path: checking for air immediately is faster than fetching NodeDefinition
				if name == "air" {
					continue
				}

				nodeDef := r.game.NodeDef(name)

				light := decodeLight(param1)
				if l := decodeLight(neighborhood.GetParam1(x+1, y, z)); l > light {
					light = l
				}
				if l := decodeLight(neighborhood.GetParam1(x, y+1, z)); l > light {
					light = l
				}
				if l := decodeLight(neighborhood.GetParam1(x, y, z+1)); l > light {
					light = l
				}

				renderableNode := RenderableNode{
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

func (r *Renderer) RenderTile(tilePos render.TilePosition, w *world.World, game *game.Game) *raster.RenderBuffer {
	tilePos.Y *= 2

	rect := image.Rect(0, 0, TileBlockWidth, TileBlockWidth)
	target := raster.NewRenderBuffer(rect)

	centerX := tilePos.Y - tilePos.X
	centerY := 0
	centerZ := tilePos.Y + tilePos.X

	yMin := int(math.Floor(float64(r.region.YBounds.Min) / float64(world.MapBlockSize)))
	yMax := int(math.Ceil(float64(r.region.YBounds.Max) / float64(world.MapBlockSize)))

	for i := yMin; i < yMax; i++ {
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
				r.renderBlock(target, &neighborhood, tileOffsetX, tileOffsetY, depthOffset)
			}
		}
	}

	return target
}

func ProjectRegion(r region.Region) region.TileRegion {
	xMin := int(math.Floor(float64((r.ZBounds.Min - r.XBounds.Max)) / 2 / world.MapBlockSize))
	xMax := int(math.Ceil(float64((r.ZBounds.Max - r.XBounds.Min)) / 2 / world.MapBlockSize))

	yMin := int(math.Floor((float64(r.ZBounds.Min+r.XBounds.Min+2*r.YBounds.Max)/4 - float64(r.YBounds.Max*YOffsetCoef)/BaseResolution) / world.MapBlockSize))
	yMax := int(math.Ceil((float64(r.ZBounds.Max+r.XBounds.Max+2*r.YBounds.Min)/4 - float64(r.YBounds.Min*YOffsetCoef)/BaseResolution) / world.MapBlockSize))

	return region.TileRegion{
		XBounds: region.Bounds{
			Min: xMin,
			Max: xMax,
		},
		YBounds: region.Bounds{
			Min: yMin,
			Max: yMax,
		},
	}
}

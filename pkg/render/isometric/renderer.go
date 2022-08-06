package isometric

import (
	"image"
	"math"

	"github.com/weqqr/panorama/pkg/game"
	"github.com/weqqr/panorama/pkg/raster"
	"github.com/weqqr/panorama/pkg/render"
	"github.com/weqqr/panorama/pkg/spatial"
	"github.com/weqqr/panorama/pkg/world"
)

type Renderer struct {
	nr NodeRasterizer

	region spatial.Region
	game   *game.Game
}

func NewRenderer(region spatial.Region, game *game.Game) *Renderer {
	return &Renderer{
		nr:     NewNodeRasterizer(),
		region: region,
		game:   game,
	}
}

func (r *Renderer) renderBlock(target *raster.RenderBuffer, blockPos spatial.BlockPos, neighborhood *BlockNeighborhood, offsetX, offsetY int, depth float32) {
	rect := image.Rect(0, 0, TileBlockWidth, TileBlockHeight)

	// FIXME: nodes must define their origin points
	originX, originY := rect.Dx()/2-BaseResolution/2, rect.Dy()/2+BaseResolution/4+2

	for z := world.MapBlockSize - 1; z >= 0; z-- {
		for y := world.MapBlockSize - 1; y >= 0; y-- {
			for x := world.MapBlockSize - 1; x >= 0; x-- {
				nodePos := spatial.NodePos{X: x, Y: y, Z: z}

				tileOffsetX := originX + BaseResolution*(z-x)/2 + offsetX
				tileOffsetY := originY + BaseResolution*(z+x)/4 + offsetY - YOffsetCoef*y

				// Fast path: Don't bother with nodes outside viewport
				nodeTileTooLow := tileOffsetX <= target.Color.Rect.Min.X-BaseResolution || tileOffsetY <= target.Color.Rect.Min.Y-BaseResolution-BaseResolution/8
				nodeTileTooHigh := tileOffsetX >= target.Color.Rect.Max.X || tileOffsetY >= target.Color.Rect.Max.Y

				if nodeTileTooLow || nodeTileTooHigh {
					continue
				}

				name, param1, param2 := neighborhood.GetNode(nodePos)

				// Fast path: checking for air immediately is faster than fetching NodeDefinition
				if name == "air" {
					continue
				}

				nodeDef := r.game.NodeDef(name)

				lightOffsets := []spatial.NodePos{
					{X: 1, Y: 0, Z: 0},
					{X: 0, Y: 1, Z: 0},
					{X: 0, Y: 0, Z: 1},
				}

				light := decodeLight(param1)
				for _, offset := range lightOffsets {
					pos := nodePos.Add(offset)
					if l := decodeLight(neighborhood.GetParam1(pos)); l > light {
						light = l
					}
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
				blockPos := spatial.BlockPos{
					X: centerX + x + i,
					Y: centerY + i,
					Z: centerZ + z + i,
				}

				neighborhood := BlockNeighborhood{}

				neighborhood.FetchBlock(w, spatial.BlockPos{0, 0, 0}, blockPos)
				neighborhood.FetchBlock(w, spatial.BlockPos{1, 0, 0}, blockPos)
				neighborhood.FetchBlock(w, spatial.BlockPos{0, 1, 0}, blockPos)
				neighborhood.FetchBlock(w, spatial.BlockPos{0, 0, 1}, blockPos)

				tileOffsetX := BaseResolution / 2 * (z - x) * world.MapBlockSize
				tileOffsetY := (BaseResolution/4*(z+x+2*i) - i*YOffsetCoef) * world.MapBlockSize

				depthOffset := (-float32(z+x+2*i)/math.Sqrt2 - 0.5*float32(i)) * world.MapBlockSize
				r.renderBlock(target, blockPos, &neighborhood, tileOffsetX, tileOffsetY, depthOffset)
			}
		}
	}

	return target
}

func ProjectRegion(region spatial.Region) spatial.TileRegion {
	xMin := int(math.Floor(float64((region.ZBounds.Min - region.XBounds.Max)) / 2 / world.MapBlockSize))
	xMax := int(math.Ceil(float64((region.ZBounds.Max - region.XBounds.Min)) / 2 / world.MapBlockSize))

	yMin := int(math.Floor((float64(region.ZBounds.Min+region.XBounds.Min+2*region.YBounds.Max)/4 - float64(region.YBounds.Max*YOffsetCoef)/BaseResolution) / world.MapBlockSize))
	yMax := int(math.Ceil((float64(region.ZBounds.Max+region.XBounds.Max+2*region.YBounds.Min)/4 - float64(region.YBounds.Min*YOffsetCoef)/BaseResolution) / world.MapBlockSize))

	return spatial.TileRegion{
		XBounds: spatial.Bounds{
			Min: xMin,
			Max: xMax,
		},
		YBounds: spatial.Bounds{
			Min: yMin,
			Max: yMax,
		},
	}
}

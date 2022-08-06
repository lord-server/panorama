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

const Epsilon = 0.0001

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

func (r *Renderer) renderNode(
	target *raster.RenderBuffer,
	pos spatial.NodePos,
	neighborhood *render.BlockNeighborhood,
	offset image.Point,
	depthOffset float32,
) {
	name, param1, param2 := neighborhood.GetNode(pos)

	// Fast path: checking for air immediately is faster than fetching NodeDefinition
	if name == "air" {
		return
	}

	nodeDef := r.game.NodeDef(name)

	// Estimate lighting by sampling neighboring nodes and using the brightest one
	neighborOffsets := []spatial.NodePos{
		{X: 1, Y: 0, Z: 0},
		{X: 0, Y: 1, Z: 0},
		{X: 0, Y: 0, Z: 1},
	}

	maxParam1 := param1
	for _, offset := range neighborOffsets {
		neighborPos := pos.Add(offset)
		if param1 := neighborhood.GetParam1(neighborPos); param1 > maxParam1 {
			maxParam1 = param1
		}
	}

	// Make underground edges visible (otherwise the edge becomes oddly thin and
	// that doesn't look good)
	if maxParam1 == render.ZeroIntensity {
		maxParam1 = render.MapEdgeIntensity
	}

	renderableNode := RenderableNode{
		Name:   name,
		Light:  render.DecodeLight(maxParam1),
		Param2: param2,
	}
	renderedNode := r.nr.Render(renderableNode, &nodeDef)

	depthOffset = -float32(pos.Z+pos.X)/math.Sqrt2 - 0.5*(float32(pos.Y)) + depthOffset
	target.OverlayDepthAware(renderedNode, offset, depthOffset)
}

func (r *Renderer) renderBlock(
	target *raster.RenderBuffer,
	blockPos spatial.BlockPos,
	neighborhood *render.BlockNeighborhood,
	offset image.Point,
	depthOffset float32,
) {
	rect := image.Rect(0, 0, TileBlockWidth, TileBlockHeight)

	// FIXME: nodes must define their origin points
	originX, originY := rect.Dx()/2-BaseResolution/2, rect.Dy()/2+BaseResolution/4+2

	for z := spatial.BlockSize - 1; z >= 0; z-- {
		for y := spatial.BlockSize - 1; y >= 0; y-- {
			for x := spatial.BlockSize - 1; x >= 0; x-- {
				nodePos := spatial.NodePos{X: x, Y: y, Z: z}
				nodeWorldPos := blockPos.AddNode(nodePos)

				if !r.region.Intersects(nodeWorldPos.Region()) {
					continue
				}

				offset := image.Point{
					X: originX + BaseResolution*(z-x)/2 + offset.X,
					Y: originY + BaseResolution*(z+x)/4 + offset.Y - YOffsetCoef*y,
				}

				r.renderNode(target, nodePos, neighborhood, offset, depthOffset)
			}
		}
	}
}

func (r *Renderer) RenderTile(
	tilePos render.TilePosition,
	world *world.World,
	game *game.Game,
) *raster.RenderBuffer {
	tilePos.Y *= 2

	rect := image.Rect(0, 0, TileBlockWidth, TileBlockWidth)
	target := raster.NewRenderBuffer(rect)

	centerX := tilePos.Y - tilePos.X
	centerY := 0
	centerZ := tilePos.Y + tilePos.X

	yMin := int(math.Floor(float64(r.region.YBounds.Min) / float64(spatial.BlockSize)))
	yMax := int(math.Ceil(float64(r.region.YBounds.Max) / float64(spatial.BlockSize)))

	for i := yMin; i < yMax; i++ {
		for z := -3; z <= 3; z++ {
			for x := -3; x <= 3; x++ {
				blockPos := spatial.BlockPos{
					X: centerX + x + i,
					Y: centerY + i,
					Z: centerZ + z + i,
				}

				neighborhood := render.BlockNeighborhood{}

				neighborhood.FetchBlock(world, spatial.BlockPos{X: 0, Y: 0, Z: 0}, blockPos)
				neighborhood.FetchBlock(world, spatial.BlockPos{X: 1, Y: 0, Z: 0}, blockPos)
				neighborhood.FetchBlock(world, spatial.BlockPos{X: 0, Y: 1, Z: 0}, blockPos)
				neighborhood.FetchBlock(world, spatial.BlockPos{X: 0, Y: 0, Z: 1}, blockPos)

				offset := image.Point{
					X: BaseResolution * (z - x) / 2 * spatial.BlockSize,
					Y: (BaseResolution*(z+x+2*i)/4 - i*YOffsetCoef) * spatial.BlockSize,
				}

				depthOffset := (-float32(z+x+2*i)/math.Sqrt2 - 0.5*float32(i)) * spatial.BlockSize
				r.renderBlock(target, blockPos, &neighborhood, offset, depthOffset)
			}
		}
	}

	return target
}

func ProjectRegion(region spatial.Region) spatial.TileRegion {
	xMin := int(math.Floor(float64((region.ZBounds.Min - region.XBounds.Max)) / 2 / spatial.BlockSize))
	xMax := int(math.Ceil(float64((region.ZBounds.Max - region.XBounds.Min)) / 2 / spatial.BlockSize))

	yMin := int(math.Floor((float64(region.ZBounds.Min+region.XBounds.Min+2*region.YBounds.Max)/4 - float64(region.YBounds.Max*YOffsetCoef)/BaseResolution) / spatial.BlockSize))
	yMax := int(math.Ceil((float64(region.ZBounds.Max+region.XBounds.Max+2*region.YBounds.Min)/4 - float64(region.YBounds.Min*YOffsetCoef)/BaseResolution) / spatial.BlockSize))

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

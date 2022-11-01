package flat

import (
	"image"
	"math"

	"github.com/weqqr/panorama/pkg/game"
	"github.com/weqqr/panorama/pkg/lm"
	"github.com/weqqr/panorama/pkg/mesh"
	"github.com/weqqr/panorama/pkg/raster"
	"github.com/weqqr/panorama/pkg/render"
	"github.com/weqqr/panorama/pkg/spatial"
	"github.com/weqqr/panorama/pkg/world"
)

var (
	YOffsetCoef = int(math.Round(render.BaseResolution * (1 + math.Sqrt2) / 4))
	TileSize    = spatial.BlockSize * render.BaseResolution
)

type FlatRenderer struct {
	nr render.NodeRasterizer

	region spatial.Region
	game   *game.Game
}

func NewRenderer(region spatial.Region, game *game.Game) *FlatRenderer {
	return &FlatRenderer{
		nr:     render.NewNodeRasterizer(lm.TopDownProjection()),
		region: region,
		game:   game,
	}
}

func (r *FlatRenderer) renderNode(
	target *raster.RenderBuffer,
	pos spatial.NodePosition,
	worldPos spatial.NodePosition,
	neighborhood *render.BlockNeighborhood,
	offset image.Point,
	depthOffset float64,
) {
	name, param1, param2 := neighborhood.GetNode(pos)

	// Fast path: checking for air immediately is faster than fetching NodeDefinition
	if name == "air" {
		return
	}

	nodeDef := r.game.NodeDef(name)

	needsAlphaBlending := true
	if nodeDef.DrawType == game.DrawTypeNormal {
		needsAlphaBlending = false
	}

	// Estimate lighting by sampling neighboring nodes and using the brightest one
	neighborOffsets := []spatial.NodePosition{
		{X: 0, Y: 1, Z: 0},
	}

	neighborFaces := []mesh.CubeFaces{
		mesh.CubeFaceTop,
	}

	maxParam1 := param1
	hiddenFaces := mesh.CubeFaces(0)
	for i, offset := range neighborOffsets {
		neighborPos := pos.Add(offset)
		neighborName, param1, _ := neighborhood.GetNode(neighborPos)
		if param1 > maxParam1 {
			maxParam1 = param1
		}

		// Compute visibility for stacked liquids
		if nodeDef.DrawType.IsLiquid() {
			hiddenFaces |= mesh.CubeFaceWest | mesh.CubeFaceDown | mesh.CubeFaceSouth

			neighborNodeDef := r.game.NodeDef(neighborName)
			if neighborNodeDef.DrawType.IsLiquid() {
				hiddenFaces |= neighborFaces[i]
			}
		}
	}

	renderableNode := render.RenderableNode{
		Name:        name,
		Light:       render.DecodeLight(maxParam1),
		Param2:      param2,
		HiddenFaces: hiddenFaces,
	}
	renderedNode := r.nr.Render(renderableNode, &nodeDef)

	depthOffset = -float64(pos.Y) + depthOffset
	if needsAlphaBlending {
		target.OverlayDepthAwareWithAlpha(renderedNode, offset, depthOffset)
	} else {
		target.OverlayDepthAware(renderedNode, offset, depthOffset)
	}
}

func (r *FlatRenderer) renderBlock(
	target *raster.RenderBuffer,
	blockPos spatial.BlockPosition,
	neighborhood *render.BlockNeighborhood,
	depthOffset float64,
) {
	for z := spatial.BlockSize - 1; z >= 0; z-- {
		for y := spatial.BlockSize - 1; y >= 0; y-- {
			for x := spatial.BlockSize - 1; x >= 0; x-- {
				nodePos := spatial.NodePosition{X: x, Y: y, Z: z}
				nodeWorldPos := blockPos.AddNode(nodePos)

				if !r.region.Intersects(nodeWorldPos.Region()) {
					continue
				}

				offset := image.Point{
					X: render.BaseResolution * x,
					Y: render.BaseResolution * z,
				}

				r.renderNode(target, nodePos, nodeWorldPos, neighborhood, offset, depthOffset)
			}
		}
	}
}

func (r *FlatRenderer) RenderTile(
	tilePos render.TilePosition,
	world *world.World,
	game *game.Game,
) *raster.RenderBuffer {
	rect := image.Rect(0, 0, TileSize, TileSize)
	target := raster.NewRenderBuffer(rect)

	yMin := int(math.Floor(float64(r.region.YBounds.Min) / float64(spatial.BlockSize)))
	yMax := int(math.Ceil(float64(r.region.YBounds.Max) / float64(spatial.BlockSize)))

	for y := yMin; y < yMax; y++ {
		blockPos := spatial.BlockPosition{
			X: tilePos.X,
			Y: y,
			Z: tilePos.Y,
		}
		neighborhood := render.BlockNeighborhood{}
		neighborhood.FetchBlock(world, spatial.BlockPosition{X: 0, Y: 0, Z: 0}, blockPos)
		neighborhood.FetchBlock(world, spatial.BlockPosition{X: 0, Y: 1, Z: 0}, blockPos)
		depthOffset := -float64(y) * spatial.BlockSize
		r.renderBlock(target, blockPos, &neighborhood, depthOffset)
	}

	return target
}

func (r *FlatRenderer) ProjectRegion(region spatial.Region) spatial.ProjectedRegion {
	return spatial.ProjectedRegion{
		XBounds: spatial.Bounds{
			Min: region.XBounds.Min / spatial.BlockSize,
			Max: region.XBounds.Max / spatial.BlockSize,
		},
		YBounds: spatial.Bounds{
			Min: region.ZBounds.Min / spatial.BlockSize,
			Max: region.ZBounds.Max / spatial.BlockSize,
		},
	}
}

func (r *FlatRenderer) Name() string {
	return "flat"
}

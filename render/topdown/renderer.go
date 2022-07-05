package topdown

import (
	"image"

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
		nr: render.NewNodeRasterizer(lm.TopDownProjection()),

		lowerLimit: lowerLimit,
		upperLimit: upperLimit,
	}
}

func (r *Renderer) renderBlock(target *raster.RenderBuffer, neighborhood *render.BlockNeighborhood, game *game.Game, depthOffset float32) {
	// FIXME: nodes must define their origin points
	for z := 0; z < world.MapBlockSize; z++ {
		for y := 0; y < world.MapBlockSize; y++ {
			for x := 0; x < world.MapBlockSize; x++ {
				tileOffsetX := x * render.BaseResolution
				tileOffsetY := (world.MapBlockSize - 1 - z) * render.BaseResolution

				name, _, param2 := neighborhood.GetNode(x, y, z)
				param1 := neighborhood.GetParam1(x, y+1, z)

				nodeDef := game.NodeDef(name)

				renderableNode := render.RenderableNode{
					Name:   name,
					Light:  render.DecodeLight(param1),
					Param2: param2,
				}
				renderedNode := r.nr.Render(renderableNode, &nodeDef)

				depthOffset := -float32(y)
				target.OverlayDepthAware(renderedNode, image.Pt(tileOffsetX, tileOffsetY), depthOffset)
			}
		}
	}
}

func (r *Renderer) RenderTile(tilePos render.TilePosition, w *world.World, game *game.Game) *image.NRGBA {
	rect := image.Rect(0, 0, render.BaseTileSize, render.BaseTileSize)
	target := raster.NewRenderBuffer(rect)

	centerX := tilePos.X
	centerY := 0
	centerZ := -tilePos.Y

	for i := r.lowerLimit; i < r.upperLimit; i++ {
		blockX := centerX
		blockY := centerY + i
		blockZ := centerZ

		neighborhood := render.BlockNeighborhood{}
		neighborhood.FetchBlock(1, 1, 1, blockX, blockY, blockZ, w)
		neighborhood.FetchBlock(1, 2, 1, blockX, blockY+1, blockZ, w)

		depthOffset := float32(blockY * world.MapBlockSize)
		r.renderBlock(target, &neighborhood, game, depthOffset)
	}

	return target.Color
}

func (r *Renderer) ListTilesWithBlock(x, y, z int) []render.TilePosition {
	return []render.TilePosition{}
}

package flat

import (
	"image"

	"github.com/lord-server/panorama/internal/game"
	"github.com/lord-server/panorama/internal/lm"
	"github.com/lord-server/panorama/internal/raster"
	"github.com/lord-server/panorama/internal/render"
	"github.com/lord-server/panorama/internal/spatial"
	"github.com/lord-server/panorama/internal/world"
)

type FlatRenderer struct {
	nr render.NodeRasterizer

	region spatial.Region
	game   *game.Game
}

func NewRenderer(region spatial.Region, game *game.Game) *FlatRenderer {
	return &FlatRenderer{
		nr:     render.NewNodeRasterizer(lm.DimetricProjection()),
		region: region,
		game:   game,
	}
}

func (r *FlatRenderer) RenderTile(
	tilePos render.TilePosition,
	world *world.World,
	game *game.Game,
) *raster.RenderBuffer {
	rect := image.Rect(0, 0, 256, 256)
	target := raster.NewRenderBuffer(rect)

	return target
}

func (r *FlatRenderer) ProjectRegion(region spatial.Region) spatial.ProjectedRegion {
	return spatial.ProjectedRegion{}
}

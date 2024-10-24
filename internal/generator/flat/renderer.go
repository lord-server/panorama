package flat

import (
	"image"

	"github.com/lord-server/panorama/internal/game"
	"github.com/lord-server/panorama/internal/generator"
	"github.com/lord-server/panorama/internal/generator/rasterizer"
	"github.com/lord-server/panorama/internal/world"
	"github.com/lord-server/panorama/pkg/geom"
	"github.com/lord-server/panorama/pkg/lm"
)

type FlatRenderer struct {
	nr rasterizer.NodeRasterizer

	region geom.Region
	game   *game.Game
}

func NewRenderer(region geom.Region, game *game.Game) *FlatRenderer {
	return &FlatRenderer{
		nr:     rasterizer.New(lm.DimetricProjection()),
		region: region,
		game:   game,
	}
}

func (r *FlatRenderer) RenderTile(
	tilePos generator.TilePosition,
	world *world.World,
	game *game.Game,
) *rasterizer.RenderBuffer {
	rect := image.Rect(0, 0, 256, 256)
	target := rasterizer.NewRenderBuffer(rect)

	return target
}

func (r *FlatRenderer) ProjectRegion(region geom.Region) geom.ProjectedRegion {
	return geom.ProjectedRegion{}
}

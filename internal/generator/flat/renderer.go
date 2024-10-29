package flat

import (
	"image"
	"log/slog"

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
	wd *world.World,
	game *game.Game,
) *rasterizer.RenderBuffer {
	rect := image.Rect(0, 0, 256, 256)
	target := rasterizer.NewRenderBuffer(rect)

	err := wd.GetBlocksAlongY(tilePos.X, tilePos.Y, func(pos geom.BlockPosition, block *world.MapBlock) error {
		return nil
	})
	if err != nil {
		slog.Error("unable to get blocks", "error", err)
	}

	return target
}

func (r *FlatRenderer) ProjectRegion(region geom.Region) geom.ProjectedRegion {
	return geom.ProjectedRegion{
		XBounds: geom.Bounds{
			Min: region.XBounds.Min / 16,
			Max: region.XBounds.Max / 16,
		},
		YBounds: geom.Bounds{
			Min: region.ZBounds.Min / 16,
			Max: region.ZBounds.Max / 16,
		},
	}
}

package tile

import (
	"github.com/lord-server/panorama/internal/game"
	"github.com/lord-server/panorama/internal/world"
	"github.com/lord-server/panorama/pkg/geom"
)

type NextTiler struct {
}

func (t *NextTiler) FullRender(game *game.Game, world *world.World, workers int, region geom.Region, createRenderer CreateRendererFunc) {

}

package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/render/isometric"
	"github.com/weqqr/panorama/world"
)

func tilePath(x, y, zoom int) string {
	return fmt.Sprintf("tiles/%v/%v/%v.png", zoom, x, y)
}

type Tiler struct {
	xMin, yMin int
	xMax, yMax int
	upperLimit int
	lowerLimit int
}

func NewTiler(region *RegionConfig) Tiler {
	return Tiler{
		xMin:       region.XBounds[0],
		yMin:       region.YBounds[0],
		xMax:       region.XBounds[1],
		yMax:       region.YBounds[1],
		upperLimit: region.UpperLimit,
		lowerLimit: region.LowerLimit,
	}
}

func (t *Tiler) FullRender(game *game.Game, world *world.World) {
	var wg sync.WaitGroup
	for x := t.xMin; x < t.xMax; x++ {
		os.MkdirAll(fmt.Sprintf("tiles/0/%v", x), os.ModePerm)
		xx := x
		wg.Add(1)
		go func() {
			defer wg.Done()

			renderer := isometric.NewRenderer(t.lowerLimit, t.upperLimit)

			for y := t.yMin; y < t.yMax; y++ {
				yy := y
				output := renderer.RenderTile(render.TilePosition{xx, yy}, world, game)
				tilePath := tilePath(xx, yy, 0)
				err := raster.SavePNG(output, tilePath)
				if err != nil {
					return
				}
				log.Printf("saved %v", tilePath)
			}
		}()
	}

	wg.Wait()
}

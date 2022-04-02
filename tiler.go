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

func worker(wg *sync.WaitGroup, game *game.Game, world *world.World, renderer isometric.Renderer, positions <-chan render.TilePosition) {
	for position := range positions {
		output := renderer.RenderTile(position, world, game)
		tilePath := tilePath(position.X, position.Y, 0)
		err := raster.SavePNG(output, tilePath)
		if err != nil {
			return
		}
		log.Printf("saved %v", tilePath)
	}

	wg.Done()
}

func (t *Tiler) FullRender(game *game.Game, world *world.World, workers int) {
	var wg sync.WaitGroup
	positions := make(chan render.TilePosition)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		renderer := isometric.NewRenderer(t.lowerLimit, t.upperLimit)
		go worker(&wg, game, world, renderer, positions)
	}

	for x := t.xMin; x < t.xMax; x++ {
		os.MkdirAll(fmt.Sprintf("tiles/0/%v", x), os.ModePerm)

		for y := t.yMin; y < t.yMax; y++ {
			positions <- render.TilePosition{X: x, Y: y}
		}
	}

	close(positions)

	wg.Wait()
}

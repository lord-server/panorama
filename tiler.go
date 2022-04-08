package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/world"
)

func tilePath(renderer string, x, y, zoom int) string {
	return fmt.Sprintf("tiles/%v/%v/%v/%v.png", renderer, zoom, x, y)
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

func worker(wg *sync.WaitGroup, game *game.Game, world *world.World, rendererName string, renderer render.Renderer, positions <-chan render.TilePosition) {
	for position := range positions {
		tilePath := tilePath(rendererName, position.X, position.Y, 0)
		output := renderer.RenderTile(position, world, game)
		err := raster.SavePNG(output, tilePath)
		if err != nil {
			log.Printf("worker encountered an error (tile skipped): %v", err)
			continue
		}
		log.Printf("saved %v", tilePath)
	}

	wg.Done()
}

func (t *Tiler) FullRender(game *game.Game, world *world.World, workers int, rendererName string, createRenderer func(int, int) render.Renderer) {
	var wg sync.WaitGroup
	positions := make(chan render.TilePosition)

	for x := t.xMin; x < t.xMax; x++ {
		os.MkdirAll(fmt.Sprintf("tiles/%v/0/%v", rendererName, x), os.ModePerm)
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		renderer := createRenderer(t.lowerLimit, t.upperLimit) // topdown.NewRenderer(t.lowerLimit, t.upperLimit)
		go worker(&wg, game, world, rendererName, renderer, positions)
	}

	for x := t.xMin; x < t.xMax; x++ {
		for y := t.yMin; y < t.yMax; y++ {
			positions <- render.TilePosition{X: x, Y: y}
		}
	}

	close(positions)

	wg.Wait()
}

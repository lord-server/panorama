package tile

import (
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/weqqr/panorama/config"
	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/render/isometric"
	"github.com/weqqr/panorama/world"
)

type Tiler struct {
	xMin, yMin int
	xMax, yMax int
	upperLimit int
	lowerLimit int

	tilesPath string
}

func NewTiler(region *config.RegionConfig, tilesPath string) Tiler {
	return Tiler{
		xMin:       region.XBounds[0],
		yMin:       region.YBounds[0],
		xMax:       region.XBounds[1],
		yMax:       region.YBounds[1],
		upperLimit: region.UpperLimit,
		lowerLimit: region.LowerLimit,

		tilesPath: tilesPath,
	}
}

func (t *Tiler) tilePath(x, y, zoom int) string {
	return fmt.Sprintf("%v/%v/%v/%v.png", t.tilesPath, -zoom, x, y)
}

func (t *Tiler) worker(wg *sync.WaitGroup, game *game.Game, world *world.World, renderer isometric.Renderer, positions <-chan render.TilePosition) {
	for position := range positions {
		output := renderer.RenderTile(position, world, game)
		tilePath := t.tilePath(position.X, position.Y, 0)
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
		go t.worker(&wg, game, world, renderer, positions)
	}

	for x := t.xMin; x < t.xMax; x++ {
		err := os.MkdirAll(fmt.Sprintf("tiles/0/%v", x), os.ModePerm)
		if err != nil {
			panic(err)
		}

		for y := t.yMin; y < t.yMax; y++ {
			positions <- render.TilePosition{X: x, Y: y}
		}
	}

	close(positions)

	wg.Wait()
}

// DownscaleTiles rescales high-resolution tiles into lower resolution ones until it reaches adequate zoom level
func (t *Tiler) DownscaleTiles() {
	mapSize := math.Max(float64(t.xMax-t.xMin), float64(t.yMax-t.yMin))
	zoomLevels := int(math.Ceil(math.Log2(mapSize)))

	log.Printf("Downscaling mapSize=%v zoomLevels=%v", mapSize, zoomLevels)

	tileDir, err := filepath.Abs("tiles/0")
	if err != nil {
		panic(err)
	}

	// Collect tile positions
	var positions []render.TilePosition
	err = filepath.WalkDir(tileDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		dir, file := filepath.Split(path)

		y, err := strconv.Atoi(strings.TrimSuffix(file, filepath.Ext(file)))
		if err != nil {
			return nil
		}

		x, err := strconv.Atoi(filepath.Base(dir))
		if err != nil {
			return nil
		}

		positions = append(positions, render.TilePosition{X: floorDiv(x, 2), Y: floorDiv(y, 2)})

		return nil
	})

	if err != nil {
		panic(err)
	}

	positions = uniquePositions(positions)

	for zoom := 1; zoom <= zoomLevels; zoom++ {
		positions = t.downscalePositions(zoom, positions)
	}
}

package tile

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/weqqr/panorama/pkg/game"
	"github.com/weqqr/panorama/pkg/lm"
	"github.com/weqqr/panorama/pkg/raster"
	"github.com/weqqr/panorama/pkg/region"
	"github.com/weqqr/panorama/pkg/render"
	"github.com/weqqr/panorama/pkg/world"
)

type Tiler struct {
	region     region.Region
	zoomLevels int
	tilesPath  string
}

func NewTiler(region region.Region, zoomLevels int, tilesPath string) Tiler {
	return Tiler{
		region:     region,
		zoomLevels: zoomLevels,
		tilesPath:  tilesPath,
	}
}

func (t *Tiler) tilePath(x, y, zoom int) string {
	return fmt.Sprintf("%v/%v/%v/%v.png", t.tilesPath, -zoom, x, y)
}

func (t *Tiler) worker(wg *sync.WaitGroup, game *game.Game, world *world.World, renderer render.Renderer, positions <-chan render.TilePosition) {
	for position := range positions {
		output := renderer.RenderTile(position, world, game)
		// Don't save empty tiles
		if !output.Dirty {
			continue
		}

		tilePath := t.tilePath(position.X, position.Y, 0)
		err := raster.SavePNG(output.Color, tilePath)
		if err != nil {
			return
		}
		log.Printf("saved %v", tilePath)
	}

	wg.Done()
}

type CreateRendererFunc func() render.Renderer

func (t *Tiler) FullRender(game *game.Game, world *world.World, workers int, region region.TileRegion, createRenderer CreateRendererFunc) {
	var wg sync.WaitGroup
	positions := make(chan render.TilePosition)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		renderer := createRenderer()
		go t.worker(&wg, game, world, renderer, positions)
	}

	for x := region.XBounds.Min; x < region.XBounds.Max; x++ {
		err := os.MkdirAll(fmt.Sprintf("%v/%v", path.Join(t.tilesPath, "0"), x), os.ModePerm)
		if err != nil {
			panic(err)
		}

		for y := region.YBounds.Min; y < region.YBounds.Max; y++ {
			positions <- render.TilePosition{X: x, Y: y}
		}
	}
	close(positions)

	wg.Wait()
}

// DownscaleTiles rescales high-resolution tiles into lower resolution ones until it reaches adequate zoom level
func (t *Tiler) DownscaleTiles() {
	log.Printf("Downscaling zoomLevels=%v", t.zoomLevels)

	tileDir, err := filepath.Abs(path.Join(t.tilesPath, "0"))
	if err != nil {
		panic(err)
	}

	// Collect tile positions
	var positions []render.TilePosition
	err = filepath.WalkDir(tileDir, func(path string, d fs.DirEntry, _ error) error {
		if d == nil || d.IsDir() {
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

		positions = append(positions, render.TilePosition{
			X: lm.FloorDiv(x, 2),
			Y: lm.FloorDiv(y, 2),
		})

		return nil
	})

	if err != nil {
		panic(err)
	}

	positions = uniquePositions(positions)

	for zoom := 1; zoom <= t.zoomLevels; zoom++ {
		log.Printf("Rescaling tiles for zoom level %v", zoom)
		positions = t.downscalePositions(zoom, positions)
	}
}

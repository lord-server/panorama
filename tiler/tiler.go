package tiler

import (
	"fmt"
	"image"
	"image/draw"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/nfnt/resize"

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
}

func NewTiler(region *config.RegionConfig) Tiler {
	return Tiler{
		xMin:       region.XBounds[0],
		yMin:       region.YBounds[0],
		xMax:       region.XBounds[1],
		yMax:       region.YBounds[1],
		upperLimit: region.UpperLimit,
		lowerLimit: region.LowerLimit,
	}
}

func tilePath(x, y, zoom int) string {
	return fmt.Sprintf("tiles/%v/%v/%v.png", -zoom, x, y)
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

func uniquePositions(input []render.TilePosition) []render.TilePosition {
	// Slices with zero or one element always contain unique elements
	if len(input) < 2 {
		return input
	}

	// Sort positions by their coordinates
	sort.Slice(input, func(i, j int) bool {
		if input[i].X < input[j].X {
			return true
		}
		if input[i].X > input[j].X {
			return false
		}
		if input[i].Y < input[j].Y {
			return true
		}
		if input[i].Y > input[j].Y {
			return false
		}
		return false
	})

	// Loop over the slice and skip repeating elements
	j := 1
	for i := 1; i < len(input); i++ {
		// Skip element if it repeats
		if input[i] == input[i-1] {
			continue
		}

		// Rewrite repeated elements with unique ones
		input[j] = input[i]
		j++
	}

	return input[:j]
}

// floorDiv returns the result of floor division. The difference compared to usual division
// is that floor division always rounds down instead of rounding towards zero.
func floorDiv(a, b int) int {
	return int(math.Floor(float64(a) / float64(b)))
}

// downscalePositions produces downscaled images for given zoom level and returns a list of produced tile positions
func downscalePositions(zoom int, positions []render.TilePosition) []render.TilePosition {
	const quadrantSize = 128

	log.Printf("zoom=%v positions=%v", zoom, positions)

	var nextPositions []render.TilePosition

	for _, pos := range positions {
		target := image.NewNRGBA(image.Rect(0, 0, 256, 256))

		for quadrantY := 0; quadrantY < 2; quadrantY++ {
			for quadrantX := 0; quadrantX < 2; quadrantX++ {
				log.Printf("quad")
				source, err := raster.LoadPNG(tilePath(pos.X*2+quadrantX, pos.Y*2+quadrantY, zoom-1))
				if err != nil {
					continue
				}

				quadrant := resize.Resize(quadrantSize, quadrantSize, source, resize.Lanczos3)

				targetX := quadrantX * quadrantSize
				targetY := quadrantY * quadrantSize
				draw.Draw(target, image.Rect(targetX, targetY, targetX+quadrantSize, targetY+quadrantSize), quadrant, image.Pt(0, 0), draw.Src)
			}
		}

		err := raster.SavePNG(target, tilePath(pos.X, pos.Y, zoom))
		if err != nil {
			panic(err)
		}

		nextPositions = append(nextPositions, render.TilePosition{X: floorDiv(pos.X, 2), Y: floorDiv(pos.Y, 2)})
	}

	nextPositions = uniquePositions(nextPositions)
	return nextPositions
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
		positions = downscalePositions(zoom, positions)
	}
}

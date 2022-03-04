package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/nfnt/resize"
	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/world"
)

func nextPowerOfTwo(value int) int {
	i := 1
	for i < value {
		i *= 2
	}
	return i
}

func tilePath(x, y, zoom int) string {
	return fmt.Sprintf("tiles/%v/%v/%v.png", zoom, x, y)
}

func renderTiles(game *game.Game, world *world.World, renderer render.Renderer, min, max int) {
	var wg sync.WaitGroup
	for x := min; x < max; x++ {
		os.MkdirAll(fmt.Sprintf("tiles/0/%v", x), os.ModePerm)
		xx := x
		wg.Add(1)
		go func() {
			defer wg.Done()
			for y := min; y < max; y++ {
				yy := y
				output := renderer.RenderTile(render.TilePosition{xx, yy * 2}, world, game)
				path := tilePath(xx, yy, 0)
				err := raster.SavePNG(output, path)
				if err != nil {
					return
				}
				log.Printf("saved %v", path)
			}
		}()
	}

	wg.Wait()
}

func downscaleTiles(min, max int) {
	mapSize := nextPowerOfTwo(max - min)
	maxZoom := int(math.Ceil(math.Log2(float64(mapSize))))

	for z := 1; z <= maxZoom; z++ {
		size := mapSize / int(math.Pow(2, float64(z-1)))
		min := -size / 2
		max := size / 2
		log.Printf("Processing zoom level %v (min=%v, max=%v)", z, min, max)

		for x := min; x < max; x++ {
			os.MkdirAll(fmt.Sprintf("tiles/%v/%v", -z, x), os.ModePerm)
			for y := min; y < max; y++ {
				target := image.NewNRGBA(image.Rect(0, 0, 256, 256))
				targetContainsTiles := false

				for quadrantX := 0; quadrantX < 2; quadrantX++ {
					for quadrantY := 0; quadrantY < 2; quadrantY++ {
						quadrant, _ := raster.LoadPNG(tilePath(2*x+quadrantX, 2*y+quadrantY, -z+1))
						if quadrant == nil {
							continue
						}

						quadrantXOffset := 128 * quadrantX
						quadrantYOffset := 128 * quadrantY
						targetRect := image.Rect(quadrantXOffset, quadrantYOffset, 128+quadrantXOffset, 128+quadrantYOffset)
						draw.Draw(target, targetRect, resize.Resize(128, 128, quadrant, resize.Lanczos3), image.Pt(0, 0), draw.Src)
						targetContainsTiles = true
					}
				}

				if targetContainsTiles {
					log.Printf("saved %s", tilePath(x, y, -z))
					raster.SavePNG(target, tilePath(x, y, -z))
				}
			}
		}
	}
}

func serveTiles() {
	staticFS := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFS))

	tilesFS := http.FileServer(http.Dir("./tiles"))
	http.Handle("/tiles/", http.StripPrefix("/tiles/", tilesFS))

	err := http.ListenAndServe(":1337", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	config := LoadConfig("panorama.toml")
	log.Printf("game path: %v\n", config.GamePath)

	descPath := path.Join(config.WorldPath, "panorama_nodes.json")
	log.Printf("game description: `%v`\n", descPath)
}

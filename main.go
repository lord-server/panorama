package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/world"
)

func savePNG(img *image.NRGBA, name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}

	if err := png.Encode(file, img); err != nil {
		file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func main() {
	web := true
	generate := true

	config := LoadConfig("panorama.toml")
	log.Printf("path: %v\n", config.Game.Path)
	log.Printf("description: `%v`\n", config.Game.Desc)

	if generate {
		game, err := game.LoadGame(config.Game.Desc, config.Game.Path)
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Loaded %v nodes, %v aliases\n", len(game.Nodes), len(game.Aliases))

		log.Printf("Using %v as backend", config.World.Backend)

		backend, err := world.NewPgBackend(config.World.Connection)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Connected to %v", config.World.Backend)

		world := world.NewWorldWithBackend(backend)
		nr := render.NewNodeRasterizer()

		var wg sync.WaitGroup
		for x := -30; x < 30; x++ {
			os.MkdirAll(fmt.Sprintf("tiles/%v", x), os.ModePerm)
			xx := x
			wg.Add(1)
			go func() {
				defer wg.Done()
				for y := -30; y < 30; y++ {
					yy := y
					output := render.RenderTile(xx-52, yy*2-3, &nr, &world, &game)
					filename := fmt.Sprintf("%v.png", yy)
					path := filepath.Join(fmt.Sprintf("tiles/%v", xx), filename)
					savePNG(output, path)
					log.Printf("saved %v", path)
				}
			}()
		}

		wg.Wait()
	}

	if web {
		staticFS := http.FileServer(http.Dir("./static"))
		http.Handle("/static/", http.StripPrefix("/static/", staticFS))

		tilesFS := http.FileServer(http.Dir("./tiles"))
		http.Handle("/tiles/", http.StripPrefix("/tiles/", tilesFS))

		err := http.ListenAndServe(":1337", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

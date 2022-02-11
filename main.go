package main

import (
	"image"
	"image/png"
	"log"
	"os"

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
	config := LoadConfig("panorama.toml")

	log.Printf("path: %v\n", config.Game.Path)
	log.Printf("description: `%v`\n", config.Game.Desc)

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
	block, err := world.GetBlock(-1, 0, -5)
	if err != nil {
		log.Panic(err)
	}

	nr := render.NewNodeRasterizer()
	output := render.RenderBlock(&nr, block, &game)

	savePNG(output, "test.png")
}

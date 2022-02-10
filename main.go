package main

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/render"
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

	node := "default:stone"
	gameNode := game.Node(node)
	renderer := render.NewNodeRasterizer()
	output := renderer.Render(&gameNode)
	savePNG(output, "test.png")
}

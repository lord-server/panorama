package main

import (
	"image"
	"image/png"
	"log"
	"os"
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

	game, err := LoadGame(config.Game.Desc, config.Game.Path)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Loaded %v nodes, %v aliases\n", len(game.nodes), len(game.aliases))

	node := "default:stone"
	nodedef := game.NodeDef(node)
	renderer := NewNodeRasterizer()
	output := renderer.Render(&nodedef)
	savePNG(output, "test.png")
}

package main

import (
	"flag"
	"log"
	"path"

	"github.com/weqqr/panorama/pkg/config"
	"github.com/weqqr/panorama/pkg/game"
	"github.com/weqqr/panorama/pkg/render"
	"github.com/weqqr/panorama/pkg/render/isometric"
	"github.com/weqqr/panorama/pkg/tile"
	"github.com/weqqr/panorama/pkg/web"
	"github.com/weqqr/panorama/pkg/world"
)

type Args struct {
	FullRender bool
	Downscale  bool
	Serve      bool
	ConfigPath string
}

var args Args

func init() {
	flag.BoolVar(&args.FullRender, "fullrender", false, "Render entire map")
	flag.BoolVar(&args.Downscale, "downscale", false, "Downscale existing tiles (--fullrender does this automatically)")
	flag.BoolVar(&args.Serve, "serve", false, "Serve tiles over the web")
	flag.StringVar(&args.ConfigPath, "config", "config.toml", "Path to config file")
	flag.Parse()
}

func main() {
	log.Printf("Config path: `%v`", args.ConfigPath)
	config, err := config.LoadConfig(args.ConfigPath)
	if err != nil {
		log.Fatalf("Unable to load config: %v\n", err)
	}

	log.Printf("Game path: `%v`\n", config.System.GamePath)

	descPath := path.Join(config.System.WorldPath, "nodes_dump.json")
	log.Printf("Game description: `%v`\n", descPath)

	game, err := game.LoadGame(descPath, config.System.GamePath)
	if err != nil {
		log.Fatalf("Unable to load game description: %v\n", err)
	}

	backend, err := world.NewPostgresBackend(config.System.WorldDSN)
	if err != nil {
		log.Fatalf("Unable to connect to world DB: %v\n", err)
	}

	world := world.NewWorldWithBackend(backend)

	tiler := tile.NewTiler(config.Region, config.Renderer.ZoomLevels, config.System.TilesPath)

	if args.FullRender {
		log.Printf("Performing a full render using %v workers", config.Renderer.Workers)

		log.Printf("Region: %v", config.Region)

		tiler.FullRender(&game, &world, config.Renderer.Workers, config.Region, func() render.Renderer {
			return isometric.NewRenderer(config.Region, &game)
		})
	}

	if args.Downscale || args.FullRender {
		tiler.DownscaleTiles()
	}

	if args.Serve {
		log.Printf("Serving tiles @ %v", config.Web.ListenAddress)
		web.Serve(&config)
	}
}

package main

import (
	"flag"
	"log"
	"path"

	"github.com/weqqr/panorama/pkg/config"
	"github.com/weqqr/panorama/pkg/game"
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

	log.Printf("Game path: `%v`\n", config.GamePath)

	descPath := path.Join(config.WorldPath, "nodes_dump.json")
	log.Printf("Game description: `%v`\n", descPath)

	game, err := game.LoadGame(descPath, config.GamePath)
	if err != nil {
		log.Fatalf("Unable to load game description: %v\n", err)
	}

	backend, err := world.NewPostgresBackend(config.WorldDSN)
	if err != nil {
		log.Fatalf("Unable to connect to world DB: %v\n", err)
	}

	world := world.NewWorldWithBackend(backend)

	tiler := tile.NewTiler(&config.RegionConfig, config.ZoomLevels, config.TilesPath)

	if args.FullRender {
		log.Printf("Performing a full render using %v workers", config.RendererWorkers)
		tiler.FullRender(&game, &world, config.RendererWorkers)
	}

	if args.Downscale || args.FullRender {
		tiler.DownscaleTiles()
	}

	if args.Serve {
		log.Printf("Serving tiles @ %v", config.ListenAddress)
		web.Serve(config.ListenAddress, &config)
	}
}

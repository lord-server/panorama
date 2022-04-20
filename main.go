package main

import (
	"flag"
	"log"
	"net/http"
	"path"

	"github.com/weqqr/panorama/config"
	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/tiler"
	"github.com/weqqr/panorama/world"
)

func serveTiles(addr string) {
	staticFS := http.FileServer(http.Dir("./static"))
	http.Handle("/", staticFS)

	tilesFS := http.FileServer(http.Dir("./tiles"))
	http.Handle("/tiles/", http.StripPrefix("/tiles/", tilesFS))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Args struct {
	FullRender bool
	Serve      bool
	ConfigPath string
}

var args Args

func init() {
	flag.BoolVar(&args.FullRender, "fullrender", false, "Render entire map")
	flag.BoolVar(&args.Serve, "serve", false, "Serve tiles over the web")
	flag.StringVar(&args.ConfigPath, "config", "config.toml", "Path to config file")
	flag.Parse()
}

func main() {
	config := config.LoadConfig(args.ConfigPath)
	log.Printf("game path: %v\n", config.GamePath)

	descPath := path.Join(config.WorldPath, "panorama_nodes.json")
	log.Printf("game description: `%v`\n", descPath)

	game, err := game.LoadGame(descPath, config.GamePath)
	if err != nil {
		panic(err)
	}

	backend, err := world.NewPostgresBackend(config.WorldDSN)
	if err != nil {
		panic(err)
	}

	world := world.NewWorldWithBackend(backend)

	tiler := tiler.NewTiler(&config.RegionConfig)

	if args.FullRender {
		tiler.FullRender(&game, &world, config.RendererWorkers)
		tiler.DownscaleTiles()
	}

	if args.Serve {
		serveTiles(config.ListenAddress)
	}
}

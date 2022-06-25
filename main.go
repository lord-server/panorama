package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"path"

	"github.com/weqqr/panorama/config"
	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/tile"
	"github.com/weqqr/panorama/world"
)

type metadataHandler struct {
	config *config.Config
}

type MapMetadata struct {
	Title      string `json:"title"`
	ZoomLevels int    `json:"zoomLevels"`
}

func (m *metadataHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MapMetadata{
		Title:      m.config.Title,
		ZoomLevels: m.config.ZoomLevels,
	})
}

func serve(addr string, config *config.Config) {
	staticFS := http.FileServer(http.Dir("./static"))
	http.Handle("/", staticFS)

	tilesFS := http.FileServer(http.Dir(config.TilesPath))
	http.Handle("/tiles/", http.StripPrefix("/tiles/", tilesFS))

	http.Handle("/metadata.json", &metadataHandler{
		config: config,
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

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
		serve(config.ListenAddress, &config)
	}
}

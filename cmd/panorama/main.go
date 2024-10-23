package main

import (
	"io/fs"
	"log/slog"
	"os"
	"path"

	"github.com/alexflint/go-arg"
	"github.com/lord-server/panorama/internal/config"
	"github.com/lord-server/panorama/internal/game"
	"github.com/lord-server/panorama/internal/render"
	"github.com/lord-server/panorama/internal/render/isometric"
	"github.com/lord-server/panorama/internal/tile"
	"github.com/lord-server/panorama/internal/web"
	"github.com/lord-server/panorama/internal/world"
)

var static fs.FS

type FullRenderArgs struct{}

type RunArgs struct{}

var args struct {
	ConfigPath string          `arg:"-c,--config" default:"config.toml"`
	FullRender *FullRenderArgs `arg:"subcommand:fullrender"`
	Run        *RunArgs        `arg:"subcommand:run"`
}

func main() {
	arg.MustParse(&args)

	config, err := config.LoadConfig(args.ConfigPath)
	if err != nil {
		slog.Error("unable to load config", "error", err)
		os.Exit(1)
	}

	switch {
	case args.Run != nil:
		err = run(config)

	case args.FullRender != nil:
		err = fullrender(config)

	default:
		slog.Warn("command not specified, proceeding with run")
		err = run(config)
	}

	if err != nil {
		os.Exit(1)
	}
}

func fullrender(config config.Config) error {
	descPath := path.Join(config.System.WorldPath, "nodes_dump.json")

	slog.Info("loading game description", "game", config.System.GamePath, "mods", config.System.ModPath, "desc", descPath)

	game, err := game.LoadGame(descPath, config.System.GamePath, config.System.ModPath)
	if err != nil {
		slog.Error("unable to load game description", "error", err)
		return err
	}

	backend, err := world.NewPostgresBackend(config.System.WorldDSN)
	if err != nil {
		slog.Error("unable to connect to world DB", "error", err)
		return err
	}

	world := world.NewWorldWithBackend(backend)

	tiler := tile.NewTiler(config.Region, config.Renderer.ZoomLevels, config.System.TilesPath)

	slog.Info("performing a full render", "workers", config.Renderer.Workers, "region", config.Region)

	tiler.FullRender(&game, &world, config.Renderer.Workers, config.Region, func() render.Renderer {
		return isometric.NewRenderer(config.Region, &game)
	})

	tiler.DownscaleTiles()

	return nil
}

func run(config config.Config) error {
	quit := make(chan bool)

	slog.Info("starting web server", "address", config.Web.ListenAddress)

	go func() {
		web.Serve(static, &config)
		quit <- true
	}()

	<-quit

	return nil
}

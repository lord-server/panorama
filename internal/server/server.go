package server

import (
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lord-server/panorama/internal/config"
)

func Serve(static fs.FS, config *config.Config) {
	router := chi.NewRouter()

	staticRootDir, err := fs.Sub(static, "ui/build")
	if err != nil {
		panic(err)
	}

	router.Handle("/*", http.FileServer(http.FS(staticRootDir)))
	router.Handle("/tiles/*", http.StripPrefix("/tiles", http.FileServer(http.Dir(config.System.TilesPath))))

	httpServer := &http.Server{
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		Addr:              config.Web.ListenAddress,
		Handler:           router,
	}

	err = httpServer.ListenAndServe()
	if err != nil {
		slog.Error("failed to start web server", "err", err)
	}
}

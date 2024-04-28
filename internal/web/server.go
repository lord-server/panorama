package web

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lord-server/panorama/internal/config"
)

type server struct {
	config *config.Config
}

func sendJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func (s *server) GetMetadata(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, map[string]any{
		"center": []int{700, 200},
	})
}

func Serve(static fs.FS, config *config.Config) {
	server := &server{
		config: config,
	}

	router := chi.NewRouter()

	staticRootDir, err := fs.Sub(static, "ui/build")
	if err != nil {
		panic(err)
	}

	router.Handle("/*", http.FileServer(http.FS(staticRootDir)))
	router.Handle("/tiles/*", http.StripPrefix("/tiles", http.FileServer(http.Dir(config.System.TilesPath))))
	router.Get("/api/v1/metadata", server.GetMetadata)

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

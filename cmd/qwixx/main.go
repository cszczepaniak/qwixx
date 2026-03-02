package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/cszczepaniak/qwixx/internal/game"
	"github.com/cszczepaniak/qwixx/internal/handlers"
	"github.com/cszczepaniak/qwixx/internal/views"
)

//go:embed static/*
var staticFS embed.FS

// responseRecorder wraps ResponseWriter to detect if a response was written.
type responseRecorder struct {
	http.ResponseWriter
	status int
	written bool
}

func (r *responseRecorder) WriteHeader(code int) {
	if !r.written {
		r.written = true
		r.status = code
		r.ResponseWriter.WriteHeader(code)
	}
}

// adapt converts a handler that returns an error into http.HandlerFunc.
// It logs errors with slog and sends an appropriate status if the handler didn't write a response.
func adapt(h handlers.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		err := h(rec, r)
		if err == nil {
			return
		}
		slog.Error("handler error", "method", r.Method, "path", r.URL.Path, "err", err)
		if !rec.written {
			var he *handlers.HTTPError
			if errors.As(err, &he) {
				http.Error(w, he.Error(), he.Code)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func main() {
	// Static assets: strip "static" prefix so /static/dist.css serves dist.css
	staticFiles, err := fs.Sub(staticFS, "static")
	if err != nil {
		slog.Error("static fs sub", "err", err)
		os.Exit(1)
	}
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles)))

	mux := http.NewServeMux()
	mux.Handle("/static/", staticHandler)
	// Use "/" + checks instead of "GET /": Go 1.26 ServeMux rejects "GET /" when "/static/" is registered (pattern conflict).
	mux.HandleFunc("/", adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return nil
		}
		if r.Method != http.MethodGet {
			return handlers.Error(http.StatusMethodNotAllowed, fmt.Errorf("method not allowed: %s", r.Method))
		}
		sigs := handlers.Signals{Game: game.State{}}
		sigs.PopulateDerived()
		gameJSON, err := json.Marshal(sigs.Game)
		if err != nil {
			return fmt.Errorf("marshal game: %w", err)
		}
		disabledJSON, err := json.Marshal(sigs.DisabledIndices)
		if err != nil {
			return fmt.Errorf("marshal disabled indices: %w", err)
		}
		unlockedJSON, err := json.Marshal(sigs.UnlockedLocks)
		if err != nil {
			return fmt.Errorf("marshal unlocked locks: %w", err)
		}
		if err := views.Index(string(gameJSON), string(disabledJSON), string(unlockedJSON), fmt.Sprintf("%d", sigs.Score)).Render(r.Context(), w); err != nil {
			return fmt.Errorf("render index: %w", err)
		}
		return nil
	}))
	mux.HandleFunc("POST /action/set-cross", adapt(handlers.SetCross))
	mux.HandleFunc("POST /action/unset-cross", adapt(handlers.UnsetCross))
	mux.HandleFunc("POST /action/set-missed", adapt(handlers.SetMissed))
	mux.HandleFunc("POST /action/unset-missed", adapt(handlers.UnsetMissed))
	mux.HandleFunc("POST /action/clear-all", adapt(handlers.ClearAll))

	port := "3000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	slog.Info("listening", "addr", "http://localhost:"+port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		slog.Error("listen failed", "err", err)
		os.Exit(1)
	}
}

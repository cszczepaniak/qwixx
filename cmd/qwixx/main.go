package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/cszczepaniak/qwixx/internal/game"
	"github.com/cszczepaniak/qwixx/internal/handlers"
	"github.com/cszczepaniak/qwixx/internal/views"
)

//go:embed static/*
var staticFS embed.FS

func main() {
	// Static assets: strip "static" prefix so /static/dist.css serves dist.css
	staticFiles, _ := fs.Sub(staticFS, "static")
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles)))

	mux := http.NewServeMux()
	mux.Handle("/static/", staticHandler)
	// Use "/" + checks instead of "GET /": Go 1.26 ServeMux rejects "GET /" when "/static/" is registered (pattern conflict).
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		sigs := handlers.Signals{Game: game.State{}}
		sigs.PopulateDerived()
		gameJSON, err := json.Marshal(sigs.Game)
		if err != nil {
			http.Error(w, "Could not marshal game", http.StatusInternalServerError)
			return
		}
		disabledJSON, err := json.Marshal(sigs.DisabledIndices)
		if err != nil {
			http.Error(w, "Could not marshal disabled indices", http.StatusInternalServerError)
			return
		}
		unlockedJSON, err := json.Marshal(sigs.UnlockedLocks)
		if err != nil {
			http.Error(w, "Could not marshal unlocked locks", http.StatusInternalServerError)
			return
		}
		views.Index(string(gameJSON), string(disabledJSON), string(unlockedJSON), fmt.Sprintf("%d", sigs.Score)).Render(r.Context(), w)
	})
	mux.HandleFunc("POST /action/set-cross", handlers.SetCross)
	mux.HandleFunc("POST /action/unset-cross", handlers.UnsetCross)
	mux.HandleFunc("POST /action/set-missed", handlers.SetMissed)
	mux.HandleFunc("POST /action/unset-missed", handlers.UnsetMissed)
	mux.HandleFunc("POST /action/clear-all", handlers.ClearAll)

	port := "3000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Listening on http://localhost:%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		log.Fatal(err)
	}
}

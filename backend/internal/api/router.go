package api

import (
	"drawl/internal/hub"
	"drawl/internal/ws"
	"net/http"
	"os"
	"path/filepath"
)

func NewRouter(h *hub.Hub, registry *ws.ClientRegistry, wsHandler *ws.Handler, handlers *Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /api/games", RateLimitMiddleware(http.HandlerFunc(handlers.CreateGame)))
	mux.HandleFunc("POST /api/games/join", handlers.JoinGame)
	mux.Handle("/ws", wsHandler)

	// Serve static frontend if the directory exists
	staticDir := "./static"
	if _, err := os.Stat(staticDir); err == nil {
		fs := http.FileServer(http.Dir(staticDir))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Try the exact file first; fall back to index.html for SPA routing
			path := filepath.Join(staticDir, r.URL.Path)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
				return
			}
			fs.ServeHTTP(w, r)
		})
	}

	return CORSMiddleware(mux)
}

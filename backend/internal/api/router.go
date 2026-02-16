package api

import (
	"drawl/internal/hub"
	"drawl/internal/ws"
	"net/http"
)

func NewRouter(h *hub.Hub, registry *ws.ClientRegistry, wsHandler *ws.Handler, handlers *Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/games", handlers.CreateGame)
	mux.HandleFunc("POST /api/games/join", handlers.JoinGame)
	mux.Handle("/ws", wsHandler)

	return CORSMiddleware(mux)
}

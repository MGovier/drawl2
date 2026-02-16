package api

import (
	"drawl/internal/game"
	"drawl/internal/hub"
	"drawl/internal/ws"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Handlers struct {
	Hub      *hub.Hub
	Registry *ws.ClientRegistry
	AI       game.AIHandler
}

type createGameRequest struct {
	PlayerName string `json:"playerName"`
}

type createGameResponse struct {
	Code     string `json:"code"`
	Token    string `json:"token"`
	PlayerID string `json:"playerId"`
}

func (h *Handlers) CreateGame(w http.ResponseWriter, r *http.Request) {
	var req createGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, "invalid request", http.StatusBadRequest)
		return
	}
	req.PlayerName = strings.TrimSpace(req.PlayerName)
	if req.PlayerName == "" {
		httpError(w, "name required", http.StatusBadRequest)
		return
	}
	if len(req.PlayerName) > 20 {
		req.PlayerName = req.PlayerName[:20]
	}

	host := game.NewHumanPlayer(req.PlayerName)
	send := func(playerID string, msg game.OutgoingMessage) {
		h.Registry.SendTo(playerID, msg)
	}

	var gameCode string
	broadcast := func(msg game.OutgoingMessage) {
		h.Registry.BroadcastToGame(gameCode, msg)
	}

	g := h.Hub.CreateGame(host, send, broadcast, h.AI)
	gameCode = g.State.Code
	log.Printf("[api] game created code=%s host=%q", gameCode, req.PlayerName)

	json.NewEncoder(w).Encode(createGameResponse{
		Code:     g.State.Code,
		Token:    host.Token,
		PlayerID: host.ID,
	})
}

type joinGameRequest struct {
	Code       string `json:"code"`
	PlayerName string `json:"playerName"`
}

type joinGameResponse struct {
	Token    string `json:"token"`
	PlayerID string `json:"playerId"`
}

func (h *Handlers) JoinGame(w http.ResponseWriter, r *http.Request) {
	var req joinGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, "invalid request", http.StatusBadRequest)
		return
	}
	req.PlayerName = strings.TrimSpace(req.PlayerName)
	req.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	if req.PlayerName == "" || req.Code == "" {
		httpError(w, "name and code required", http.StatusBadRequest)
		return
	}
	if len(req.PlayerName) > 20 {
		req.PlayerName = req.PlayerName[:20]
	}

	g := h.Hub.GetGame(req.Code)
	if g == nil {
		httpError(w, "game not found", http.StatusNotFound)
		return
	}

	player := game.NewHumanPlayer(req.PlayerName)
	g.HandleJoin(player)
	log.Printf("[api] player %q joined game %s", req.PlayerName, req.Code)

	json.NewEncoder(w).Encode(joinGameResponse{
		Token:    player.Token,
		PlayerID: player.ID,
	})
}

func httpError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

package ws

import (
	"context"
	"drawl/internal/game"
	"drawl/internal/hub"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

type Handler struct {
	Hub      *hub.Hub
	Registry *ClientRegistry
}

func NewHandler(h *hub.Hub, registry *ClientRegistry) *Handler {
	return &Handler{Hub: h, Registry: registry}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow all origins in dev
	})
	if err != nil {
		log.Printf("websocket accept error: %v", err)
		return
	}
	conn.SetReadLimit(10 * 1024 * 1024) // 10MB for drawing data

	token := r.URL.Query().Get("token")
	gameCode := r.URL.Query().Get("game")

	if token == "" || gameCode == "" {
		conn.Close(websocket.StatusPolicyViolation, "missing token or game code")
		return
	}

	g := h.Hub.GetGame(gameCode)
	if g == nil {
		conn.Close(websocket.StatusPolicyViolation, "game not found")
		return
	}

	// Find player by token
	player := g.State.FindPlayerByToken(token)
	if player == nil {
		conn.Close(websocket.StatusPolicyViolation, "invalid token")
		return
	}

	client := NewClient(player.ID, gameCode, conn)
	h.Registry.Add(client)
	defer func() {
		h.Registry.Remove(player.ID)
		g.HandleDisconnect(player.ID)
		client.Close()
		if g.IsEmpty() {
			h.Hub.RemoveGame(gameCode)
		}
	}()

	// Send initial game state
	client.Send(game.OutgoingMessage{
		Type: game.MsgGameState,
		Data: map[string]interface{}{
			"code":        g.State.Code,
			"phase":       g.State.Phase,
			"players":     g.State.Players,
			"round":       g.State.Round,
			"totalRounds": g.State.TotalRounds,
			"hostId":      g.State.HostID,
			"playerId":    player.ID,
		},
	})

	ctx := r.Context()
	for {
		msg, err := client.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil || websocket.CloseStatus(err) != -1 {
				return
			}
			if context.Cause(ctx) != nil {
				return
			}
			log.Printf("read error from %s: %v", player.ID, err)
			return
		}
		g.HandleMessage(player.ID, msg)
	}
}

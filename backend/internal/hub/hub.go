package hub

import (
	"drawl/internal/game"
	"sync"
	"time"
)

type Hub struct {
	mu    sync.RWMutex
	games map[string]*game.Game
}

func New() *Hub {
	h := &Hub{
		games: make(map[string]*game.Game),
	}
	go h.cleanupLoop()
	return h
}

func (h *Hub) CreateGame(host *game.Player, send game.SendFunc, broadcast game.BroadcastFunc, ai game.AIHandler) *game.Game {
	h.mu.Lock()
	defer h.mu.Unlock()

	var code string
	for {
		code = game.GenerateCode()
		if _, exists := h.games[code]; !exists {
			break
		}
	}

	g := game.NewGame(code, host, send, broadcast, ai)
	h.games[code] = g
	return g
}

func (h *Hub) GetGame(code string) *game.Game {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.games[code]
}

func (h *Hub) RemoveGame(code string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.games, code)
}

func (h *Hub) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		h.mu.Lock()
		for code, g := range h.games {
			if g.IsEmpty() {
				delete(h.games, code)
			}
		}
		h.mu.Unlock()
	}
}

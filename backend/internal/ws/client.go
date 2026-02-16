package ws

import (
	"context"
	"drawl/internal/game"
	"sync"

	"nhooyr.io/websocket"
)

type Client struct {
	PlayerID string
	GameCode string
	conn     *websocket.Conn
	mu       sync.Mutex
}

func NewClient(playerID, gameCode string, conn *websocket.Conn) *Client {
	return &Client{
		PlayerID: playerID,
		GameCode: gameCode,
		conn:     conn,
	}
}

func (c *Client) Send(msg game.OutgoingMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, err := EncodeMessage(msg)
	if err != nil {
		return err
	}
	return c.conn.Write(context.Background(), websocket.MessageText, data)
}

func (c *Client) ReadMessage(ctx context.Context) (game.IncomingMessage, error) {
	_, data, err := c.conn.Read(ctx)
	if err != nil {
		return game.IncomingMessage{}, err
	}
	return DecodeMessage(data)
}

func (c *Client) Close() {
	c.conn.Close(websocket.StatusNormalClosure, "")
}

// ClientRegistry tracks all connected clients.
type ClientRegistry struct {
	mu      sync.RWMutex
	clients map[string]*Client // playerID -> client
}

func NewClientRegistry() *ClientRegistry {
	return &ClientRegistry{
		clients: make(map[string]*Client),
	}
}

func (cr *ClientRegistry) Add(c *Client) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.clients[c.PlayerID] = c
}

func (cr *ClientRegistry) Remove(playerID string) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	delete(cr.clients, playerID)
}

func (cr *ClientRegistry) Get(playerID string) *Client {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	return cr.clients[playerID]
}

func (cr *ClientRegistry) SendTo(playerID string, msg game.OutgoingMessage) {
	c := cr.Get(playerID)
	if c != nil {
		c.Send(msg)
	}
}

func (cr *ClientRegistry) BroadcastToGame(gameCode string, msg game.OutgoingMessage) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	for _, c := range cr.clients {
		if c.GameCode == gameCode {
			go c.Send(msg)
		}
	}
}

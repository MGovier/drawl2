package ws

import (
	"drawl/internal/game"
	"encoding/json"
)

func EncodeMessage(msg game.OutgoingMessage) ([]byte, error) {
	return json.Marshal(msg)
}

func DecodeMessage(data []byte) (game.IncomingMessage, error) {
	var msg game.IncomingMessage
	err := json.Unmarshal(data, &msg)
	return msg, err
}

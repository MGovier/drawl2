package ai

import "drawl/internal/game"

// Handler implements game.AIHandler using OpenAI APIs.
// When APIKey is empty, it acts as a no-op (returns placeholders).
type Handler struct {
	APIKey string
}

func NewHandler(apiKey string) game.AIHandler {
	if apiKey == "" {
		return &noopHandler{}
	}
	return &Handler{APIKey: apiKey}
}

type noopHandler struct{}

func (h *noopHandler) GuessDrawing(imageDataURL string) (string, error) {
	return "???", nil
}

func (h *noopHandler) DrawPrompt(prompt string) (string, error) {
	return "", nil
}

package ai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (h *Handler) DrawPrompt(prompt string) (string, error) {
	if prompt == "" || prompt == "???" {
		return "", nil
	}

	body := map[string]any{
		"model":      "gpt-image-1-mini",
		"prompt":     fmt.Sprintf("You are a human playing a Pictionary-style game drawing a basic illustration on a small whiteboard with only black, blue, red, yellow and green marker pens. You need to draw an image representing '%s'. Only show the finished drawing, not the whiteboard or pen/hand. Do not use text.", prompt),
		"n":          1,
		"size":       "1024x1024",
		"quality":    "low",
		"moderation": "low",
	}

	data, _ := json.Marshal(body)
	log.Printf("[ai/draw] POST images/generations model=gpt-image-1-mini prompt=%q", prompt)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ai/draw] OpenAI returned %d: %s", resp.StatusCode, string(respBody))
		return "", fmt.Errorf("openai returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data []struct {
			B64JSON string `json:"b64_json"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil || len(result.Data) == 0 {
		log.Printf("[ai/draw] failed to parse response: %s", string(respBody))
		return "", fmt.Errorf("failed to parse openai response: %s", string(respBody))
	}

	decoded, err := base64.StdEncoding.DecodeString(result.Data[0].B64JSON)
	if err != nil {
		return "", fmt.Errorf("decode b64: %w", err)
	}
	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(decoded)
	log.Printf("[ai/draw] generated image, size=%d bytes", len(dataURL))
	return dataURL, nil
}

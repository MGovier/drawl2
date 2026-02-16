package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (h *Handler) GuessDrawing(imageDataURL string) (string, error) {
	if imageDataURL == "" {
		return "???", nil
	}

	// Resize to limit tokens/cost on the vision API
	resized, err := resizeDataURL(imageDataURL)
	if err != nil {
		log.Printf("[ai/vision] resize failed, using original: %v", err)
		resized = imageDataURL
	}

	body := map[string]any{
		"model":             "gpt-5-mini",
		"max_output_tokens": 200,
		"input": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{
						"type": "input_text",
						"text": "This is a drawing from a Pictionary-like game. What is this a drawing of? Reply with just 1-3 words, no punctuation.",
					},
					{
						"type":      "input_image",
						"image_url": resized,
					},
				},
			},
		},
	}

	data, _ := json.Marshal(body)
	log.Printf("[ai/vision] POST responses model=gpt-5-mini image_size=%d", len(resized))

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "???", fmt.Errorf("openai request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ai/vision] OpenAI returned %d: %s", resp.StatusCode, string(respBody))
		return "???", fmt.Errorf("openai returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Output []struct {
			Type    string `json:"type"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("[ai/vision] failed to parse response: %s", string(respBody))
		return "???", fmt.Errorf("failed to parse openai response: %s", string(respBody))
	}

	// Extract text from the first message output
	var guess string
	for _, out := range result.Output {
		if out.Type == "message" {
			for _, c := range out.Content {
				if c.Type == "output_text" && c.Text != "" {
					guess = c.Text
					break
				}
			}
		}
		if guess != "" {
			break
		}
	}
	if guess == "" {
		log.Printf("[ai/vision] no text in response: %s", string(respBody))
		return "???", fmt.Errorf("no text in openai response")
	}
	log.Printf("[ai/vision] guessed: %q", guess)
	return guess, nil
}

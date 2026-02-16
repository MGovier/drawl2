package main

import (
	"drawl/internal/ai"
	"drawl/internal/api"
	"drawl/internal/hub"
	"drawl/internal/ws"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	aiKey := os.Getenv("OPENAI_API_KEY")
	aiHandler := ai.NewHandler(aiKey)
	if aiKey != "" {
		log.Printf("OpenAI API key configured (%d chars)", len(aiKey))
	} else {
		log.Printf("No OPENAI_API_KEY set — AI players will use placeholders")
	}

	h := hub.New()
	registry := ws.NewClientRegistry()
	wsHandler := ws.NewHandler(h, registry)
	handlers := &api.Handlers{Hub: h, Registry: registry, AI: aiHandler}
	router := api.NewRouter(h, registry, wsHandler, handlers)

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

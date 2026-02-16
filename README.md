# Drawl

A multiplayer drawing telephone game. Players take turns drawing prompts and guessing what others drew, forming chains that inevitably go off the rails.

## How It Works

1. **Lobby** — Host creates a game, shares the 5-letter code. Up to 8 players (human or AI bots) can join.
2. **Rounds** — Each player gets their own chain starting with a random word. Rounds alternate between drawing and guessing. Each round, players rotate to a different chain, so everyone contributes to every chain.
3. **Reveal & Voting** — After all rounds complete, the full chains are revealed. Players vote thumbs-up on chains that survived the telephone game (awarding a point to the chain owner) and pick a favourite drawing (bonus point to the artist).
4. **Play Again** — Host can restart from the lobby with scores preserved.

## Project Structure

```
backend/                  Go server (HTTP + WebSocket)
  cmd/server/             Entry point
  internal/
    ai/                   OpenAI integration (vision + image generation for AI bots)
    api/                  REST handlers, router, middleware
    game/                 Core game logic (state, rounds, chains, scoring)
    hub/                  Game room registry
    ws/                   WebSocket client management
frontend/                 React + TypeScript (Vite)
  src/
    components/           UI components
    hooks/                useGameState (reducer), useWebSocket
    lib/                  Protocol constants & types, canvas utils
    styles/               CSS
```

## Development

### Prerequisites

- Go 1.24+
- Node.js 20+

### Running locally

```sh
# Backend (port 8080)
cd backend
go run ./cmd/server

# Frontend (port 5173, proxies API/WS to backend)
cd frontend
npm install
npm run dev
```

Open `http://localhost:5173`.

### Environment

| Variable | Description |
|---|---|
| `PORT` | Backend port (default `8080`) |
| `OPENAI_API_KEY` | Enables AI bot players. Without it, bots submit placeholders. |

Copy `.env.example` or create `.env` in the project root.

### Building

```sh
cd backend && go build -o server ./cmd/server
cd frontend && npm run build   # outputs to frontend/dist/
```

## Protocol

Communication is over a single WebSocket per player. Messages are JSON `{ type, data }`.

**Client -> Server:** `start_game`, `submit_drawing`, `submit_guess`, `add_ai`, `kick_player`, `submit_votes`, `play_again`

**Server -> Client:** `game_state`, `player_joined`, `player_left`, `game_started`, `turn_start`, `turn_tick`, `waiting`, `round_complete`, `game_over`, `score_update`, `return_to_lobby`, `error`, `ai_error`

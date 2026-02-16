package game

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// Message types
const (
	// Client -> Server
	MsgJoin          = "join"
	MsgAddAI         = "add_ai"
	MsgStartGame     = "start_game"
	MsgSubmitDrawing  = "submit_drawing"
	MsgSubmitGuess   = "submit_guess"
	MsgKickPlayer    = "kick_player"

	MsgSubmitVotes = "submit_votes"
	MsgPlayAgain   = "play_again"

	// Server -> Client
	MsgGameState     = "game_state"
	MsgPlayerJoined  = "player_joined"
	MsgPlayerLeft    = "player_left"
	MsgGameStarted   = "game_started"
	MsgTurnStart     = "turn_start"
	MsgTurnTick      = "turn_tick"
	MsgWaiting       = "waiting"
	MsgRoundComplete = "round_complete"
	MsgGameOver      = "game_over"
	MsgAIError       = "ai_error"
	MsgError         = "error"
	MsgScoreUpdate   = "score_update"
	MsgReturnToLobby = "return_to_lobby"
)

type IncomingMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type OutgoingMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type SendFunc func(playerID string, msg OutgoingMessage)
type BroadcastFunc func(msg OutgoingMessage)

type AIHandler interface {
	GuessDrawing(imageDataURL string) (string, error)
	DrawPrompt(prompt string) (string, error)
}

type Game struct {
	mu        sync.Mutex
	State     *GameState
	send      SendFunc
	broadcast BroadcastFunc
	ai        AIHandler
	timer      *time.Timer
	tickCancel chan struct{}   // closed to stop the tick goroutine
	submitted  map[string]bool // tracks submissions per round
}

func NewGame(code string, host *Player, send SendFunc, broadcast BroadcastFunc, ai AIHandler) *Game {
	return &Game{
		State:     NewGameState(code, host),
		send:      send,
		broadcast: broadcast,
		ai:        ai,
		submitted: make(map[string]bool),
	}
}

func (g *Game) HandleMessage(playerID string, msg IncomingMessage) {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch msg.Type {
	case MsgAddAI:
		g.handleAddAI(playerID)
	case MsgStartGame:
		g.handleStartGame(playerID)
	case MsgSubmitDrawing:
		g.handleSubmitDrawing(playerID, msg.Data)
	case MsgSubmitGuess:
		g.handleSubmitGuess(playerID, msg.Data)
	case MsgKickPlayer:
		g.handleKickPlayer(playerID, msg.Data)
	case MsgSubmitVotes:
		g.handleSubmitVotes(playerID, msg.Data)
	case MsgPlayAgain:
		g.handlePlayAgain(playerID)
	}
}

func (g *Game) HandleJoin(player *Player) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.State.Phase != PhaseLobby {
		g.send(player.ID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "game already started"}})
		return
	}
	if g.State.PlayerCount() >= 8 {
		g.send(player.ID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "game is full"}})
		return
	}

	g.State.AddPlayer(player)
	g.broadcast(OutgoingMessage{Type: MsgPlayerJoined, Data: map[string]interface{}{
		"player": player,
	}})
	g.sendGameState(player.ID)
}

func (g *Game) HandleDisconnect(playerID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	p := g.State.FindPlayer(playerID)
	if p == nil {
		return
	}

	if g.State.Phase == PhasePlaying && !g.submitted[playerID] {
		// Auto-submit blank
		chainIdx, turnType := g.State.GetAssignment(p.Index)
		if turnType == TurnDraw {
			g.State.Chains[chainIdx].Entries = append(g.State.Chains[chainIdx].Entries, ChainEntry{
				PlayerID: playerID, Type: TurnDraw, Drawing: "",
			})
		} else {
			g.State.Chains[chainIdx].Entries = append(g.State.Chains[chainIdx].Entries, ChainEntry{
				PlayerID: playerID, Type: TurnGuess, Guess: "???",
			})
		}
		g.submitted[playerID] = true
		g.checkRoundComplete()
	}

	if g.State.Phase == PhaseReveal && !g.State.VotesSubmitted[playerID] {
		g.State.VotesSubmitted[playerID] = true
		g.State.Votes[playerID] = &PlayerVote{}
		g.checkAllVotesIn()
	}

	g.State.RemovePlayer(playerID)
	g.broadcast(OutgoingMessage{Type: MsgPlayerLeft, Data: map[string]interface{}{
		"playerId": playerID,
		"hostId":   g.State.HostID,
	}})
}

func (g *Game) sendGameState(playerID string) {
	g.send(playerID, OutgoingMessage{Type: MsgGameState, Data: map[string]interface{}{
		"code":        g.State.Code,
		"phase":       g.State.Phase,
		"players":     g.State.Players,
		"round":       g.State.Round,
		"totalRounds": g.State.TotalRounds,
		"hostId":      g.State.HostID,
		"scores":      g.State.Scores,
	}})
}

func (g *Game) handleAddAI(playerID string) {
	if playerID != g.State.HostID {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "only host can add AI"}})
		return
	}
	if g.State.PlayerCount() >= 8 {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "game is full"}})
		return
	}
	ai := NewAIPlayer()
	g.State.AddPlayer(ai)
	g.broadcast(OutgoingMessage{Type: MsgPlayerJoined, Data: map[string]interface{}{
		"player": ai,
	}})
}

func (g *Game) handleStartGame(playerID string) {
	if playerID != g.State.HostID {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "only host can start"}})
		return
	}
	if g.State.PlayerCount() < 2 {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "need at least 2 players"}})
		return
	}
	log.Printf("[game %s] starting with %d players", g.State.Code, g.State.PlayerCount())
	g.State.Phase = PhasePlaying
	g.State.InitChains()
	g.submitted = make(map[string]bool)

	g.broadcast(OutgoingMessage{Type: MsgGameStarted, Data: nil})
	g.startTurn()
}

func (g *Game) startTurn() {
	turnType := "draw"
	if g.State.Round%2 == 1 {
		turnType = "guess"
	}
	log.Printf("[game %s] round %d/%d starting (%s)", g.State.Code, g.State.Round+1, g.State.TotalRounds, turnType)

	g.submitted = make(map[string]bool)
	infos := g.State.GetTurnInfos()

	for _, info := range infos {
		p := g.State.FindPlayer(info.PlayerID)
		if p == nil {
			continue
		}

		turnData := map[string]interface{}{
			"round":      g.State.Round,
			"totalRounds": g.State.TotalRounds,
			"turnType":   info.TurnType,
			"prompt":     info.Prompt,
			"timeLimit":  g.State.TurnTime,
		}

		if p.Type == AIPlayer {
			go g.handleAITurn(info)
		} else {
			g.send(info.PlayerID, OutgoingMessage{Type: MsgTurnStart, Data: turnData})
		}
	}

	// Start turn timer
	g.startTimer()
}

func (g *Game) stopTimer() {
	if g.timer != nil {
		g.timer.Stop()
		g.timer = nil
	}
	if g.tickCancel != nil {
		close(g.tickCancel)
		g.tickCancel = nil
	}
}

func (g *Game) startTimer() {
	g.stopTimer()

	g.timer = time.AfterFunc(time.Duration(g.State.TurnTime)*time.Second, func() {
		g.mu.Lock()
		defer g.mu.Unlock()
		g.forceSubmitAll()
		g.checkRoundComplete()
	})

	// Tick every second — exits when tickCancel is closed
	done := make(chan struct{})
	g.tickCancel = done
	go func() {
		remaining := g.State.TurnTime
		for remaining > 0 {
			select {
			case <-done:
				return
			case <-time.After(1 * time.Second):
			}
			remaining--
			g.mu.Lock()
			if g.State.Phase != PhasePlaying {
				g.mu.Unlock()
				return
			}
			g.broadcast(OutgoingMessage{Type: MsgTurnTick, Data: map[string]int{"remaining": remaining}})
			g.mu.Unlock()
		}
	}()
}

func (g *Game) forceSubmitAll() {
	for _, p := range g.State.Players {
		if g.submitted[p.ID] {
			continue
		}
		chainIdx, turnType := g.State.GetAssignment(p.Index)
		if turnType == TurnDraw {
			g.State.Chains[chainIdx].Entries = append(g.State.Chains[chainIdx].Entries, ChainEntry{
				PlayerID: p.ID, Type: TurnDraw, Drawing: "",
			})
		} else {
			g.State.Chains[chainIdx].Entries = append(g.State.Chains[chainIdx].Entries, ChainEntry{
				PlayerID: p.ID, Type: TurnGuess, Guess: "???",
			})
		}
		g.submitted[p.ID] = true
	}
}

func (g *Game) handleAITurn(info TurnInfo) {
	g.mu.Lock()
	if g.ai == nil {
		// No AI handler, submit placeholder
		if info.TurnType == TurnDraw {
			g.State.SubmitDrawing(info.PlayerID, "")
		} else {
			g.State.SubmitGuess(info.PlayerID, "???")
		}
		g.submitted[info.PlayerID] = true
		g.checkRoundComplete()
		g.mu.Unlock()
		return
	}
	ai := g.ai
	g.mu.Unlock()

	player := g.State.FindPlayer(info.PlayerID)
	playerName := "unknown"
	if player != nil {
		playerName = player.Name
	}

	var result string
	var aiErr error
	if info.TurnType == TurnGuess {
		log.Printf("[game] AI %q guessing drawing (round %d, chain %d, prompt_size=%d)",
			playerName, g.State.Round, info.ChainIdx, len(info.Prompt))
		result, aiErr = ai.GuessDrawing(info.Prompt)
		if aiErr != nil {
			log.Printf("[game] AI %q guess failed: %v", playerName, aiErr)
			result = "???"
		}
	} else {
		log.Printf("[game] AI %q drawing prompt=%q (round %d, chain %d)",
			playerName, info.Prompt, g.State.Round, info.ChainIdx)
		result, aiErr = ai.DrawPrompt(info.Prompt)
		if aiErr != nil {
			log.Printf("[game] AI %q draw failed: %v", playerName, aiErr)
			result = ""
		}
	}

	g.mu.Lock()
	if aiErr != nil {
		g.broadcast(OutgoingMessage{Type: MsgAIError, Data: map[string]string{
			"playerName": playerName,
			"message":    fmt.Sprintf("AI player %q failed: %v", playerName, aiErr),
		}})
	}
	defer g.mu.Unlock()
	if info.TurnType == TurnDraw {
		g.State.SubmitDrawing(info.PlayerID, result)
	} else {
		g.State.SubmitGuess(info.PlayerID, result)
	}
	g.submitted[info.PlayerID] = true
	g.checkRoundComplete()
}

const maxDrawingBytes = 5 * 1024 * 1024 // 5MB max for drawing data URLs

type submitDrawingData struct {
	Drawing string `json:"drawing"`
}

func (g *Game) handleSubmitDrawing(playerID string, data json.RawMessage) {
	if g.submitted[playerID] {
		return
	}
	var d submitDrawingData
	if err := json.Unmarshal(data, &d); err != nil {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "invalid data"}})
		return
	}
	if len(d.Drawing) > maxDrawingBytes {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "drawing too large"}})
		return
	}
	if !g.State.SubmitDrawing(playerID, d.Drawing) {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "cannot submit drawing now"}})
		return
	}
	g.submitted[playerID] = true
	g.send(playerID, OutgoingMessage{Type: MsgWaiting, Data: nil})
	g.checkRoundComplete()
}

type submitGuessData struct {
	Guess string `json:"guess"`
}

func (g *Game) handleSubmitGuess(playerID string, data json.RawMessage) {
	if g.submitted[playerID] {
		return
	}
	var d submitGuessData
	if err := json.Unmarshal(data, &d); err != nil {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "invalid data"}})
		return
	}
	if !g.State.SubmitGuess(playerID, d.Guess) {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "cannot submit guess now"}})
		return
	}
	g.submitted[playerID] = true
	g.send(playerID, OutgoingMessage{Type: MsgWaiting, Data: nil})
	g.checkRoundComplete()
}

func (g *Game) checkRoundComplete() {
	if !g.State.AllSubmitted() {
		return
	}
	log.Printf("[game %s] round %d complete, all submitted", g.State.Code, g.State.Round+1)
	g.stopTimer()

	g.broadcast(OutgoingMessage{Type: MsgRoundComplete, Data: map[string]int{
		"round": g.State.Round,
	}})

	if g.State.AdvanceRound() {
		// Game over — enter reveal phase, prepare for voting
		g.State.Votes = make(map[string]*PlayerVote)
		g.State.VotesSubmitted = make(map[string]bool)

		// AI players auto-submit empty votes
		for _, p := range g.State.Players {
			if p.Type == AIPlayer {
				g.State.VotesSubmitted[p.ID] = true
				g.State.Votes[p.ID] = &PlayerVote{}
			}
		}

		g.broadcast(OutgoingMessage{Type: MsgGameOver, Data: map[string]interface{}{
			"chains": g.State.GetChains(),
			"scores": g.State.Scores,
		}})
		return
	}

	// Small delay before next turn
	time.AfterFunc(500*time.Millisecond, func() {
		g.mu.Lock()
		defer g.mu.Unlock()
		if g.State.Phase == PhasePlaying {
			g.startTurn()
		}
	})
}

func (g *Game) checkAllVotesIn() {
	for _, p := range g.State.HumanPlayers() {
		if !g.State.VotesSubmitted[p.ID] {
			return
		}
	}
	log.Printf("[game %s] all votes in", g.State.Code)

	// Tally success votes: each thumbs-up on a chain gives 1 point to the chain owner
	for _, vote := range g.State.Votes {
		for _, chainIdx := range vote.SuccessChains {
			if chainIdx >= 0 && chainIdx < len(g.State.Chains) {
				ownerID := g.State.Chains[chainIdx].OwnerID
				g.State.Scores[ownerID]++
			}
		}
	}

	// Tally favourite drawing votes: most-picked drawing's author gets bonus point
	favCounts := make(map[string]int)
	for _, vote := range g.State.Votes {
		if vote.FavDrawing != "" {
			favCounts[vote.FavDrawing]++
		}
	}
	bestKey := ""
	bestCount := 0
	for key, count := range favCounts {
		if count > bestCount {
			bestCount = count
			bestKey = key
		}
	}
	if bestKey != "" {
		var ci, ei int
		if _, err := fmt.Sscanf(bestKey, "%d:%d", &ci, &ei); err == nil {
			if ci >= 0 && ci < len(g.State.Chains) && ei >= 0 && ei < len(g.State.Chains[ci].Entries) {
				authorID := g.State.Chains[ci].Entries[ei].PlayerID
				g.State.Scores[authorID]++
			}
		}
	}

	g.broadcast(OutgoingMessage{Type: MsgScoreUpdate, Data: map[string]interface{}{
		"scores":     g.State.Scores,
		"favDrawing": bestKey,
	}})
}

type kickData struct {
	PlayerID string `json:"playerId"`
}

func (g *Game) handleKickPlayer(playerID string, data json.RawMessage) {
	if playerID != g.State.HostID {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "only host can kick"}})
		return
	}
	var d kickData
	if err := json.Unmarshal(data, &d); err != nil {
		return
	}
	if d.PlayerID == playerID {
		return // can't kick yourself
	}
	g.State.RemovePlayer(d.PlayerID)
	g.broadcast(OutgoingMessage{Type: MsgPlayerLeft, Data: map[string]interface{}{
		"playerId": d.PlayerID,
		"hostId":   g.State.HostID,
	}})
}

type submitVotesData struct {
	SuccessChains []int  `json:"successChains"`
	FavDrawing    string `json:"favDrawing"` // "chainIdx:entryIdx"
}

func (g *Game) handleSubmitVotes(playerID string, data json.RawMessage) {
	if g.State.Phase != PhaseReveal {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "not in reveal phase"}})
		return
	}
	if g.State.VotesSubmitted[playerID] {
		return
	}
	var d submitVotesData
	if err := json.Unmarshal(data, &d); err != nil {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "invalid data"}})
		return
	}
	g.State.Votes[playerID] = &PlayerVote{SuccessChains: d.SuccessChains, FavDrawing: d.FavDrawing}
	g.State.VotesSubmitted[playerID] = true
	g.send(playerID, OutgoingMessage{Type: MsgWaiting, Data: nil})
	g.checkAllVotesIn()
}

func (g *Game) handlePlayAgain(playerID string) {
	if playerID != g.State.HostID {
		g.send(playerID, OutgoingMessage{Type: MsgError, Data: map[string]string{"message": "only host can restart"}})
		return
	}
	log.Printf("[game %s] play again requested by host", g.State.Code)
	g.State.ResetForNewGame()
	g.submitted = make(map[string]bool)
	g.broadcast(OutgoingMessage{Type: MsgReturnToLobby, Data: map[string]interface{}{
		"players": g.State.Players,
		"scores":  g.State.Scores,
		"hostId":  g.State.HostID,
	}})
}

func (g *Game) IsEmpty() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	humanCount := 0
	for _, p := range g.State.Players {
		if p.Type == HumanPlayer {
			humanCount++
		}
	}
	return humanCount == 0
}

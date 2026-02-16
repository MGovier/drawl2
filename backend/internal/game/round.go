package game

// InitChains sets up chains for all players with random words.
func (gs *GameState) InitChains() {
	n := gs.PlayerCount()
	gs.TotalRounds = n
	gs.Round = 0
	words := RandomWords(n)
	gs.Chains = make([]*Chain, n)
	for i, p := range gs.Players {
		gs.Chains[i] = &Chain{
			OriginalWord: words[i],
			OwnerID:      p.ID,
		}
	}
}

// CurrentTurnInfo returns what each player should be doing this round.
type TurnInfo struct {
	PlayerID string
	ChainIdx int
	TurnType TurnType
	Prompt   string // word or drawing data URL to work from
}

func (gs *GameState) GetTurnInfos() []TurnInfo {
	var infos []TurnInfo
	for _, p := range gs.Players {
		chainIdx, turnType := gs.GetAssignment(p.Index)
		chain := gs.Chains[chainIdx]

		var prompt string
		if gs.Round == 0 {
			prompt = chain.OriginalWord
		} else {
			lastEntry := chain.Entries[len(chain.Entries)-1]
			if turnType == TurnDraw {
				prompt = lastEntry.Guess
			} else {
				prompt = lastEntry.Drawing
			}
		}

		infos = append(infos, TurnInfo{
			PlayerID: p.ID,
			ChainIdx: chainIdx,
			TurnType: turnType,
			Prompt:   prompt,
		})
	}
	return infos
}

// SubmitDrawing records a drawing submission for a player.
func (gs *GameState) SubmitDrawing(playerID string, drawing string) bool {
	p := gs.FindPlayer(playerID)
	if p == nil {
		return false
	}
	chainIdx, turnType := gs.GetAssignment(p.Index)
	if turnType != TurnDraw {
		return false
	}
	gs.Chains[chainIdx].Entries = append(gs.Chains[chainIdx].Entries, ChainEntry{
		PlayerID: playerID,
		Type:     TurnDraw,
		Drawing:  drawing,
	})
	return true
}

// SubmitGuess records a guess submission for a player.
func (gs *GameState) SubmitGuess(playerID string, guess string) bool {
	p := gs.FindPlayer(playerID)
	if p == nil {
		return false
	}
	chainIdx, turnType := gs.GetAssignment(p.Index)
	if turnType != TurnGuess {
		return false
	}
	gs.Chains[chainIdx].Entries = append(gs.Chains[chainIdx].Entries, ChainEntry{
		PlayerID: playerID,
		Type:     TurnGuess,
		Guess:    guess,
	})
	return true
}

// AllSubmitted checks if all players have submitted for the current round.
func (gs *GameState) AllSubmitted() bool {
	for _, p := range gs.Players {
		chainIdx, _ := gs.GetAssignment(p.Index)
		chain := gs.Chains[chainIdx]
		if len(chain.Entries) <= gs.Round {
			return false
		}
	}
	return true
}

// AdvanceRound moves to the next round. Returns true if game is over.
func (gs *GameState) AdvanceRound() bool {
	gs.Round++
	if gs.Round >= gs.TotalRounds {
		gs.Phase = PhaseReveal
		return true
	}
	return false
}

// GetChains returns all chains for reveal.
func (gs *GameState) GetChains() []*Chain {
	return gs.Chains
}

// HumanPlayers returns a list of human players.
func (gs *GameState) HumanPlayers() []*Player {
	var humans []*Player
	for _, p := range gs.Players {
		if p.Type == HumanPlayer {
			humans = append(humans, p)
		}
	}
	return humans
}

// ResetForNewGame resets game state for a new game while keeping players and scores.
func (gs *GameState) ResetForNewGame() {
	gs.Phase = PhaseLobby
	gs.Chains = nil
	gs.Round = 0
	gs.TotalRounds = 0
	gs.Votes = make(map[string]*PlayerVote)
	gs.VotesSubmitted = make(map[string]bool)
}

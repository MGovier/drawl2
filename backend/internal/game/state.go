package game

type GamePhase int

const (
	PhaseLobby GamePhase = iota
	PhasePlaying
	PhaseReveal
)

type TurnType int

const (
	TurnDraw TurnType = iota
	TurnGuess
)

type ChainEntry struct {
	PlayerID string   `json:"playerId"`
	Type     TurnType `json:"type"`
	Drawing  string   `json:"drawing,omitempty"` // base64 PNG data URL
	Guess    string   `json:"guess,omitempty"`
}

type Chain struct {
	OriginalWord string       `json:"originalWord"`
	OwnerID      string       `json:"ownerId"`
	Entries      []ChainEntry `json:"entries"`
}

type PlayerVote struct {
	SuccessChains []int  `json:"successChains"` // chain indices player thinks succeeded
	FavDrawing    string `json:"favDrawing"`    // "chainIdx:entryIdx" of favourite drawing
}

type GameState struct {
	Code        string    `json:"code"`
	Phase       GamePhase `json:"phase"`
	Players     []*Player `json:"players"`
	Chains      []*Chain  `json:"-"`
	Round       int       `json:"round"`
	TotalRounds int       `json:"totalRounds"`
	HostID      string    `json:"hostId"`
	TurnTime    int       `json:"turnTime"` // seconds per turn

	Scores         map[string]int         `json:"scores"`  // playerID → total points
	Votes          map[string]*PlayerVote `json:"-"`        // playerID → their votes at reveal
	VotesSubmitted map[string]bool        `json:"-"`        // tracks who has voted
}

func NewGameState(code string, host *Player) *GameState {
	host.Index = 0
	return &GameState{
		Code:            code,
		Phase:           PhaseLobby,
		Players:         []*Player{host},
		HostID:          host.ID,
		TurnTime:        60,
		Scores:         make(map[string]int),
		Votes:          make(map[string]*PlayerVote),
		VotesSubmitted: make(map[string]bool),
	}
}

func (gs *GameState) AddPlayer(p *Player) {
	p.Index = len(gs.Players)
	gs.Players = append(gs.Players, p)
}

func (gs *GameState) RemovePlayer(id string) {
	for i, p := range gs.Players {
		if p.ID == id {
			gs.Players = append(gs.Players[:i], gs.Players[i+1:]...)
			break
		}
	}
	// Re-index
	for i, p := range gs.Players {
		p.Index = i
	}
	// Promote host if needed
	if gs.HostID == id && len(gs.Players) > 0 {
		for _, p := range gs.Players {
			if p.Type == HumanPlayer {
				gs.HostID = p.ID
				break
			}
		}
	}
}

func (gs *GameState) FindPlayer(id string) *Player {
	for _, p := range gs.Players {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func (gs *GameState) FindPlayerByToken(token string) *Player {
	for _, p := range gs.Players {
		if p.Token == token {
			return p
		}
	}
	return nil
}

func (gs *GameState) PlayerCount() int {
	return len(gs.Players)
}

// GetAssignment returns (chainIndex, turnType) for a given player in the current round.
func (gs *GameState) GetAssignment(playerIdx int) (int, TurnType) {
	n := len(gs.Players)
	chainIdx := ((playerIdx - gs.Round) % n + n) % n
	turnType := TurnDraw
	if gs.Round%2 == 1 {
		turnType = TurnGuess
	}
	return chainIdx, turnType
}

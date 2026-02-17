package game

import "testing"

func setupGame(n int) *GameState {
	host := NewHumanPlayer("P0")
	gs := NewGameState("TEST1", host)
	for i := 1; i < n; i++ {
		gs.AddPlayer(NewHumanPlayer("P" + string(rune('0'+i))))
	}
	gs.InitChains()
	return gs
}

func TestInitChains(t *testing.T) {
	gs := setupGame(3)

	if len(gs.Chains) != 3 {
		t.Fatalf("Chains len = %d, want 3", len(gs.Chains))
	}
	if gs.Round != 0 {
		t.Errorf("Round = %d, want 0", gs.Round)
	}
	if gs.TotalRounds != 3 {
		t.Errorf("TotalRounds = %d, want 3", gs.TotalRounds)
	}
	for i, c := range gs.Chains {
		if c.OriginalWord == "" {
			t.Errorf("chain %d has empty word", i)
		}
		if c.OwnerID != gs.Players[i].ID {
			t.Errorf("chain %d owner = %q, want %q", i, c.OwnerID, gs.Players[i].ID)
		}
	}
}

func TestSubmitDrawing_Valid(t *testing.T) {
	gs := setupGame(3)
	// Round 0 is a draw round
	p := gs.Players[0]
	chainIdx, _ := gs.GetAssignment(p.Index)

	ok := gs.SubmitDrawing(p.ID, "data:image/png;base64,abc")
	if !ok {
		t.Error("SubmitDrawing returned false for valid submission")
	}
	if len(gs.Chains[chainIdx].Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(gs.Chains[chainIdx].Entries))
	}
	entry := gs.Chains[chainIdx].Entries[0]
	if entry.Drawing != "data:image/png;base64,abc" {
		t.Error("drawing content mismatch")
	}
	if entry.Type != TurnDraw {
		t.Error("entry type should be TurnDraw")
	}
}

func TestSubmitDrawing_WrongTurnType(t *testing.T) {
	gs := setupGame(3)
	gs.Round = 1 // guess round
	p := gs.Players[0]

	ok := gs.SubmitDrawing(p.ID, "data:image/png;base64,abc")
	if ok {
		t.Error("SubmitDrawing should return false on guess round")
	}
}

func TestSubmitDrawing_UnknownPlayer(t *testing.T) {
	gs := setupGame(3)
	ok := gs.SubmitDrawing("nonexistent", "data:image/png;base64,abc")
	if ok {
		t.Error("SubmitDrawing should return false for unknown player")
	}
}

func TestSubmitGuess_Valid(t *testing.T) {
	gs := setupGame(3)
	// Need entries for round 0 first, then advance to round 1
	for _, p := range gs.Players {
		gs.SubmitDrawing(p.ID, "drawing")
	}
	gs.Round = 1

	p := gs.Players[0]
	chainIdx, _ := gs.GetAssignment(p.Index)

	ok := gs.SubmitGuess(p.ID, "cat")
	if !ok {
		t.Error("SubmitGuess returned false for valid submission")
	}
	entry := gs.Chains[chainIdx].Entries[len(gs.Chains[chainIdx].Entries)-1]
	if entry.Guess != "cat" {
		t.Errorf("guess = %q, want %q", entry.Guess, "cat")
	}
}

func TestSubmitGuess_WrongTurnType(t *testing.T) {
	gs := setupGame(3)
	// Round 0 is draw
	ok := gs.SubmitGuess(gs.Players[0].ID, "cat")
	if ok {
		t.Error("SubmitGuess should return false on draw round")
	}
}

func TestAllSubmitted(t *testing.T) {
	gs := setupGame(3)

	if gs.AllSubmitted() {
		t.Error("AllSubmitted should be false before any submissions")
	}

	for _, p := range gs.Players {
		gs.SubmitDrawing(p.ID, "drawing")
	}

	if !gs.AllSubmitted() {
		t.Error("AllSubmitted should be true after all submissions")
	}
}

func TestAdvanceRound(t *testing.T) {
	gs := setupGame(3)

	gs.Round = 0
	done := gs.AdvanceRound()
	if done {
		t.Error("AdvanceRound should return false when rounds remain")
	}
	if gs.Round != 1 {
		t.Errorf("Round = %d, want 1", gs.Round)
	}

	gs.Round = gs.TotalRounds - 1
	done = gs.AdvanceRound()
	if !done {
		t.Error("AdvanceRound should return true at final round")
	}
	if gs.Phase != PhaseReveal {
		t.Errorf("Phase = %d, want PhaseReveal", gs.Phase)
	}
}

func TestGetTurnInfos_Round0(t *testing.T) {
	gs := setupGame(3)

	infos := gs.GetTurnInfos()
	if len(infos) != 3 {
		t.Fatalf("infos len = %d, want 3", len(infos))
	}
	for _, info := range infos {
		if info.TurnType != TurnDraw {
			t.Errorf("round 0 should be TurnDraw")
		}
		// Prompt should be the original word
		chain := gs.Chains[info.ChainIdx]
		if info.Prompt != chain.OriginalWord {
			t.Errorf("prompt = %q, want %q", info.Prompt, chain.OriginalWord)
		}
	}
}

func TestGetTurnInfos_Round1(t *testing.T) {
	gs := setupGame(3)

	// Submit drawings for round 0
	for _, p := range gs.Players {
		gs.SubmitDrawing(p.ID, "drawing-by-"+p.ID)
	}
	gs.Round = 1

	infos := gs.GetTurnInfos()
	for _, info := range infos {
		if info.TurnType != TurnGuess {
			t.Errorf("round 1 should be TurnGuess")
		}
		// Prompt should be the drawing from round 0
		chain := gs.Chains[info.ChainIdx]
		if info.Prompt != chain.Entries[0].Drawing {
			t.Errorf("prompt = %q, want %q", info.Prompt, chain.Entries[0].Drawing)
		}
	}
}

func TestHumanPlayers(t *testing.T) {
	host := NewHumanPlayer("Alice")
	gs := NewGameState("TEST1", host)
	gs.AddPlayer(NewAIPlayer())
	gs.AddPlayer(NewHumanPlayer("Bob"))

	humans := gs.HumanPlayers()
	if len(humans) != 2 {
		t.Errorf("HumanPlayers len = %d, want 2", len(humans))
	}
	for _, h := range humans {
		if h.Type != HumanPlayer {
			t.Error("HumanPlayers returned a non-human")
		}
	}
}

func TestResetForNewGame(t *testing.T) {
	gs := setupGame(3)
	gs.Phase = PhaseReveal
	gs.Round = 2
	gs.Scores["p1"] = 5

	gs.ResetForNewGame()

	if gs.Phase != PhaseLobby {
		t.Errorf("Phase = %d, want PhaseLobby", gs.Phase)
	}
	if gs.Chains != nil {
		t.Error("Chains should be nil")
	}
	if gs.Round != 0 {
		t.Errorf("Round = %d, want 0", gs.Round)
	}
	// Scores should be preserved
	if gs.Scores["p1"] != 5 {
		t.Errorf("Scores should be preserved, got %d", gs.Scores["p1"])
	}
	// Players should be preserved
	if len(gs.Players) != 3 {
		t.Errorf("Players len = %d, want 3", len(gs.Players))
	}
}

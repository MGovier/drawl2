package game

import "testing"

func TestNewGameState(t *testing.T) {
	host := NewHumanPlayer("Alice")
	gs := NewGameState("ABCDE", host)

	if gs.Code != "ABCDE" {
		t.Errorf("Code = %q, want %q", gs.Code, "ABCDE")
	}
	if gs.Phase != PhaseLobby {
		t.Errorf("Phase = %d, want PhaseLobby", gs.Phase)
	}
	if gs.HostID != host.ID {
		t.Errorf("HostID = %q, want %q", gs.HostID, host.ID)
	}
	if len(gs.Players) != 1 {
		t.Fatalf("Players len = %d, want 1", len(gs.Players))
	}
	if gs.Players[0].Index != 0 {
		t.Errorf("host Index = %d, want 0", gs.Players[0].Index)
	}
	if gs.TurnTime != 60 {
		t.Errorf("TurnTime = %d, want 60", gs.TurnTime)
	}
}

func TestAddPlayer(t *testing.T) {
	host := NewHumanPlayer("Alice")
	gs := NewGameState("ABCDE", host)

	p2 := NewHumanPlayer("Bob")
	gs.AddPlayer(p2)

	if len(gs.Players) != 2 {
		t.Fatalf("Players len = %d, want 2", len(gs.Players))
	}
	if p2.Index != 1 {
		t.Errorf("p2 Index = %d, want 1", p2.Index)
	}

	p3 := NewHumanPlayer("Charlie")
	gs.AddPlayer(p3)
	if p3.Index != 2 {
		t.Errorf("p3 Index = %d, want 2", p3.Index)
	}
}

func TestRemovePlayer_Reindex(t *testing.T) {
	host := NewHumanPlayer("Alice")
	gs := NewGameState("ABCDE", host)
	p2 := NewHumanPlayer("Bob")
	p3 := NewHumanPlayer("Charlie")
	gs.AddPlayer(p2)
	gs.AddPlayer(p3)

	gs.RemovePlayer(p2.ID)

	if len(gs.Players) != 2 {
		t.Fatalf("Players len = %d, want 2", len(gs.Players))
	}
	if gs.Players[0].Index != 0 || gs.Players[1].Index != 1 {
		t.Errorf("indices = [%d, %d], want [0, 1]", gs.Players[0].Index, gs.Players[1].Index)
	}
}

func TestRemovePlayer_HostPromotion(t *testing.T) {
	host := NewHumanPlayer("Alice")
	gs := NewGameState("ABCDE", host)
	p2 := NewHumanPlayer("Bob")
	gs.AddPlayer(p2)

	gs.RemovePlayer(host.ID)

	if gs.HostID != p2.ID {
		t.Errorf("HostID = %q, want %q (promoted)", gs.HostID, p2.ID)
	}
}

func TestRemovePlayer_HostPromotionSkipsAI(t *testing.T) {
	host := NewHumanPlayer("Alice")
	gs := NewGameState("ABCDE", host)
	ai := NewAIPlayer()
	gs.AddPlayer(ai)
	human := NewHumanPlayer("Bob")
	gs.AddPlayer(human)

	gs.RemovePlayer(host.ID)

	if gs.HostID != human.ID {
		t.Errorf("HostID = %q, want %q (should skip AI)", gs.HostID, human.ID)
	}
}

func TestFindPlayer(t *testing.T) {
	host := NewHumanPlayer("Alice")
	gs := NewGameState("ABCDE", host)

	found := gs.FindPlayer(host.ID)
	if found == nil || found.ID != host.ID {
		t.Errorf("FindPlayer(%q) = %v, want player", host.ID, found)
	}
	if gs.FindPlayer("nonexistent") != nil {
		t.Error("FindPlayer(nonexistent) should return nil")
	}
}

func TestFindPlayerByToken(t *testing.T) {
	host := NewHumanPlayer("Alice")
	gs := NewGameState("ABCDE", host)

	found := gs.FindPlayerByToken(host.Token)
	if found == nil || found.ID != host.ID {
		t.Errorf("FindPlayerByToken returned wrong player")
	}
	if gs.FindPlayerByToken("bad-token") != nil {
		t.Error("FindPlayerByToken(bad-token) should return nil")
	}
}

func TestGetAssignment(t *testing.T) {
	host := NewHumanPlayer("Alice")
	gs := NewGameState("ABCDE", host)
	gs.AddPlayer(NewHumanPlayer("Bob"))
	gs.AddPlayer(NewHumanPlayer("Charlie"))
	n := len(gs.Players)

	// Round 0: draw turn, each player gets a different chain
	gs.Round = 0
	seen := make(map[int]bool)
	for i := 0; i < n; i++ {
		chainIdx, turnType := gs.GetAssignment(i)
		if turnType != TurnDraw {
			t.Errorf("round 0: turnType = %d, want TurnDraw", turnType)
		}
		if chainIdx < 0 || chainIdx >= n {
			t.Errorf("round 0: chainIdx = %d out of range", chainIdx)
		}
		seen[chainIdx] = true
	}
	if len(seen) != n {
		t.Errorf("round 0: expected %d unique chains, got %d", n, len(seen))
	}

	// Round 1: guess turn
	gs.Round = 1
	seen = make(map[int]bool)
	for i := 0; i < n; i++ {
		chainIdx, turnType := gs.GetAssignment(i)
		if turnType != TurnGuess {
			t.Errorf("round 1: turnType = %d, want TurnGuess", turnType)
		}
		seen[chainIdx] = true
	}
	if len(seen) != n {
		t.Errorf("round 1: expected %d unique chains, got %d", n, len(seen))
	}

	// Rotation: no player gets the same chain in consecutive rounds
	for i := 0; i < n; i++ {
		gs.Round = 0
		c0, _ := gs.GetAssignment(i)
		gs.Round = 1
		c1, _ := gs.GetAssignment(i)
		if c0 == c1 {
			t.Errorf("player %d got same chain %d in rounds 0 and 1", i, c0)
		}
	}
}

package api

import (
	"bytes"
	"drawl/internal/hub"
	"drawl/internal/ws"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestHandlers(password string) *Handlers {
	return &Handlers{
		Hub:          hub.New(),
		Registry:     ws.NewClientRegistry(),
		GamePassword: password,
	}
}

func TestCreateGame_Success(t *testing.T) {
	h := newTestHandlers("")
	body := `{"playerName":"Alice"}`
	req := httptest.NewRequest("POST", "/api/games", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateGame(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp createGameResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Code == "" {
		t.Error("Code should not be empty")
	}
	if resp.Token == "" {
		t.Error("Token should not be empty")
	}
	if resp.PlayerID == "" {
		t.Error("PlayerID should not be empty")
	}
}

func TestCreateGame_MissingName(t *testing.T) {
	h := newTestHandlers("")
	body := `{"playerName":"  "}`
	req := httptest.NewRequest("POST", "/api/games", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateGame(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestCreateGame_WrongPassword(t *testing.T) {
	h := newTestHandlers("secret")
	body := `{"playerName":"Alice","password":"wrong"}`
	req := httptest.NewRequest("POST", "/api/games", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateGame(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestCreateGame_CorrectPassword(t *testing.T) {
	h := newTestHandlers("secret")
	body := `{"playerName":"Alice","password":"secret"}`
	req := httptest.NewRequest("POST", "/api/games", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateGame(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

func TestJoinGame_Success(t *testing.T) {
	h := newTestHandlers("")

	// Create a game first
	createBody := `{"playerName":"Alice"}`
	createReq := httptest.NewRequest("POST", "/api/games", bytes.NewBufferString(createBody))
	createW := httptest.NewRecorder()
	h.CreateGame(createW, createReq)

	var createResp createGameResponse
	json.NewDecoder(createW.Body).Decode(&createResp)

	// Join the game
	joinBody, _ := json.Marshal(joinGameRequest{Code: createResp.Code, PlayerName: "Bob"})
	joinReq := httptest.NewRequest("POST", "/api/games/join", bytes.NewBuffer(joinBody))
	joinW := httptest.NewRecorder()

	h.JoinGame(joinW, joinReq)

	if joinW.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", joinW.Code)
	}
	var joinResp joinGameResponse
	json.NewDecoder(joinW.Body).Decode(&joinResp)
	if joinResp.Token == "" {
		t.Error("Token should not be empty")
	}
}

func TestJoinGame_MissingFields(t *testing.T) {
	h := newTestHandlers("")
	body := `{"code":"","playerName":""}`
	req := httptest.NewRequest("POST", "/api/games/join", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.JoinGame(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestJoinGame_GameNotFound(t *testing.T) {
	h := newTestHandlers("")
	body := `{"code":"ZZZZZ","playerName":"Bob"}`
	req := httptest.NewRequest("POST", "/api/games/join", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.JoinGame(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestCreateGame_InvalidJSON(t *testing.T) {
	h := newTestHandlers("")
	req := httptest.NewRequest("POST", "/api/games", bytes.NewBufferString("not json"))
	w := httptest.NewRecorder()

	h.CreateGame(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

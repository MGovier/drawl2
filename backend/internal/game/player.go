package game

import (
	"math/rand"

	"github.com/google/uuid"
)

type PlayerType int

const (
	HumanPlayer PlayerType = iota
	AIPlayer
)

type Player struct {
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	Type  PlayerType `json:"type"`
	Token string     `json:"-"`
	Index int        `json:"index"`
}

func NewHumanPlayer(name string) *Player {
	token := uuid.New().String()
	return &Player{
		ID:    token[:8],
		Name:  name,
		Type:  HumanPlayer,
		Token: token,
	}
}

var botAdjectives = []string{
	"Sneaky", "Fuzzy", "Wobbly", "Sparkly", "Grumpy",
	"Zippy", "Dizzy", "Chunky", "Spooky", "Bouncy",
	"Sassy", "Cranky", "Giggly", "Snooty", "Wacky",
	"Funky", "Jolly", "Zappy", "Loopy", "Cheeky",
}

var botNouns = []string{
	"Penguin", "Noodle", "Potato", "Cactus", "Walrus",
	"Muffin", "Pickle", "Badger", "Waffle", "Goblin",
	"Turnip", "Otter", "Pretzel", "Squid", "Toaster",
	"Llama", "Pancake", "Moose", "Nugget", "Dingus",
}

func randomBotName() string {
	adj := botAdjectives[rand.Intn(len(botAdjectives))]
	noun := botNouns[rand.Intn(len(botNouns))]
	return adj + " " + noun
}

func NewAIPlayer() *Player {
	id := uuid.New().String()[:8]
	return &Player{
		ID:   id,
		Name: randomBotName(),
		Type: AIPlayer,
	}
}

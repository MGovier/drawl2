package game

import "testing"

func TestRandomWords_Count(t *testing.T) {
	for _, n := range []int{1, 3, 8} {
		words := RandomWords(n)
		if len(words) != n {
			t.Errorf("RandomWords(%d) returned %d words", n, len(words))
		}
	}
}

func TestRandomWords_NonEmpty(t *testing.T) {
	words := RandomWords(10)
	for i, w := range words {
		if w == "" {
			t.Errorf("RandomWords: word %d is empty", i)
		}
	}
}

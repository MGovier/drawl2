package game

import (
	"strings"
	"testing"
)

func TestGenerateCode_Length(t *testing.T) {
	code := GenerateCode()
	if len(code) != codeLength {
		t.Errorf("GenerateCode() len = %d, want %d", len(code), codeLength)
	}
}

func TestGenerateCode_CharacterSet(t *testing.T) {
	for i := 0; i < 100; i++ {
		code := GenerateCode()
		for _, c := range code {
			if !strings.ContainsRune(codeChars, c) {
				t.Errorf("code %q contains invalid char %q", code, c)
			}
		}
	}
}

func TestGenerateCode_Uniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code := GenerateCode()
		seen[code] = true
	}
	// With 30^5 ≈ 24M possibilities, 100 codes should all be unique
	if len(seen) < 95 {
		t.Errorf("only %d unique codes out of 100, expected near-100%% uniqueness", len(seen))
	}
}

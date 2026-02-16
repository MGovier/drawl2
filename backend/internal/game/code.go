package game

import (
	"math/rand"
	"strings"
)

const codeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const codeLength = 5

func GenerateCode() string {
	var b strings.Builder
	for i := 0; i < codeLength; i++ {
		b.WriteByte(codeChars[rand.Intn(len(codeChars))])
	}
	return b.String()
}

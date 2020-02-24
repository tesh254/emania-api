package helpers

import (
	"math/rand"
	"time"
)

const charSet = "abcdefghijklmnopqrstuvwxyz" +
"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seedRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// StringWithCharset generates a random string based on length and charset
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charSet[seedRand.Intn(len(charset))]
	}

	return string(b)
}

// String returns the generated string
func String(length int) string {
	return StringWithCharset(length, charSet)
}
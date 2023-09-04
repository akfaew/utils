package webhandler

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

var (
	DisableMeditation = false
)

func meditation() string {
	if DisableMeditation { // for tests
		return "06210bea27acd4fa"
	}
	r := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, r); err != nil {
		return "I can't even meditate"
	}
	return hex.EncodeToString(r)
}

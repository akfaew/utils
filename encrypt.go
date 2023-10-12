package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

type Key []byte

//nolint:unused
func genkey() string {
	randbytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, randbytes); err != nil {
		log.Fatal(err)
	}
	// randbytes := []byte("01234567890123456789012345678901")

	return hex.EncodeToString(randbytes)
}

func ParseKey(enckey string) (ret Key) {
	if newkey, err := hex.DecodeString(enckey); err != nil {
		log.Fatal(err)
	} else if len(newkey) != 32 {
		log.Fatalf("Encoded key must be of length 32, has length: %d", len(newkey))
	} else {
		return newkey
	}

	return nil // this will never be reached
}

func (key Key) Encrypt(text string) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("error: %v", err)
	}

	return gcm.Seal(nonce, nonce, []byte(text), nil)
}

func (key Key) Decrypt(ciphertext []byte) string {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		log.Fatalf("error: len(ciphertext) < nonceSize: %d < %d",
			len(ciphertext), nonceSize)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return string(plaintext)
}

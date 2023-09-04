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

func NewKey(enckey string) (ret Key) {
	if newkey, err := hex.DecodeString(enckey); err != nil {
		log.Fatal(err)
	} else if len(newkey) != 32 {
		log.Fatalf("Encoded key must be of length 32, has length: %d", len(newkey))
	} else {
		return newkey
	}

	return nil
}

func (key Key) Encrypt(text []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, Errorc(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, Errorc(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, Errorc(err)
	}

	return gcm.Seal(nonce, nonce, text, nil), nil
}

func (key Key) Decrypt(ciphertext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, Errorc(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, Errorc(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, Errorc(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, Errorc(err)
	}

	return plaintext, nil
}

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	key := ParseKey("303132333435363738393031323" +
		"3343536373839303132333435363738393031")

	t.Run("sometext", func(t *testing.T) {
		plain := []byte("sometext")

		enc := key.Encrypt(plain)
		dec := key.Decrypt(enc)

		require.Equal(t, plain, dec)
	})

	t.Run("empty", func(t *testing.T) {
		plain := []byte("")

		enc := key.Encrypt(plain)
		dec := key.Decrypt(enc)

		require.Equal(t, "", string(dec))
	})

	t.Run("very empty", func(t *testing.T) {
		dec := key.Decrypt([]byte{})

		require.Equal(t, "", string(dec))
	})
}

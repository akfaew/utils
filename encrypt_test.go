package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	key := NewKey("303132333435363738393031323" +
		"3343536373839303132333435363738393031")

	plain := []byte("sometext")

	enc, err := key.Encrypt(plain)
	require.NoError(t, err)

	dec, err := key.Decrypt(enc)
	require.NoError(t, err)

	require.Equal(t, dec, plain)
}

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompress(t *testing.T) {
	t.Run("sometext", func(t *testing.T) {
		plain := "sometext"

		comp := Compress(plain)
		var dec string
		Decompress(comp, &dec)

		require.Equal(t, dec, plain)
	})

	t.Run("empty", func(t *testing.T) {
		type Dupa struct {
			A string
			B int
		}
		plain := Dupa{A: "a", B: 2}

		comp := Compress(plain)
		var dec Dupa
		Decompress(comp, &dec)

		require.Equal(t, dec, plain)
	})

	t.Run("very empty", func(t *testing.T) {
		var (
			plain *string
			dec   *string
		)
		plain = nil

		comp := Compress(plain)
		Decompress(comp, &dec)

		require.Equal(t, dec, plain)
	})
}

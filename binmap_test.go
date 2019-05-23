package utils

import (
	"testing"

	"github.com/akfaew/test"
)

func TestBitMap(t *testing.T) {
	var bm BitMap
	test.False(t, bm.IsSet(1))
	test.False(t, bm.IsSet(2))
	test.False(t, bm.IsSet(4))

	// Set something
	bm.Set(2)
	test.False(t, bm.IsSet(1))
	test.True(t, bm.IsSet(2))
	test.False(t, bm.IsSet(4))

	// Again
	bm.Set(2)
	test.False(t, bm.IsSet(1))
	test.True(t, bm.IsSet(2))
	test.False(t, bm.IsSet(4))

	// Set something else
	bm.Set(4)
	test.False(t, bm.IsSet(1))
	test.True(t, bm.IsSet(2))
	test.True(t, bm.IsSet(4))

	// Clear something
	bm.Clear(2)
	test.False(t, bm.IsSet(1))
	test.False(t, bm.IsSet(2))
	test.True(t, bm.IsSet(4))

	// Clear everything
	bm.Clear(4)
	test.False(t, bm.IsSet(1))
	test.False(t, bm.IsSet(2))
	test.False(t, bm.IsSet(4))
}

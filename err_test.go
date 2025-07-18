package utils

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorc(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.True(t, Errorc(nil) == nil)
	})

	t.Run("path", func(t *testing.T) {
		err := Errorc(errors.New("oups"))
		t.Logf("err.Error() = %s", err.Error())
		require.Regexp(t, `utils/err_test.go:\d+ oups`, err.Error())
	})

	t.Run("is", func(t *testing.T) {
		err := Errorc(os.ErrNotExist)
		require.True(t, errors.Is(err, os.ErrNotExist))
	})

	t.Run("nice", func(t *testing.T) {
		err := UserErrorfc(os.ErrNotExist, "Something bad happened %d", 5)
		require.True(t, errors.Is(err, os.ErrNotExist))
	})
}

func TestRootCause(t *testing.T) {
	x := errors.New("x")
	a := Errorc(x)
	b := Errorc(a)
	require.Equal(t, x, RootCause(b))
}

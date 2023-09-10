package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Init(t *testing.T) {
	Init("/a/b/c/d/e.go")
	require.Equal(t, "/a/b/c/d/", trimprefix)

	Init("")
	file, _ := logctx(0)
	require.Equal(t, file, "utils_test.go")
}

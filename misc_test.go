package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/akfaew/utils/fixture"
)

func TestSlash(t *testing.T) {
	tests := []struct {
		input string
		ret1  string
		ret2  string
	}{
		{"", "", ""},
		{"test", "test", ""},
		{"test/", "test", ""},
		{"test///", "test", "//"},
		{"test/case", "test", "case"},
		{"test/case/a/b/c", "test", "case/a/b/c"},
		{"/case/a/b/c", "", "case/a/b/c"},
	}

	for _, tc := range tests {
		ret1, ret2 := Slash(tc.input)
		require.Equal(t, ret1, tc.ret1)
		require.Equal(t, ret2, tc.ret2)
	}
}

func TestRandEmail(t *testing.T) {
	val := RandEmail()
	require.Len(t, val, 16)
}

func TestCrc32(t *testing.T) {
	for _, s := range []string{"a", "", "://!%$"} {
		t.Run(Crc32(s), func(t *testing.T) {
			fixture.Fixture(t, fmt.Sprintf("%s -> %s\n", s, Crc32(s)))
		})
	}
}

func TestUniq(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		input := []int{5, 3, 2, 3, 1, 2, 5}
		original := append([]int(nil), input...)
		output := Uniq(input)
		require.Equal(t, []int{1, 2, 3, 5}, output)
		require.Equal(t, original, input)
	})

	t.Run("string", func(t *testing.T) {
		input := []string{"b", "a", "b", "c", "a"}
		output := Uniq(input)
		require.Equal(t, []string{"a", "b", "c"}, output)
	})

	t.Run("nil", func(t *testing.T) {
		var input []int
		output := Uniq(input)
		require.Nil(t, output)
	})

	t.Run("empty", func(t *testing.T) {
		input := []int{}
		output := Uniq(input)
		require.Len(t, output, 0)
	})
}

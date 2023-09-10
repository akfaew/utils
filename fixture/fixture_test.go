package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_makeFixturePath(t *testing.T) {
	require.Equal(t, makeFixturePath(t, ""), "testdata/output/Test_makeFixturePath.fixture")

	t.Run("Sub Test", func(t *testing.T) {
		require.Equal(t, makeFixturePath(t, ""), "testdata/output/Test_makeFixturePath-Sub_Test.fixture")
	})

	t.Run("Sub Test With Extra", func(t *testing.T) {
		require.Equal(t, makeFixturePath(t, "extra"), "testdata/output/Test_makeFixturePath-Sub_Test_With_Extra-extra.fixture")
	})
}

func Test_Fixture(t *testing.T) {
	t.Run("bytes", func(t *testing.T) {
		Fixture(t, []byte("an array of bytes"))
	})

	t.Run("string", func(t *testing.T) {
		Fixture(t, "a string of text")
	})

	t.Run("struct", func(t *testing.T) {
		Fixture(t, struct {
			A string
			B int
		}{
			"something",
			1234,
		})
	})

	r := *regen
	t.Run("regen", func(t *testing.T) {
		b := []byte(fmt.Sprintf("%v", time.Now()))

		*regen = true
		require.True(t, Regen())
		Fixture(t, b)

		*regen = false
		require.False(t, Regen())
		Fixture(t, b)

		os.Remove(makeFixturePath(t, ""))
	})
	*regen = r
}

func Test_InputFixture(t *testing.T) {
	input := InputFixture(t, "input.fixture")
	require.Equal(t, string(input), "foo")
}

func Test_InputFixtureJson(t *testing.T) {
	a := struct {
		A string
		B string
		C int
	}{
		"aaa", "bbb", 123,
	}

	b := struct {
		A string
		B string
		C int
	}{}
	InputFixtureJson(t, "struct.json", &b)
	require.Equal(t, a, b)
}

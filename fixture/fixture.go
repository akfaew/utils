package test

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	permissions = 0644
)

var (
	regen = flag.Bool("regen", false, "Regenerate fixtures")

	FixtureInputPath  = "testdata/input/"
	FixtureOutputPath = "testdata/output/"
)

func Regen() bool {
	return *regen
}

// makeFixturePath makes a path from the test name, and optionally appends "extra".
func makeFixturePath(t *testing.T, extra string) string {
	t.Helper()

	name := strings.Replace(t.Name(), "/", "-", -1)
	path := FixtureOutputPath + name
	if extra != "" {
		path += "-" + extra
	}
	path += ".fixture"

	return path
}

// Fixture ensures that 'data' is equal to what's stored on disk.
//
// If 'data' is a string it gets written verbatim, otherwise it's json-encoded.
//
// The filename of the fixture is generated from the test name. To use multiple fixtures in one test see FixtureExtra()
func Fixture(t *testing.T, data interface{}) {
	t.Helper()

	FixtureExtra(t, "", data)
}

// FixtureExtra ensures that data is equal to what's stored on disk.
//
// If 'data' is a string it gets written verbatim, otherwise it's json-encoded.
//
// The filename of the fixture is generated from the test name with 'extra' appended.
func FixtureExtra(t *testing.T, extra string, data interface{}) {
	t.Helper()

	// Write strings verbatim, otherwise json-encode.
	var got []byte
	if b, ok := data.(string); ok {
		got = []byte(b)
	} else {
		var err error
		got, err = json.MarshalIndent(data, "", "\t")
		require.NoError(t, err)
	}

	path := makeFixturePath(t, extra)
	// If -regen then write and return
	if *regen {
		if err := os.WriteFile(path, []byte(got), permissions); err != nil {
			t.Fatalf("Error writing file %q: %v", path, err)
		}
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Error reading file %q: %v", path, err)
	}

	if !bytes.Equal(got, want) {
		if err := os.WriteFile("/tmp/got", got, permissions); err != nil {
			t.Fatalf("Error writing file /tmp/got: %v", err)
		}
		t.Fatalf("Error comparing with fixture. See: diff /tmp/got %s", path)
	}
}

// InputFixture returns the contents of a fixture file
func InputFixture(t *testing.T, filename string) []byte {
	t.Helper()

	input, err := os.ReadFile(FixtureInputPath + filename)
	if err != nil {
		t.Fatalf("Error reading fixture: %v", err)
	}

	return input
}

// InputFixtureJson returns the contents of a json fixture file, and unmarshals it
func InputFixtureJson(t *testing.T, filename string, v interface{}) {
	t.Helper()

	data := InputFixture(t, filename)
	require.NoError(t, json.Unmarshal(data, v))
}

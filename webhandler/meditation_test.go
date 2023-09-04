package webhandler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMeditation(t *testing.T) {
	m := meditation()
	t.Logf("meditation=%s, len=%d", m, len(m))
	require.Len(t, m, 16)
}

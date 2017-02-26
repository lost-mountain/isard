package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractNames(t *testing.T) {
	cases := []struct {
		n   string
		cn  string
		san []string
	}{
		{"www.cabal.io", "cabal.io", []string{"cabal.io", "www.cabal.io"}},
		{"cabal.io", "cabal.io", []string{"cabal.io", "www.cabal.io"}},
		{"test.cabal.io", "test.cabal.io", []string{"test.cabal.io"}},
	}

	for _, c := range cases {
		g, err := ExtractNames(c.n)
		require.NoError(t, err)
		require.Equal(t, c.cn, g.CN)
		require.Equal(t, c.san, g.SAN)
	}
}

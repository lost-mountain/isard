package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadHeader(t *testing.T) {
	v, err := readHeader("netlify.com", "netlify.com", "Server")
	require.NoError(t, err)
	require.Equal(t, "Netlify", v)
}

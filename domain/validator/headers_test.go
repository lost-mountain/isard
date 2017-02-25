package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidHeader(t *testing.T) {
	v := ValidHeader("netlify.com", "Server", "Netlify")
	require.True(t, v)

	v = ValidHeader("netlify.com", "Server", "nginx")
	require.False(t, v)
}

func TestReadHeader(t *testing.T) {
	v, err := readHeader("netlify.com", "netlify.com", "Server")
	require.NoError(t, err)
	require.Equal(t, "Netlify", v)
}

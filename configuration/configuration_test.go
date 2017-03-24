package configuration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	c, err := Load("testdata/basic.json")
	require.NoError(t, err)

	require.Equal(t, "/tmp/cert.pem", c.TLS.CertFile)
	require.Equal(t, "/tmp/key.pem", c.TLS.KeyFile)
	require.Equal(t, 8080, c.TCP.Port)
	require.Equal(t, productionDirectory, c.ACME.DefaultProductionDirectory)
	require.Equal(t, stagingDirectory, c.ACME.DefaultStagingDirectory)
}

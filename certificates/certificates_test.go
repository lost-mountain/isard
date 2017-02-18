package certificates

import (
	"os"
	"testing"

	"github.com/lost-mountain/isard/account"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	directoryURL string
}

func (s testSuite) TestRegisterAccount() {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(s.T(), err)
	a.DirectoryURL = s.directoryURL

	c, err := NewClient(a)
	require.NoError(s.T(), err)
	err = c.Register()
	require.NoError(s.T(), err)
}

func TestCertificates(t *testing.T) {
	directoryURL := os.Getenv("ISARD_TEST_ACME_DIRECTORY")
	if directoryURL == "" {
		t.Skip(`Certificates test suite skipped.
Set ISARD_TEST_ACME_DIRECTORY with the ACME test directory URL to enable them`)
	}

	suite.Run(t, &testSuite{
		directoryURL: directoryURL,
	})
}

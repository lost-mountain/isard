package certificates

import (
	"os"
	"testing"

	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/domain"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	directoryURL string
}

func (s testSuite) TestRegisterAccount() {
	c := s.newClient()

	err := c.Register()
	require.NoError(s.T(), err)
}

func (s testSuite) TestAuthorizeDomain() {
	c := s.newClient()
	d := domain.NewDomain(c.account, "example.com")

	err := c.Register()
	require.NoError(s.T(), err)

	authz, err := c.AuthorizeDomain(d)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), authz.URI)
}

func (s testSuite) newClient() *Client {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(s.T(), err)
	a.DirectoryURL = s.directoryURL

	c, err := NewClient(a)
	require.NoError(s.T(), err)

	return c
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

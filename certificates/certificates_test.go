package certificates

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"golang.org/x/crypto/acme"

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
	c := s.newClient("")

	err := c.Register()
	require.NoError(s.T(), err)
}

func (s testSuite) TestAuthorizeDomain() {
	c := s.newClient("")
	d := domain.NewDomain(c.account, "example.com")

	err := c.Register()
	require.NoError(s.T(), err)

	authz, err := c.AuthorizeDomain(d)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), authz.URI)
}

func (s testSuite) TestCompleteChallenge() {
	httpTestListenerAddr := os.Getenv("ISARD_TEST_HTTP_LISTENER_ADDR")
	if httpTestListenerAddr == "" {
		s.T().Skip(`HTTP complete challenge test suite skipped.
		Set ISARD_TEST_HTTP_LISTENER_ADDR with the listener address (i.e: 127.0.0.1:5002) to enable them`)
	}
	c := s.newClient("")
	d := domain.NewDomain(c.account, "test-cert-isard.cabal.io")

	err := c.Register()
	require.NoError(s.T(), err)

	authz, err := c.AuthorizeDomain(d)
	require.NoError(s.T(), err)

	var chal *acme.Challenge
	for _, c := range authz.Challenges {
		if c.Type == "http-01" {
			chal = c
			break
		}
	}

	require.NotNil(s.T(), chal)

	d, err = c.PrepareChallenge(d, chal)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), d.HTTP01ChallengePath)
	require.NotEmpty(s.T(), d.HTTP01ChallengeResponse)

	mux := http.NewServeMux()
	mux.HandleFunc(d.HTTP01ChallengePath, func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, d.HTTP01ChallengeResponse)
	})
	ts := httptest.NewUnstartedServer(mux)
	l, err := net.Listen("tcp", httpTestListenerAddr)
	require.NoError(s.T(), err)
	ts.Listener = l
	defer ts.Close()
	ts.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.client.WaitAuthorization(ctx, authz.URI)
	require.NoError(s.T(), err)
}

func (s testSuite) newClient(apiKey string) *Client {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(s.T(), err)
	a.DirectoryURL = s.directoryURL

	c, err := NewClientWithAPIKey(a, apiKey)
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

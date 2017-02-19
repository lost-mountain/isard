package challenges

import (
	"os"
	"testing"

	"golang.org/x/crypto/acme"

	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/domain"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type dnsTestSuite struct {
	suite.Suite
	resolver *DNSResolver
}

func (s *dnsTestSuite) TearDownSuite() {
	s.resolver.ns1Client.Zones.Delete("test-dns-isard.cabal.io")
}

func (s *dnsTestSuite) TestResolveAndCleanup() {
	d := &domain.Domain{
		Name: "cabal.io",
	}

	chal := &acme.Challenge{
		Type:  "dns-01",
		Token: "123==",
	}

	_, err := s.resolver.Resolve(d, chal)
	require.NoError(s.T(), err)

	err = s.resolver.Cleanup(d)
	require.NoError(s.T(), err)
}

func (s *dnsTestSuite) TestResolveAndCleanupWithMissingZone() {
	d := &domain.Domain{
		Name: "test-dns-isard.cabal.io",
	}

	chal := &acme.Challenge{
		Type:  "dns-01",
		Token: "123==",
	}

	_, err := s.resolver.Resolve(d, chal)
	require.NoError(s.T(), err)

	err = s.resolver.Cleanup(d)
	require.NoError(s.T(), err)
}

func TestDNSResolver(t *testing.T) {
	apiKey := os.Getenv("ISARD_TEST_NS1_API_KEY")
	if apiKey == "" {
		t.Skip(`DNSResolver test suite skipped.
Set ISARD_TEST_NS1_API_KEY with the NS1 api key to enable them`)
	}

	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(t, err)

	pk, err := a.PrivateKey()
	require.NoError(t, err)

	c := &acme.Client{
		Key: pk,
	}

	suite.Run(t, &dnsTestSuite{
		resolver: NewDNSResolver(apiKey, c),
	})
}

package challenges

import (
	"os"
	"testing"

	"golang.org/x/crypto/acme"

	"github.com/stretchr/testify/suite"
)

type dnsTestSuite struct {
	suite.Suite
	resolver *DNSResolver
}

func (s *dnsTestSuite) TestResolve() {
}

func (s *dnsTestSuite) TestCleanup() {
}

func TestDNSResolver(t *testing.T) {
	apiKey := os.Getenv("ISARD_TEST_NS1_API_KEY")
	if apiKey == "" {
		t.Skip(`DNSResolver test suite skipped.
Set ISARD_TEST_NS1_API_KEY with the NS1 api key to enable them`)
	}

	suite.Run(t, &dnsTestSuite{
		resolver: NewDNSResolver(apiKey, &acme.Client{}),
	})
}

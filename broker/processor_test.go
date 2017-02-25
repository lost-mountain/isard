package broker

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/configuration"
	"github.com/lost-mountain/isard/domain"
	"github.com/lost-mountain/isard/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type noopBroker struct{}

func (b *noopBroker) Close() error                                       { return nil }
func (b *noopBroker) Publish(topic TopicType, payload interface{}) error { return nil }
func (b *noopBroker) Subscribe(processor Processor) error                { return nil }

type testSuite struct {
	suite.Suite
	processor *DomainProcessor
	account   *account.Account
}

func (s *testSuite) TestCreateDomain() {
	s.createDefaultDomain("test.cabal.io")
}

func (s *testSuite) TestValidateDomain() {
	s.createDefaultDomain("netlify.com")

	c := &DomainPayload{
		AccountID:  s.account.ID,
		DomainName: "netlify.com",
	}

	m := NewMessage(c)

	err := s.processor.ValidateDomain(m)
	require.NoError(s.T(), err)

	d, err := s.processor.db.GetDomain(s.account.ID, "netlify.com")
	require.NoError(s.T(), err)
	require.Equal(s.T(), domain.Verified, d.State)
}

func (s *testSuite) TestValidateDomainWithInvalidHeaders() {
	s.createDefaultDomain("invalid.cabal.io")

	c := &DomainPayload{
		AccountID:  s.account.ID,
		DomainName: "invalid.cabal.io",
	}

	m := NewMessage(c)

	err := s.processor.ValidateDomain(m)
	require.EqualError(s.T(), err, "domain validation failed for domain: invalid.cabal.io")
}

func (s *testSuite) createDefaultDomain(name string) {
	c := &CreateDomainPayload{
		AccountID:    s.account.ID,
		AccountToken: s.account.Token,
		DomainName:   name,
	}

	m := NewMessage(c)

	err := s.processor.CreateDomain(m)
	require.NoError(s.T(), err)
}

func TestProcessor(t *testing.T) {
	f, err := ioutil.TempFile("", "isard-")
	require.NoError(t, err)

	defer os.Remove(f.Name())
	err = f.Close()
	require.NoError(t, err)

	b, err := storage.NewBoltBucket(&storage.Configuration{
		URL: f.Name(),
	})
	require.NoError(t, err)

	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(t, err)

	err = b.SaveAccount(a)
	require.NoError(t, err)

	c := &configuration.Configuration{}
	c.Domains.HeaderValidator.Name = "Server"
	c.Domains.HeaderValidator.Value = "Netlify"

	p := &DomainProcessor{
		db:     b,
		broker: &noopBroker{},
		config: c.Domains,
	}

	s := &testSuite{
		processor: p,
		account:   a,
	}

	suite.Run(t, s)
}

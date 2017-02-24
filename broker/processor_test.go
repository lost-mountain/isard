package broker

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/decoding"
	"github.com/lost-mountain/isard/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	processor Processor
	account   *account.Account
}

func (s *testSuite) TestCreateDomain() {
	c := &CreateDomainPayload{
		AccountID:    s.account.ID,
		AccountToken: s.account.Token,
		DomainName:   "test.cabal.io",
	}
	b, err := json.Marshal(c)
	require.NoError(s.T(), err)

	m := NewMessage(bytes.NewReader(b))

	err = s.processor.CreateDomain(m)
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

	p := &DomainProcessor{
		db:      b,
		decoder: decoding.JSONDecoder,
	}

	s := &testSuite{
		processor: p,
		account:   a,
	}

	suite.Run(t, s)
}

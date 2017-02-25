package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	bucket Bucket
}

func (s *testSuite) TestGetAccount() {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(s.T(), err)

	err = s.bucket.SaveAccount(a)
	require.NoError(s.T(), err)

	acc, err := s.bucket.GetAccount(a.ID, a.Token)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), a.ID, acc.ID)

	acc, err = s.bucket.GetAccount(a.ID, uuid.New())
	require.Nil(s.T(), acc)
	require.Error(s.T(), err, "unable to get account with invalid token")
}

func (s *testSuite) TestGetDomain() {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), a.ID)

	d, err := domain.NewDomain(a, "example.com")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), a.ID)
	err = s.bucket.SaveDomain(d)
	require.NoError(s.T(), err)

	dom, err := s.bucket.GetDomain(a.ID, "example.com")
	require.NoError(s.T(), err)
	require.Equal(s.T(), d.ID, dom.ID)

	_, err = s.bucket.GetDomain(a.ID, "foobar.com")
	require.Error(s.T(), err, "unable to get domain with missing name")
}

func (s *testSuite) TestSaveAccount() {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), a.ID)
	require.NotEmpty(s.T(), a.Token)

	err = s.bucket.SaveAccount(a)
	require.NoError(s.T(), err)
}

func (s *testSuite) TestSaveDomain() {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), a.ID)

	d, err := domain.NewDomain(a, "example.com")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), a.ID)
	err = s.bucket.SaveDomain(d)
	require.NoError(s.T(), err)
}

func TestBoltBucket(t *testing.T) {
	f, err := ioutil.TempFile("", "isard-")
	require.NoError(t, err)

	defer os.Remove(f.Name())
	err = f.Close()
	require.NoError(t, err)

	b, err := NewBoltBucket(&Configuration{
		URL: f.Name(),
	})
	require.NoError(t, err)

	suite.Run(t, &testSuite{bucket: b})
}

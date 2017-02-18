package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/lost-mountain/isard/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	bucket Bucket
}

func (s *Suite) TestSaveAccount() {
	a := account.NewAccount("david.calavera@gmail.com")
	require.NotEmpty(s.T(), a.ID)
	require.NotEmpty(s.T(), a.Token)

	err := s.bucket.SaveAccount(a)
	require.NoError(s.T(), err)
}

func (s *Suite) TestGetAccount() {
	a := account.NewAccount("david.calavera@gmail.com")
	err := s.bucket.SaveAccount(a)
	require.NoError(s.T(), err)

	acc, err := s.bucket.GetAccount(a.ID, a.Token)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), a.ID, acc.ID)

	acc, err = s.bucket.GetAccount(a.ID, uuid.New())
	require.Nil(s.T(), acc)
	require.Error(s.T(), err, "unable to get account with a fake token")
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

	suite.Run(t, &Suite{bucket: b})
}

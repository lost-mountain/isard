package domain

import (
	"testing"

	"github.com/lost-mountain/isard/account"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	account *account.Account
}

func (s *testSuite) TestNewDomain() {
	d, err := NewDomain(s.account, "test.cabal.io")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), d.ID)
	require.NotEmpty(s.T(), d.CreatedAt)
	require.NotEmpty(s.T(), d.UpdatedAt)
	require.Equal(s.T(), "test.cabal.io", d.Name)
	require.Equal(s.T(), Pending, d.State)
	require.Equal(s.T(), "http-01", d.ChallengeType)

	names := d.SANNames()
	require.Len(s.T(), names, 1)
	require.Equal(s.T(), []string{"test.cabal.io"}, names)
}

func (s *testSuite) TestNewDomainWithChallengeType() {
	d, err := NewDomainWithChallengeType(s.account, "test.cabal.io", "dns-01")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), d.ID)
	require.NotEmpty(s.T(), d.CreatedAt)
	require.NotEmpty(s.T(), d.UpdatedAt)
	require.Equal(s.T(), "test.cabal.io", d.Name)
	require.Equal(s.T(), Pending, d.State)
	require.Equal(s.T(), "dns-01", d.ChallengeType)

	names := d.SANNames()
	require.Len(s.T(), names, 1)
	require.Equal(s.T(), []string{"test.cabal.io"}, names)
}

func (s *testSuite) TestNewDomainEmptyWithChallengeType() {
	d, err := NewDomainWithChallengeType(s.account, "test.cabal.io", "")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), d.ID)
	require.NotEmpty(s.T(), d.CreatedAt)
	require.NotEmpty(s.T(), d.UpdatedAt)
	require.Equal(s.T(), "test.cabal.io", d.Name)
	require.Equal(s.T(), Pending, d.State)
	require.Equal(s.T(), "http-01", d.ChallengeType)

	names := d.SANNames()
	require.Len(s.T(), names, 1)
	require.Equal(s.T(), []string{"test.cabal.io"}, names)
}

func (s *testSuite) TestNewDomainWithExtension() {
	d, err := NewDomainWithChallengeType(s.account, "cabal.io", "")
	require.NoError(s.T(), err)
	require.Equal(s.T(), "cabal.io", d.Name)

	names := d.SANNames()
	require.Len(s.T(), names, 2)
	require.Contains(s.T(), names, "cabal.io")
	require.Contains(s.T(), names, "www.cabal.io")

	d, err = NewDomainWithChallengeType(s.account, "www.cabal.io", "")
	require.NoError(s.T(), err)
	require.Equal(s.T(), "cabal.io", d.Name)

	names = d.SANNames()
	require.Len(s.T(), names, 2)
	require.Contains(s.T(), names, "cabal.io")
	require.Contains(s.T(), names, "www.cabal.io")
}

func (s *testSuite) TestAddSANName() {
	d, err := NewDomain(s.account, "test.cabal.io")
	require.NoError(s.T(), err)

	err = d.AddSANName("beta.cabal.io")
	require.NoError(s.T(), err)

	names := d.SANNames()
	require.Len(s.T(), names, 2)
	require.Contains(s.T(), names, "beta.cabal.io")

	err = d.AddSANName("beta.cabal.io")
	require.EqualError(s.T(), err, ErrDuplicatedSANName.Error())
}

func (s *testSuite) TestRemoveSANName() {
	d, err := NewDomain(s.account, "test.cabal.io")
	require.NoError(s.T(), err)

	err = d.AddSANName("beta.cabal.io")
	require.NoError(s.T(), err)

	d.RemoveSANName("beta.cabal.io")
	names := d.SANNames()
	require.Len(s.T(), names, 1)
	require.Contains(s.T(), names, "test.cabal.io")
}

func TestDomain(t *testing.T) {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(t, err)

	s := &testSuite{
		account: a,
	}

	suite.Run(t, s)
}

package domain

import (
	"testing"

	"github.com/lost-mountain/isard/account"
	"github.com/stretchr/testify/require"
)

func TestNewDomain(t *testing.T) {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(t, err)

	d, err := NewDomain(a, "test.cabal.io")
	require.NoError(t, err)
	require.NotEmpty(t, d.ID)
	require.NotEmpty(t, d.CreatedAt)
	require.NotEmpty(t, d.UpdatedAt)
	require.Equal(t, "cabal.io", d.Name)
	require.Equal(t, Pending, d.State)
	require.Equal(t, "http-01", d.ChallengeType)

	names := d.SANNames()
	require.Len(t, names, 1)
	require.Equal(t, "test.cabal.io", names[0])
}

func TestNewDomainWithChallengeType(t *testing.T) {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(t, err)

	d, err := NewDomainWithChallengeType(a, "test.cabal.io", "dns-01")
	require.NoError(t, err)
	require.NotEmpty(t, d.ID)
	require.NotEmpty(t, d.CreatedAt)
	require.NotEmpty(t, d.UpdatedAt)
	require.Equal(t, "cabal.io", d.Name)
	require.Equal(t, Pending, d.State)
	require.Equal(t, "dns-01", d.ChallengeType)

	names := d.SANNames()
	require.Len(t, names, 1)
	require.Equal(t, "test.cabal.io", names[0])
}

func TestAddSANName(t *testing.T) {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(t, err)

	d, err := NewDomain(a, "test.cabal.io")
	require.NoError(t, err)

	err = d.AddSANName("beta.cabal.io")
	require.NoError(t, err)

	names := d.SANNames()
	require.Len(t, names, 2)
	require.Contains(t, names, "beta.cabal.io")

	err = d.AddSANName("beta.cabal.io")
	require.EqualError(t, err, ErrDuplicatedSANName.Error())
}

func TestRemoveSANName(t *testing.T) {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(t, err)

	d, err := NewDomain(a, "test.cabal.io")
	require.NoError(t, err)

	err = d.AddSANName("beta.cabal.io")
	require.NoError(t, err)

	d.RemoveSANName("beta.cabal.io")
	require.Len(t, d.SANNames(), 1)

	d.RemoveSANName("test.cabal.io")
	require.Len(t, d.SANNames(), 1)
}

package domain

import (
	"testing"

	"github.com/lost-mountain/isard/account"
	"github.com/stretchr/testify/require"
)

func TestNewDomain(t *testing.T) {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(t, err)

	d := NewDomain(a, "test.cabal.io")
	require.NotEmpty(t, d.ID)
	require.NotEmpty(t, d.CreatedAt)
	require.NotEmpty(t, d.UpdatedAt)
	require.Equal(t, "test.cabal.io", d.Name)
	require.Equal(t, Pending, d.State)
	require.Equal(t, "http-01", d.ChallengeType)
}

func TestNewDomainWithChallengeType(t *testing.T) {
	a, err := account.NewAccount("david.calavera@gmail.com")
	require.NoError(t, err)

	d := NewDomainWithChallengeType(a, "test.cabal.io", "dns-01")
	require.NotEmpty(t, d.ID)
	require.NotEmpty(t, d.CreatedAt)
	require.NotEmpty(t, d.UpdatedAt)
	require.Equal(t, "test.cabal.io", d.Name)
	require.Equal(t, Pending, d.State)
	require.Equal(t, "dns-01", d.ChallengeType)
}

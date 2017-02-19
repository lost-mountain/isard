package domain

import (
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestListARecords(t *testing.T) {
	answers, err := listARecords("dollsanddoughnuts.com")
	require.NoError(t, err)
	require.Len(t, answers, 1)

	ar, ok := answers[0].(*dns.A)
	require.True(t, ok, "expected A record: %v", ar)

	require.Equal(t, "104.198.14.52", ar.A.String())
}

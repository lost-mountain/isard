package domain

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

// Names is a struct that holds
// domain names information.
type Names struct {
	CN  string
	SAN []string
}

// ExtractNames normalizes a domain
// name to extract the common name
// and the SAN names.
func ExtractNames(name string) (*Names, error) {
	normal := strings.TrimLeft(name, "www.")

	dn, err := publicsuffix.Domain(name)
	if err != nil {
		return nil, errors.Wrap(err, "error looking up the correct domain name")
	}

	if normal == dn {
		return &Names{
			CN:  dn,
			SAN: []string{dn, "www." + dn},
		}, nil
	}

	return &Names{
		CN:  name,
		SAN: []string{name},
	}, nil
}

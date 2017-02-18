package challenges

import (
	"github.com/lost-mountain/isard/domain"
	"golang.org/x/crypto/acme"
)

// Resolver decides how to act upon an ACME challenge.
type Resolver interface {
	Cleanup(d *domain.Domain, keyAuth string) error
	Resolve(d *domain.Domain, chal *acme.Challenge) error
}

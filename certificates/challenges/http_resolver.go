package challenges

import (
	"github.com/lost-mountain/isard/domain"
	"github.com/pkg/errors"
	"golang.org/x/crypto/acme"
)

// HTTPResolver validates http-01 challenges.
type HTTPResolver struct {
	acmeClient *acme.Client
}

// Cleanup is a NOOP for the HTTPResolver.
func (r *HTTPResolver) Cleanup(d *domain.Domain) error {
	return nil
}

// Resolve stores the HTTP path and response challenges in the domain.
// So it can pass it along when the ACME verification is triggered.
func (r *HTTPResolver) Resolve(d *domain.Domain, challenge *acme.Challenge) (*domain.Domain, error) {
	d.HTTP01ChallengePath = r.acmeClient.HTTP01ChallengePath(challenge.Token)
	res, err := r.acmeClient.HTTP01ChallengeResponse(challenge.Token)
	if err != nil {
		return nil, errors.Wrapf(err, "error generating response for http-01 challenge: %s", d.Name)
	}

	d.HTTP01ChallengeResponse = res
	return d, nil
}

// NewHTTPResolver initializes a new http challenge resolver.
func NewHTTPResolver(ac *acme.Client) *HTTPResolver {
	return &HTTPResolver{
		acmeClient: ac,
	}
}

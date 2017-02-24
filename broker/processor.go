package broker

import (
	"github.com/lost-mountain/isard/decoding"
	"github.com/lost-mountain/isard/domain"
	"github.com/lost-mountain/isard/storage"
	"github.com/pkg/errors"
)

// DomainProcessor controls the lifecycle of a domain.
// It moves the domain across the state machine
// accordingly to the previous operation and its result.
type DomainProcessor struct {
	db      storage.Bucket
	decoder decoding.Decoder
}

// CreateDomain creates a new domain.
// If the job succeeds, it moves the domain to the
// validation state. Otherwise, it leaves to the broker
// to decide what to do with the message.
func (p *DomainProcessor) CreateDomain(m *Message) error {
	c := &CreateDomainPayload{}
	if err := p.decoder(c, m.Payload); err != nil {
		return errors.Wrap(err, "error creating domain")
	}

	a, err := p.db.GetAccount(c.AccountID, c.AccountToken)
	if err != nil {
		return err
	}

	d, err := domain.NewDomainWithChallengeType(a, c.DomainName, c.ChallengeType)
	if err != nil {
		return err
	}

	return p.db.SaveDomain(d)
}

// ModifyDomain modifies a domain.
// If the job succeeds, it moves the domain to the
// validation state. Otherwise, it leaves to the broker
// to decide what to do with the message.
func (p *DomainProcessor) ModifyDomain(*Message) error { return nil }

// ValidateDomain validates that a domain is correctly configured
// before authorizing its issuing.
// If the job succeeds, it moves the domain to the
// authorization state. Otherwise, it leaves to the broker
// to decide what to do with the message.
func (p *DomainProcessor) ValidateDomain(*Message) error { return nil }

// AuthorizeDomain sends an authorization request to the CA.
// If the job succeeds, it moves the domain to the
// certificate request state. Otherwise, it leaves to the broker
// to decide what to do with the message.
func (p *DomainProcessor) AuthorizeDomain(*Message) error { return nil }

// RequestDomainCertificate sends a request to retrieve a certificate to the CA.
// If the job succeeds, it marks the domain as issued and removes it from any
// processing queue. Otherwise, it leaves to the broker
// to decide what to do with the message.
func (p *DomainProcessor) RequestDomainCertificate(*Message) error { return nil }

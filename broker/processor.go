package broker

import (
	"github.com/lost-mountain/isard/certificates"
	"github.com/lost-mountain/isard/configuration"
	"github.com/lost-mountain/isard/domain"
	"github.com/lost-mountain/isard/domain/validator"
	"github.com/lost-mountain/isard/storage"
	"github.com/pkg/errors"
	"golang.org/x/crypto/acme"
)

// Processor defines an interface to process messages
// publised by the Broker.
type Processor interface {
	AuthorizeDomain(*Message) error
	CreateDomain(*Message) error
	ModifyDomain(*Message) error
	RequestDomainCertificate(*Message) error
	ValidateDomain(*Message) error
}

// DomainProcessor controls the lifecycle of a domain.
// It moves the domain across the state machine
// accordingly to the previous operation and its result.
type DomainProcessor struct {
	db     storage.Bucket
	broker Broker
	config configuration.DomainsConfiguration
}

// AuthorizeDomain sends an authorization request to the CA.
// If the job succeeds, it moves the domain to the
// certificate request state. Otherwise, it leaves to the broker
// to decide what to do with the message.
func (p *DomainProcessor) AuthorizeDomain(m *Message) error {
	v, ok := m.Payload.(*DomainPayload)
	if !ok {
		return errors.Errorf("error authorizing domain, invalid payload message: %v", m.Payload)
	}

	d, err := p.db.GetDomain(v.AccountID, v.DomainName)
	if err != nil {
		return err
	}

	c, err := certificates.NewClientWithAPIKey(d.Account, p.config.Ns1APIKey)
	if err != nil {
		return err
	}

	if d.AuthorizationURL != "" {
		return p.checkAuthzState(c, d)
	}
	return p.startAuthProcess(c, d)
}

// CreateDomain creates a new domain.
// If the job succeeds, it moves the domain to the
// validation state. Otherwise, it leaves to the broker
// to decide what to do with the message.
func (p *DomainProcessor) CreateDomain(m *Message) error {
	v, ok := m.Payload.(*CreateDomainPayload)
	if !ok {
		return errors.Errorf("error creating domain, invalid payload message: %v", m.Payload)
	}

	a, err := p.db.GetAccount(v.AccountID, v.AccountToken)
	if err != nil {
		return err
	}

	d, err := domain.NewDomainWithChallengeType(a, v.DomainName, v.ChallengeType)
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

// RequestDomainCertificate sends a request to retrieve a certificate to the CA.
// If the job succeeds, it marks the domain as issued and removes it from any
// processing queue. Otherwise, it leaves to the broker
// to decide what to do with the message.
func (p *DomainProcessor) RequestDomainCertificate(*Message) error { return nil }

// ValidateDomain validates that a domain is correctly configured
// before authorizing its issuing.
// If the job succeeds, it moves the domain to the
// authorization state. Otherwise, it leaves to the broker
// to decide what to do with the message.
func (p *DomainProcessor) ValidateDomain(m *Message) error {
	v, ok := m.Payload.(*DomainPayload)
	if !ok {
		return errors.Errorf("error verifying domain, invalid payload message: %v", m.Payload)
	}

	d, err := p.db.GetDomain(v.AccountID, v.DomainName)
	if err != nil {
		return err
	}

	next := domain.Verified
	var vErr error

	h := p.config.HeaderValidator
	valid := validator.ValidHeader(d.Name, h.Name, h.Value)

	if !valid {
		next = domain.Invalid
		vErr = errors.Errorf("domain validation failed for domain: %s", d.Name)
	}

	d.State = next
	if err := p.db.SaveDomain(d); err != nil {
		return err
	}

	return vErr
}

func (p *DomainProcessor) checkAuthzState(c *certificates.Client, d *domain.Domain) error {
	authz, err := c.GetAuthorization(d)
	if err != nil {
		return err
	}

	m := NewMessage(&DomainPayload{
		AccountID:  d.Account.ID,
		DomainName: d.Name,
	})

	if authz.Status != acme.StatusPending && authz.Status != acme.StatusProcessing {
		d.State = domain.Authorized
		if err := p.db.SaveDomain(d); err != nil {
			return err
		}

		return p.broker.Publish(CertRequest, m)
	}

	return p.broker.Publish(Authorization, m)
}

func (p *DomainProcessor) startAuthProcess(c *certificates.Client, d *domain.Domain) error {
	authz, err := c.AuthorizeDomain(d)
	if err != nil {
		return err
	}

	d.AuthorizationURL = authz.URI
	if err := p.db.SaveDomain(d); err != nil {
		return err
	}

	var chal *acme.Challenge
	for _, c := range authz.Challenges {
		if c.Type == d.ChallengeType {
			chal = c
			break
		}
	}

	if chal == nil {
		return errors.Errorf("unable to find a valid challenge for domain: %s", d.Name)
	}

	d, err = c.PrepareChallenge(d, chal)
	if err := p.db.SaveDomain(d); err != nil {
		return err
	}

	if _, err := c.AcceptChallenge(d, chal); err != nil {
		return err
	}

	m := NewMessage(&DomainPayload{
		AccountID:  d.Account.ID,
		DomainName: d.Name,
	})

	return p.broker.Publish(Authorization, m)
}

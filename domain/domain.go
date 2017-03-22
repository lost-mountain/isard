package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/cryptopolis"
	"github.com/pkg/errors"
)

// State is the status the certificate is in
// its lifecycle.
type State int

const (
	// Pending is the state of a domain when it enter the system.
	Pending State = iota
	// Validating is the state of a domain when the system is validating it.
	Validating
	// Invalid is the state of a domain when there is a validation error.
	Invalid
	// Verified is the state of a domain when the verification succeeded.
	Verified
	// Provisioning is the state of a domain when the certificate provisioning has started.
	Provisioning
	// Authorized is the state of a domain after the ACME authority authorizes the domain.
	Authorized
	// Requesting is the state of a domain when it's requesting the certificate.
	Requesting
	// Issued is the state of a domain after the certificate has been issued.
	Issued
	// Cancelling is the state of a domain when it's being cancelled.
	Cancelling
	// Cancelled is the state of a domain after the certificate has been cancelled.
	Cancelled

	defaultChallengeType = "http-01"
)

// ErrDuplicatedSANName is an error returned when a name
// already exists in the certificate's names list.
var ErrDuplicatedSANName = errors.Errorf("domain already includes SAN name")

// Domain stores information
// about a registered domain
// and its certificate authority.
type Domain struct {
	ID                      uuid.UUID
	Name                    string
	ChallengeType           string
	AuthorizationURL        string
	State                   State
	AccountID               string
	Account                 *account.Account
	CreatedAt               time.Time
	UpdatedAt               time.Time
	HTTP01ChallengePath     string
	HTTP01ChallengeResponse string

	Certificate *cryptopolis.Certificate

	initSAN map[string]struct{}
	san     map[string]struct{}
}

// AddSANName appends a name to the
// certificate's names list.
// It returns an error if the name is already
// in the list. This prevents reaching out
// limits with duplicated certificates.
func (d *Domain) AddSANName(name string) error {
	if d.san == nil {
		d.san = map[string]struct{}{}
	}
	if _, exist := d.san[name]; exist {
		return ErrDuplicatedSANName
	}
	d.san[name] = struct{}{}
	return nil
}

// RemoveSANName deletes a name from the
// certificate's names list.
// It doesn't allow to remove the initial
// name added when the domain was created.
func (d *Domain) RemoveSANName(name string) {
	if _, i := d.initSAN[name]; i {
		return
	}

	if d.san != nil {
		delete(d.san, name)
	}
}

// SANNames returns the certificate's names list.
func (d *Domain) SANNames() []string {
	if d.san == nil {
		return make([]string, 0)
	}
	names := make([]string, 0, len(d.san))
	for k := range d.san {
		names = append(names, k)
	}
	return names
}

// NewDomain initializes a new domain.
func NewDomain(account *account.Account, name string) (*Domain, error) {
	return NewDomainWithChallengeType(account, name, defaultChallengeType)
}

// NewDomainWithChallengeType initializes a new domain.
// It uses the Public Suffix list of domains to assing
// the name to the domain. The given name is added to the
// SAN names list.
func NewDomainWithChallengeType(account *account.Account, name, challengeType string) (*Domain, error) {
	names, err := ExtractNames(name)
	if err != nil {
		return nil, err
	}

	if challengeType == "" {
		challengeType = defaultChallengeType
	}

	d := &Domain{
		AccountID:     account.ID.String(),
		Account:       account,
		Name:          names.CN,
		State:         Pending,
		ChallengeType: challengeType,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		initSAN:       map[string]struct{}{names.CN: struct{}{}},
	}

	for _, s := range names.SAN {
		if err := d.AddSANName(s); err != nil {
			return nil, err
		}
	}
	return d, nil
}

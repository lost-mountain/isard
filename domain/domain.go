package domain

import (
	"github.com/google/uuid"
	"github.com/lost-mountain/isard/account"
)

// State is the status the certificate is in
// its lifecycle.
type State int

const (
	// Pending is the state of a domain when it enter the system.
	Pending State = iota
	// Invalid is the state of a domain when there is a validation error.
	Invalid
	// Provisioning is the state of a domain when the certificate provisioning has started.
	Provisioning
	// Authorized is the state of a domain after the ACME authority authorizes the domain.
	Authorized
	// Issued is the state of a domain after the certificate has been issued.
	Issued
	// Cancelled is the state of a domain after the certificate has been cancelled.
	Cancelled
)

// Domain stores information
// about a registered domain
// and its certificate authority.
type Domain struct {
	ID               uuid.UUID
	Name             string
	AuthorizationURL string
	State            State
	Account          *account.Account
}

// NewDomain initializes a new domain.
func NewDomain(account *account.Account, name string) *Domain {
	return &Domain{
		Account: account,
		Name:    name,
	}
}

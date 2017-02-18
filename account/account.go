package account

import (
	"github.com/google/uuid"
	"github.com/lost-mountain/isard/configuration"
)

// Account stores information
// about a registered account that
// issues domain certificates.
type Account struct {
	ID           uuid.UUID
	Token        uuid.UUID
	Certificate  string
	DirectoryURL string
	Owners       []string
}

// NewAccount initializes a new account for a set of owners.
func NewAccount(owners ...string) *Account {
	return NewAccountWithCertificate("", owners...)
}

// NewAccountWithCertificate initializes a new account for a set of owners
// with an existent certificate.
func NewAccountWithCertificate(certificate string, owners ...string) *Account {
	return &Account{
		ID:           uuid.New(),
		Token:        uuid.New(),
		Certificate:  certificate,
		DirectoryURL: configuration.StagingDirectory,
		Owners:       owners,
	}
}

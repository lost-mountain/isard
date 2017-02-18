package account

import (
	"crypto"
	"crypto/rand"

	"github.com/google/uuid"
	"github.com/lost-mountain/isard/configuration"
)

var rander = rand.Reader

// Account stores information
// about a registered account that
// issues domain certificates.
type Account struct {
	ID           uuid.UUID
	Token        uuid.UUID
	Key          string
	DirectoryURL string
	Owners       []string
}

// Contacts generates a list of mail contacts from the
// account owners.
func (a *Account) Contacts() []string {
	contacts := make([]string, len(a.Owners))
	for i, o := range a.Owners {
		contacts[i] = "mailto:" + o
	}
	return contacts
}

// PrivateKey returns the account's private key.
func (a *Account) PrivateKey() (crypto.Signer, error) {
	return extractPEMSigner(a.Key)
}

// NewAccount initializes a new account for a set of owners.
func NewAccount(owners ...string) (*Account, error) {
	return NewAccountWithKey("", owners...)
}

// NewAccountWithKey initializes a new account for a set of owners
// with an existent certificate.
func NewAccountWithKey(key string, owners ...string) (*Account, error) {
	account := &Account{
		ID:           uuid.New(),
		Token:        uuid.New(),
		Key:          key,
		DirectoryURL: configuration.StagingDirectory,
		Owners:       owners,
	}

	if key == "" {
		b, err := generatePEM()
		if err != nil {
			return nil, err
		}
		account.Key = string(b)
	} else {
		_, err := extractPEMSigner(key)
		if err != nil {
			return nil, err
		}
		account.Key = key
	}

	return account, nil
}

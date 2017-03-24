package account

import (
	"crypto"
	"time"

	"github.com/google/uuid"
	"github.com/lost-mountain/isard/cryptopolis"
)

// Account stores information
// about a registered account that
// issues domain certificates.
type Account struct {
	ID           uuid.UUID
	Token        uuid.UUID
	Key          string
	DirectoryURL string
	Owners       []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	return cryptopolis.ExtractPEMSigner(a.Key)
}

// NewAccount initializes a new account for a set of owners.
func NewAccount(owners ...string) (*Account, error) {
	return NewAccountWithKey("", owners...)
}

// NewAccountWithKey initializes a new account for a set of owners
// with an existent certificate.
func NewAccountWithKey(key string, owners ...string) (*Account, error) {
	account := &Account{
		ID:        uuid.New(),
		Token:     uuid.New(),
		Key:       key,
		Owners:    owners,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if key == "" {
		b, err := cryptopolis.GenerateECPrivateKeyPEM()
		if err != nil {
			return nil, err
		}
		account.Key = string(b)
	} else {
		_, err := cryptopolis.ExtractPEMSigner(key)
		if err != nil {
			return nil, err
		}
		account.Key = key
	}

	return account, nil
}

package storage

import (
	"github.com/google/uuid"
	"github.com/lost-mountain/isard/account"
)

// Bucket defines an interface to store information
// in a database.
type Bucket interface {
	Close() error
	GetAccount(id, token uuid.UUID) (*account.Account, error)
	SaveAccount(account *account.Account) error
}

// Configuration holds information to configure
// a storage backend.
type Configuration struct {
	URL string
}

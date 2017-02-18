package storage

import (
	"github.com/google/uuid"
	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/domain"
)

// Bucket defines an interface to store information
// in a database.
type Bucket interface {
	Close() error
	GetAccount(id, token uuid.UUID) (*account.Account, error)
	GetDomain(account *account.Account, name string) (*domain.Domain, error)
	SaveAccount(account *account.Account) error
	SaveDomain(domain *domain.Domain) error
}

// Configuration holds information to configure
// a storage backend.
type Configuration struct {
	URL string
}

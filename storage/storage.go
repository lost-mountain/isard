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
	GetDomain(accountID uuid.UUID, name string) (*domain.Domain, error)
	SaveAccount(account *account.Account) error
	SaveDomain(domain *domain.Domain) error
}

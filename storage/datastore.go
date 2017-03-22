package storage

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/domain"
	"github.com/pkg/errors"
)

// Datastore implements the Storage interface
// using Google Cloud Datastore as a backend.
type Datastore struct {
	client *datastore.Client
}

// Close closes the connection with the Datastore server.
func (d *Datastore) Close() error {
	return d.client.Close()
}

// GetAccount searches for an account with a given ID and Token.
func (d *Datastore) GetAccount(id, token uuid.UUID) (*account.Account, error) {
	key := datastore.NameKey("Account", id.String(), nil)
	var account account.Account
	if err := d.client.Get(context.Background(), key, &account); err != nil {
		return nil, errors.Wrapf(err, "error retrieving account %s", id)
	}

	return &account, nil
}

// GetDomain searches for a domain with a given name.
func (d *Datastore) GetDomain(accountID uuid.UUID, name string) (*domain.Domain, error) {
	query := datastore.NewQuery("Domain").
		Filter("AccountID =", accountID.String()).
		Filter("Name =", name).
		Limit(1)

	var res []*domain.Domain
	_, err := d.client.GetAll(context.Background(), query, &res)

	if err != nil {
		return nil, errors.Wrapf(err, "error retrieving domain %s", name)
	}

	if len(res) == 0 {
		return nil, errors.Wrapf(err, "error retrieving domain %s", name)
	}

	return res[0], nil
}

// SaveAccount saves an account in a bucket.
func (d *Datastore) SaveAccount(a *account.Account) error {
	key := datastore.NameKey("Account", a.ID.String(), nil)
	_, err := d.client.Put(context.Background(), key, a)

	return errors.Wrapf(err, "error saving account %s", a.ID)
}

// SaveDomain saves a domain in a bucket.
func (d *Datastore) SaveDomain(dm *domain.Domain) error {
	key := datastore.NameKey("Domain", dm.ID.String(), nil)
	_, err := d.client.Put(context.Background(), key, dm)

	return errors.Wrapf(err, "error saving domain %s", dm.ID)
}

// NewDatastoreBucket connects with the Datastore server.
func NewDatastoreBucket(projectID string) (*Datastore, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &Datastore{client}, nil
}

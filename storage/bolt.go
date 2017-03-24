package storage

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/domain"
	"github.com/pkg/errors"
)

// Bolt implements the Storage interface
// using BoltDB as a backend.
type Bolt struct {
	db *bolt.DB
}

// Close closes the connection with the BoltDB bucket.
func (b *Bolt) Close() error {
	return b.db.Close()
}

// GetAccount searches for an account with a given ID and Token.
func (b *Bolt) GetAccount(id, token uuid.UUID) (*account.Account, error) {
	var account account.Account

	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("accounts"))
		v := b.Get([]byte(id.String()))
		return json.Unmarshal(v, &account)
	})

	if err != nil {
		return nil, errors.Wrapf(err, "error retrieving account %s", id)
	}

	if account.Token != token {
		return nil, errors.Errorf("error retrieving account %s", id)
	}

	return &account, nil
}

// GetDomain searches for a domain with a given name.
func (b *Bolt) GetDomain(accountID uuid.UUID, name string) (*domain.Domain, error) {
	var domain domain.Domain

	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("domains"))

		key := fmt.Sprintf("%s@@%s", accountID, name)
		v := b.Get([]byte(key))

		return json.Unmarshal(v, &domain)
	})

	if err != nil {
		return nil, errors.Wrapf(err, "error retrieving domain %s", name)
	}

	return &domain, nil
}

// SaveAccount saves an account in a bucket.
func (b *Bolt) SaveAccount(a *account.Account) error {
	j, err := json.Marshal(a)
	if err != nil {
		return errors.Wrapf(err, "error saving account %s", a.ID)
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("accounts"))
		return b.Put([]byte(a.ID.String()), j)
	})

	return errors.Wrapf(err, "error saving account %s", a.ID)
}

// SaveDomain saves a domain in a bucket.
func (b *Bolt) SaveDomain(d *domain.Domain) error {
	j, err := json.Marshal(d)
	if err != nil {
		return errors.Wrapf(err, "error saving domain %s", d.ID)
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("domains"))
		key := fmt.Sprintf("%s@@%s", d.Account.ID, d.Name)
		return b.Put([]byte(key), j)
	})

	return errors.Wrapf(err, "error saving domain %s", d.ID)
}

// NewBoltBucket connects with a BoltDB database.
func NewBoltBucket(url string) (*Bolt, error) {
	db, err := bolt.Open(url, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := initializeBuckets(db); err != nil {
		return nil, err
	}

	return &Bolt{db}, nil
}

func initializeBuckets(db *bolt.DB) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.CreateBucketIfNotExists([]byte("accounts"))
	if err != nil {
		return err
	}

	_, err = tx.CreateBucketIfNotExists([]byte("domains"))
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

package storage

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/lost-mountain/isard/account"
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

// NewBoltBucket connects with a BoltDB database.
func NewBoltBucket(c *Configuration) (*Bolt, error) {
	db, err := bolt.Open(c.URL, 0600, nil)
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

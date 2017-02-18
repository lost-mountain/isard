package certificates

import (
	"context"

	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/domain"
	"github.com/pkg/errors"

	"golang.org/x/crypto/acme"
)

// Client uses an account to negotiate
// certificate operations with an ACME service.
type Client struct {
	account *account.Account
	client  *acme.Client
}

// NewClient initializes a new certificate client
// to handle ACME requests.
func NewClient(a *account.Account) (*Client, error) {
	pk, err := a.PrivateKey()
	if err != nil {
		return nil, err
	}

	c := &acme.Client{
		Key:          pk,
		DirectoryURL: a.DirectoryURL,
	}

	return &Client{a, c}, nil
}

// Register sends the account information to the ACME service.
func (c *Client) Register() error {
	ctx := context.Background()

	aa := &acme.Account{
		Contact: c.account.Contacts(),
	}

	_, err := c.client.Register(ctx, aa, acme.AcceptTOS)
	return err
}

// AuthorizeDomain initiates a domain name registration
// by sending the initial authorization.
func (c *Client) AuthorizeDomain(d *domain.Domain) (*acme.Authorization, error) {
	ctx := context.Background()
	return c.client.Authorize(ctx, d.Name)
}

// AcceptChallenge decides which challenge to use for a domain and uses
// the right resolver to accept the challenge.
func (c *Client) AcceptChallenge(d *domain.Domain, challenge *acme.Challenge) error {
	ctx := context.Background()
	_, err := c.client.Accept(ctx, challenge)
	if err != nil {
		return errors.Wrapf(err, "error accepting challenge for domain: %s", d.Name)
	}
	return nil
}

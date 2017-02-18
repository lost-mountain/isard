package certificates

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/lost-mountain/isard/account"

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

func certRequest(key crypto.Signer, domain string) ([]byte, error) {
	req := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: domain},
		DNSNames: []string{domain},
	}
	return x509.CreateCertificateRequest(rand.Reader, req, key)
}

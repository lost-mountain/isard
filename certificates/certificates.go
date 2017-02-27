package certificates

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"time"

	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/certificates/challenges"
	"github.com/lost-mountain/isard/domain"
	"github.com/pkg/errors"

	"golang.org/x/crypto/acme"
)

// Client uses an account to negotiate
// certificate operations with an ACME service.
type Client struct {
	account   *account.Account
	client    *acme.Client
	ns1ApiKey string
}

// AcceptChallenge sends the request to the ACME service to accept a challenge.
func (c *Client) AcceptChallenge(d *domain.Domain, chal *acme.Challenge) (*acme.Challenge, error) {
	ctx := context.Background()
	newC, err := c.client.Accept(ctx, chal)
	if err != nil {
		return nil, errors.Wrapf(err, "error accepting challenge for domain: %s", d.Name)
	}

	return newC, nil
}

// AuthorizeDomain initiates a domain name registration
// by sending the initial authorization.
func (c *Client) AuthorizeDomain(d *domain.Domain) (*acme.Authorization, error) {
	return c.client.Authorize(context.Background(), d.Name)
}

// GetAuthorization requests an authorization object that
// has already been issued.
func (c *Client) GetAuthorization(d *domain.Domain) (*acme.Authorization, error) {
	return c.client.GetAuthorization(context.Background(), d.AuthorizationURL)
}

// PrepareChallenge uses a challenge resolver to prepare a challenge.
func (c *Client) PrepareChallenge(d *domain.Domain, chal *acme.Challenge) (*domain.Domain, error) {
	var res challenges.Resolver
	switch chal.Type {
	case "dns-01":
		res = challenges.NewDNSResolver(c.ns1ApiKey, c.client)
	case "http-01":
		res = challenges.NewHTTPResolver(c.client)
	default:
		return nil, errors.Errorf("unsupported ACME challenge: %s", chal.Type)
	}

	d, err := res.Resolve(d, chal)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// RequestCertificate retrieves a certificate once the challenge has been completed
// and the authorization is valid.
func (c *Client) RequestCertificate(d *domain.Domain) (*x509.Certificate, error) {
	pk, err := d.Account.PrivateKey()
	if err != nil {
		return nil, errors.Wrapf(err, "error requesting certificate for domain: %s", d.Name)
	}

	req := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: d.Name},
		DNSNames: d.SANNames(),
	}

	cr, err := x509.CreateCertificateRequest(rand.Reader, req, pk)
	if err != nil {
		return nil, errors.Wrapf(err, "error requesting certificate for domain: %s", d.Name)
	}

	ctx := context.Background()
	der, _, err := c.client.CreateCert(ctx, cr, 0, true)
	if err != nil {
		return nil, errors.Wrapf(err, "error requesting certificate for domain: %s", d.Name)
	}

	l, err := validateCertificate(d.Name, pk, der)
	if err != nil {
		return nil, err
	}

	return l, nil
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

// NewClient initializes a new certificate client
// to handle ACME requests.
func NewClient(a *account.Account) (*Client, error) {
	return NewClientWithAPIKey(a, "")
}

// NewClientWithAPIKey initializes a new certificate client
// to handle ACME requests. It uses the api key
// to contact the DNS provider.
func NewClientWithAPIKey(a *account.Account, apiKey string) (*Client, error) {
	pk, err := a.PrivateKey()
	if err != nil {
		return nil, err
	}

	c := &acme.Client{
		Key:          pk,
		DirectoryURL: a.DirectoryURL,
	}

	return &Client{
		account:   a,
		ns1ApiKey: apiKey,
		client:    c,
	}, nil
}

// validateCertificate parses the certificate to ensure it's valid.
// Extracted from acme/autocert:
// https://github.com/golang/crypto/blob/9b1a210a06ea1176ec1f0a1ddf83ad7463b8ea3e/acme/autocert/autocert.go#L714
func validateCertificate(domain string, key crypto.Signer, der [][]byte) (*x509.Certificate, error) {
	var n int
	for _, b := range der {
		n += len(b)
	}

	pub := make([]byte, n)
	n = 0
	for _, b := range der {
		n += copy(pub[n:], b)
	}

	x509Cert, err := x509.ParseCertificates(pub)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid public key for domain: %s", domain)
	}
	if len(x509Cert) == 0 {
		return nil, errors.Errorf("invalid public key for domain: %s", domain)
	}

	leaf := x509Cert[0]
	now := time.Now()
	if now.Before(leaf.NotBefore) {
		return nil, errors.Errorf("invalid public key for domain: %s", domain)
	}
	if now.After(leaf.NotAfter) {
		return nil, errors.Errorf("expired certificate for domain: %s", domain)
	}
	if err := leaf.VerifyHostname(domain); err != nil {
		return nil, errors.Wrapf(err, "invalid hostname for certificate: %s", domain)
	}

	switch pubKey := leaf.PublicKey.(type) {
	case *rsa.PublicKey:
		prv, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.Errorf("mismatched public and private keys for certificate: %s", domain)
		}

		if pubKey.N.Cmp(prv.N) != 0 {
			return nil, errors.Errorf("mismatched public and private keys for certificate: %s", domain)
		}
	case *ecdsa.PublicKey:
		prv, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, errors.Errorf("mismatched public and private keys for certificate: %s", domain)
		}
		if pubKey.X.Cmp(prv.X) != 0 || pubKey.Y.Cmp(prv.Y) != 0 {
			return nil, errors.Errorf("mismatched public and private keys for certificate: %s", domain)
		}
	default:
		return nil, errors.Errorf("unsupported certificate type for domain: %s", domain)
	}

	return leaf, nil
}

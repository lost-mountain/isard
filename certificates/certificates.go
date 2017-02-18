package certificates

import (
	"context"

	"github.com/lost-mountain/isard/account"

	"golang.org/x/crypto/acme"
)

// RegisterAccount sends the account information to the ACME service.
func RegisterAccount(account *account.Account) error {
	pk, err := account.PrivateKey()
	if err != nil {
		return err
	}

	ctx := context.Background()
	c := &acme.Client{
		Key:          pk,
		DirectoryURL: account.DirectoryURL,
	}

	aa := &acme.Account{
		Contact: account.Contacts(),
	}

	_, err = c.Register(ctx, aa, acme.AcceptTOS)
	return err
}

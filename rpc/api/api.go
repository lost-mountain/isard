package api

import (
	"context"

	"github.com/lost-mountain/isard/account"
	"github.com/lost-mountain/isard/configuration"
	"github.com/lost-mountain/isard/rpc"
	"github.com/lost-mountain/isard/storage"
)

// API implements the GRPC server definition.
// It's the main entry point to manage accounts and domains.
type API struct {
	configuration *configuration.Configuration
	bucket        storage.Bucket
}

// CreateAccount creates a new domain account.
func (a *API) CreateAccount(ctx context.Context, req *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, error) {
	var (
		acc *account.Account
		err error
	)

	if req.Key != "" {
		acc, err = account.NewAccountWithKey(req.Key, req.Owner)
	} else {
		acc, err = account.NewAccount(req.Owner)
	}

	if err != nil {
		return nil, err
	}

	if req.Environment == rpc.AccountEnvironment_PRODUCTION {
		acc.DirectoryURL = a.configuration.ACME.DefaultProductionDirectory
	}

	if err := a.bucket.SaveAccount(acc); err != nil {
		return nil, err
	}

	return &rpc.CreateAccountResponse{
		Id:    acc.ID.String(),
		Token: acc.Token.String(),
	}, nil
}

// UpdateAccount updates the environment information for an account.
func (a *API) UpdateAccount(context.Context, *rpc.UpdateAccountRequest) (*rpc.UpdateAccountResponse, error) {
	return nil, nil
}

// CreateCertificate starts the process to request a domain certificate.
// It creates a new domain and negotiates the challenge type.
func (a *API) CreateCertificate(context.Context, *rpc.CreateCertificateRequest) (*rpc.CreateCertificateResponse, error) {
	return nil, nil
}

// ResolveCertificateChallenge returns the information required to resolve a challenge.
func (a *API) ResolveCertificateChallenge(context.Context, *rpc.ResolveChallengeRequest) (*rpc.ResolveChallengeResponse, error) {
	return nil, nil
}

// CheckCertificateState returns the state of a certificate.
func (a *API) CheckCertificateState(context.Context, *rpc.CertificateStateRequest) (*rpc.CertificateStateResponse, error) {
	return nil, nil
}

// GetCertificate returns the domain certificate once it has been authorized by the CA.
func (a *API) GetCertificate(context.Context, *rpc.GetCertificateRequest) (*rpc.GetCertificateResponse, error) {
	return nil, nil
}

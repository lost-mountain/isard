package account

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

func generatePEM() ([]byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rander)
	if err != nil {
		return nil, errors.Wrap(err, "error generating a key for a new account")
	}

	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "error generating a key for a new account")
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "ECDSA PRIVATE KEY",
			Bytes: b,
		},
	), nil
}

func extractPEMSigner(data string) (crypto.Signer, error) {
	block, _ := pem.Decode([]byte(data))
	if block == nil {
		return nil, errors.Errorf("error decoding PEM block")
	}

	switch block.Type {
	case "ECDSA PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, errors.Errorf("invalid private key type: %s", block.Type)
	}
}

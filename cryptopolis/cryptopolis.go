package cryptopolis

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"

	"github.com/pkg/errors"
)

const (
	certBlockType  = "CERTIFICATE"
	ecdsaBlockType = "EC PRIVATE KEY"
	rsaBlockType   = "RSA PRIVATE TYPE"
)

var rander = rand.Reader

// GenerateECPrivateKeyPEM generates an ECDSA private key.
func GenerateECPrivateKeyPEM() ([]byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rander)
	if err != nil {
		return nil, errors.Wrap(err, "error generating a key for a new account")
	}

	var buf bytes.Buffer
	if err := encodeECDSAKeyPEM(&buf, key); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExtractPEMSigner decodes a private key in PEM format.
func ExtractPEMSigner(data string) (crypto.Signer, error) {
	block, _ := pem.Decode([]byte(data))
	if block == nil {
		return nil, errors.Errorf("error decoding PEM block")
	}

	switch block.Type {
	case ecdsaBlockType:
		return x509.ParseECPrivateKey(block.Bytes)
	case rsaBlockType:
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, errors.Errorf("invalid private key type: %s", block.Type)
	}
}

// EncodeCertificate encodes a certificate and its chain in PEM format.
func EncodeCertificate(cert *tls.Certificate) ([]byte, error) {
	var (
		buf bytes.Buffer
		err error
	)

	switch key := cert.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		err = encodeECDSAKeyPEM(&buf, key)
	case *rsa.PrivateKey:
		err = encodeRSAKeyPEM(&buf, key)
	default:
		return nil, errors.Errorf("invalid private key type encoding certificate")
	}

	if err != nil {
		return nil, err
	}

	for _, b := range cert.Certificate {
		if err := encodeCertPEM(&buf, b); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func encodeECDSAKeyPEM(w io.Writer, key *ecdsa.PrivateKey) error {
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return errors.Wrap(err, "error encoding ECDSA key")
	}

	block := &pem.Block{
		Type:  ecdsaBlockType,
		Bytes: b,
	}

	if err := pem.Encode(w, block); err != nil {
		return errors.Wrap(err, "error encoding ECDSA key")
	}
	return nil
}

func encodeRSAKeyPEM(w io.Writer, key *rsa.PrivateKey) error {
	b := x509.MarshalPKCS1PrivateKey(key)
	pb := &pem.Block{
		Type:  rsaBlockType,
		Bytes: b,
	}

	if err := pem.Encode(w, pb); err != nil {
		return errors.Wrap(err, "error encoding RSA key")
	}
	return nil
}

func encodeCertPEM(w io.Writer, b []byte) error {
	pb := &pem.Block{Type: certBlockType, Bytes: b}
	if err := pem.Encode(w, pb); err != nil {
		return errors.Wrap(err, "error encoding public certificate")
	}
	return nil
}

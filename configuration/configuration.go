package configuration

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

const (
	// productionDirectory is the address where the production ACME directory is.
	productionDirectory = "http://localhost:4000/directory"
	// stagingDirectory is the address where the staging ACME directory is.
	stagingDirectory = "http://localhost:4000/directory"

	defaultTCPPort = 8473
)

// Configuration hold setup
// information for the service
// to work.
type Configuration struct {
	ACME struct {
		DefaultProductionDirectory string
		DefaultStagingDirectory    string
	}

	TCP struct {
		Port int
	}

	TLS *struct {
		CertFile string
		KeyFile  string
	}

	GC *struct {
		Project        string
		AccountKeyFile string
	}

	Domains *DomainsConfiguration
}

// DomainsConfiguration holds setup
// information to request certificates
// and validate domains.
type DomainsConfiguration struct {
	Ns1APIKey       string
	HeaderValidator struct {
		Name  string
		Value string
	}
}

// Load parses a file to generate
// a configuration structure.
func Load(p string) (*Configuration, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open configuration file: %s", p)
	}

	var c Configuration
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, errors.Wrapf(err, "cannot parse configuration file: %s", p)
	}

	if c.ACME.DefaultProductionDirectory == "" {
		c.ACME.DefaultProductionDirectory = productionDirectory
	}

	if c.ACME.DefaultStagingDirectory == "" {
		c.ACME.DefaultStagingDirectory = stagingDirectory
	}

	if c.TCP.Port == 0 {
		c.TCP.Port = 8080
	}

	return &c, nil
}

package configuration

const (
	// ProductionDirectory is the address where the production ACME directory is.
	ProductionDirectory = "http://localhost:4000/directory"
	// StagingDirectory is the address where the staging ACME directory is.
	StagingDirectory = "http://localhost:4000/directory"
)

// Configuration hold setup
// information for the service
// to work.
type Configuration struct {
	Domains DomainsConfiguration
	ACME    struct {
		DefaultProductionDirectory string
		DefaultStagingDirectory    string
	}
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

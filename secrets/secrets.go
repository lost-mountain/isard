package secrets

// Secrets is an interface that defines
// how to encrypt and decrypt data.
type Secrets interface {
	Encrypt(data string) string
	Decrypt(data string) string
}

// Noop is a secrets backend that
// doesn't execute any encryption.
type Noop struct{}

// Encrypt returns the same data that receives.
func (Noop) Encrypt(data string) string { return data }

// Decrypt returns the same data that receives.
func (Noop) Decrypt(data string) string { return data }

package broker

// TopicType defines the operations
// the broker knows about.
// They are used as a state machine
// to control the domain's lifecycle.
type TopicType string

const (
	// Creation is the topic to create new domains.
	Creation TopicType = "creation"
	// Modification is the topic to modify domains and their certificate names.
	Modification TopicType = "modification"
	// Validation is the topic to validate domains before authorize a certificate.
	Validation TopicType = "validation"
	// Authorization is the topic to authorize domain certificates.
	Authorization TopicType = "authorization"
	// CertRequest is the topic to request domain certificates after they have been authorized.
	CertRequest TopicType = "cert_request"
)

var allTopics = []TopicType{
	Creation,
	Modification,
	Validation,
	Authorization,
	CertRequest,
}

// Broker defines an interface to publish
// messages in a queue and subscribe to them.
type Broker interface {
	Close() error
	Publish(topic TopicType, payload interface{}) error
	Subscribe(processor Processor) error
}

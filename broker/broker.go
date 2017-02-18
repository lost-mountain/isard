package broker

// Processor defines an interface to process messages
// publised by the Broker.
type Processor func(message *Message) error

// Broker defines an interface to publish
// messages in a queue and subscribe to them.
type Broker interface {
	Close() error
	Publish(topic string, payload interface{}) error
	Subscribe(processor Processor) error
}

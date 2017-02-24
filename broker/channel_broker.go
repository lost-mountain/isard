package broker

import "io"

// ChannelBroker implements broker.Broker using channels as a backend.
// This interface is only suitable for testing.
// It offers no guarantees about the elements pushed and pulled from the queue.
type ChannelBroker struct {
	c map[TopicType]chan *Message
	e chan struct{}
}

// NewChannelBroker initializes the channel broker.
func NewChannelBroker() *ChannelBroker {
	c := map[TopicType]chan *Message{}
	for _, t := range allTopics {
		c[t] = make(chan *Message)
	}

	return &ChannelBroker{
		c: c,
		e: make(chan struct{}),
	}
}

// Close sends a message to the exit channel
// to stop processing messages.
func (b *ChannelBroker) Close() error {
	b.e <- struct{}{}
	return nil
}

// Publish sends messages to the channel for a specific job.
func (b *ChannelBroker) Publish(topic TopicType, payload io.Reader) error {
	m := NewMessage(payload)
	b.c[topic] <- m
	return nil
}

// Subscribe receives messages from the channel to process them.
func (b *ChannelBroker) Subscribe(processor Processor) {
	go func() {
		for {
			select {
			case msg := <-b.c[Creation]:
				processor.CreateDomain(msg)
			case msg := <-b.c[Modification]:
				processor.ModifyDomain(msg)
			case msg := <-b.c[Validation]:
				processor.ValidateDomain(msg)
			case msg := <-b.c[Authorization]:
				processor.AuthorizeDomain(msg)
			case msg := <-b.c[CertRequest]:
				processor.RequestDomainCertificate(msg)
			case <-b.e:
				break
			}
		}
	}()
}

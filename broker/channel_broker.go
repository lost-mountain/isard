package broker

import "github.com/google/uuid"

type topicMessage struct {
	Topic string
	*Message
}

// ChannelBroker implements broker.Broker using a channel as a backend.
// This interface is only suitable for testing.
// It offers no guarantees about the elements pushed and pulled from the queue.
type ChannelBroker struct {
	c chan *topicMessage
	e chan struct{}
}

// NewChannelBroker initializes the channel broker.
func NewChannelBroker() *ChannelBroker {
	return &ChannelBroker{
		c: make(chan *topicMessage),
	}
}

// Close sends a message to the exit channel
// to stop processing messages.
func (b *ChannelBroker) Close() error {
	b.e <- struct{}{}
	return nil
}

// Publish sends messages to the channel for a specific job.
func (b *ChannelBroker) Publish(topic string, payload interface{}) error {
	m := NewMessage(uuid.New(), payload)
	b.c <- &topicMessage{topic, m}
	return nil
}

// Subscribe receives messages from the channel to process them.
func (b *ChannelBroker) Subscribe(processor Processor) {
	go func() {
		for {
			select {
			case msg := <-b.c:
				go processor(msg.Message)
			case <-b.e:
				break
			}
		}
	}()
}

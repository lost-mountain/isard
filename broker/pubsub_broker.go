package broker

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"

	"github.com/pkg/errors"

	"cloud.google.com/go/pubsub"
)

// PubSubBroker implements broker.Broker using Google Cloud's PubSub as a backend.
type PubSubBroker struct {
	client      *pubsub.Client
	subs        *pubsub.Subscription
	subsCancel  context.CancelFunc
	subsContext context.Context
}

// NewPubSubBroker initializes a new broker and stablish a
// client connection with the remote endpoint.
// It returns an error if connection cannot be stablished.
func NewPubSubBroker(projectID string) (*PubSubBroker, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	cctx, cancel := context.WithCancel(ctx)
	b := &PubSubBroker{
		client:      client,
		subsCancel:  cancel,
		subsContext: cctx,
	}

	if err := b.setupPubSub(); err != nil {
		return nil, errors.Wrap(err, "unable to create PubSub topics")
	}

	return b, nil
}

// Close cancels the subscription listener.
func (b *PubSubBroker) Close() error {
	b.subsCancel()
	return nil
}

// Publish sends messages to the broker for a specific job.
func (b *PubSubBroker) Publish(topic TopicType, payload interface{}) error {
	m := NewMessage(payload)
	m.Topic = topic

	d, err := json.Marshal(m)
	if err != nil {
		return err
	}

	t := b.client.Topic(string(topic))
	ctx := context.Background()
	_, err = t.Publish(ctx, &pubsub.Message{
		Data: d,
	}).Get(ctx)

	return err
}

// Subscribe receives messages from the broker to process them.
func (b *PubSubBroker) Subscribe(processor Processor) error {
	go func() {
		b.subs.Receive(b.subsContext, func(ctx context.Context, msg *pubsub.Message) {
			var umsg Message
			if err := json.Unmarshal(msg.Data, &umsg); err != nil {
				msg.Ack()
				return
			}

			var err error
			switch umsg.Topic {
			case Creation:
				err = processor.CreateDomain(&umsg)
			case Modification:
				err = processor.ModifyDomain(&umsg)
			case Validation:
				err = processor.ValidateDomain(&umsg)
			case Authorization:
				err = processor.AuthorizeDomain(&umsg)
			case CertRequest:
				err = processor.RequestDomainCertificate(&umsg)
			}

			if err != nil {
				msg.Nack()
				return
			}

			msg.Ack()
		})
	}()
	return nil
}

func (b *PubSubBroker) setupPubSub() error {
	ctx := context.Background()
	topic, err := b.client.CreateTopic(ctx, "isard-topic")
	if err != nil {
		return err
	}

	s, err := b.client.CreateSubscription(ctx, "isard-subscription", topic, 20*time.Second, nil)
	if err != nil {
		return err
	}

	b.subs = s

	return nil
}

package broker

import (
	"io"

	"github.com/google/uuid"
)

// Message is the structure that the broker sends and receives.
type Message struct {
	JobUUID uuid.UUID // unique identifiler for the job that trigerred this message
	Payload io.Reader
}

// NewMessage creates new messages with default ids.
func NewMessage(payload io.Reader) *Message {
	return NewMessageWithID(uuid.New(), payload)
}

// NewMessageWithID creates a new message with a given id.
func NewMessageWithID(jobID uuid.UUID, payload io.Reader) *Message {
	return &Message{
		JobUUID: jobID,
		Payload: payload,
	}
}

package broker

import "github.com/google/uuid"

// Message is the structure that the broker sends and receives.
type Message struct {
	JobUUID uuid.UUID // unique identifiler for the job that trigerred this message
	Payload interface{}
}

// NewMessage creates new messages.
func NewMessage(jobID uuid.UUID, payload interface{}) *Message {
	return &Message{
		JobUUID: jobID,
		Payload: payload,
	}
}

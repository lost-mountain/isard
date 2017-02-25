package broker

import "github.com/google/uuid"

// CreateDomainPayload is the payload
// sent by a client to create a new
// domain.
type CreateDomainPayload struct {
	AccountID     uuid.UUID `json:"account_id"`
	AccountToken  uuid.UUID `json:"account_token"`
	DomainName    string    `json:"domain_name"`
	ChallengeType string    `json:"challenge_type"`
}

// DomainPayload is the payload
// sent by a client to verify a
// domain.
type DomainPayload struct {
	AccountID  uuid.UUID `json:"account_id"`
	DomainName string    `json:"domain_name"`
}

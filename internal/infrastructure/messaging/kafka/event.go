package kafka

import (
	"time"

	"github.com/google/uuid"
)

type Event[T any] struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Payload   T         `json:"payload"`
}

func NewEvent[T any](tenantID, eventType string, payload T) Event[T] {
	return Event[T]{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
}

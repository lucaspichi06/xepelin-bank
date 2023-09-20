package domain

import (
	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID  `json:"transaction_id"`
	AccountID     uuid.UUID  `json:"account_id" binding:"required"`
	DestinationID *uuid.UUID `json:"destination_id,omitempty"`
	Type          EventType  `json:"type" binding:"required"`
	Amount        float64    `json:"amount" binding:"required"`
	Timestamp     string     `json:"timestamp"`
}

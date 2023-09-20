package domain

import "github.com/google/uuid"

const (
	Create   EventType = "create"
	Deposit  EventType = "deposit"
	WithDraw EventType = "withdraw"
	Transfer EventType = "transfer"
	Balance  EventType = "balance"
)

type Event interface {
	Process() (Account, error)
}

type EventType string

type DefaultEvent struct {
	AccId uuid.UUID
	Type  EventType
}
